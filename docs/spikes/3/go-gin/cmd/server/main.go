package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
	"github.com/hibiken/asynq"
	"golang.org/x/sync/errgroup"

	"surface-volume/spike-3/go-gin/internal/api"
	"surface-volume/spike-3/go-gin/internal/cfg"
	"surface-volume/spike-3/go-gin/internal/db"
	"surface-volume/spike-3/go-gin/internal/db/sqlc"
	"surface-volume/spike-3/go-gin/internal/jobs"
)

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, nil)))
	if err := run(); err != nil {
		slog.Error("fatal", "err", err)
		os.Exit(1)
	}
}

func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	c, err := cfg.Load()
	if err != nil {
		return err
	}

	pool, err := db.New(ctx, c.DatabaseURL)
	if err != nil {
		return err
	}
	defer pool.Close()
	q := sqlc.New(pool)

	redisOpt := asynq.RedisClientOpt{Addr: c.RedisAddr}
	asynqClient := asynq.NewClient(redisOpt)
	defer asynqClient.Close()

	worker := jobs.NewWorker(slog.Default())
	asynqServer := asynq.NewServer(redisOpt, asynq.Config{
		Concurrency:     4,
		ShutdownTimeout: 8 * time.Second,
	})

	// Build go-oidc verifier once. Gin middleware is gin.HandlerFunc, so we
	// write native Gin middleware
	var verifier *oidc.IDTokenVerifier
	if c.OIDCIssuerURL != "" {
		provider, err := oidc.NewProvider(ctx, c.OIDCIssuerURL)
		if err != nil {
			return err
		}
		verifier = provider.Verifier(&oidc.Config{
			SkipClientIDCheck: true,
		})
	} else {
		slog.Warn("auth disabled: OIDC_ISSUER_URL unset")
	}

	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	h := api.New(c, q, asynqClient)
	r.GET("/api/v1/ping", h.Ping)

	// Gin-native JWT middleware: validates Bearer token using go-oidc verifier,
	// skipping /ping and /health. gin.HandlerFunc uses *gin.Context
	api := r.Group("/api/v1")
	if verifier != nil {
		api.Use(func(c *gin.Context) {
			if c.FullPath() == "/api/v1/ping" || c.FullPath() == "/health" {
				c.Next()
				return
			}
			h := c.Request.Header.Get("Authorization")
			if !strings.HasPrefix(h, "Bearer ") {
				c.JSON(http.StatusUnauthorized, map[string]string{"error": "missing bearer token"})
				c.Abort()
				return
			}
			raw := strings.TrimPrefix(h, "Bearer ")
			if _, err := verifier.Verify(c.Request.Context(), raw); err != nil {
				c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid token"})
				c.Abort()
				return
			}
			c.Next()
		})
	}
	api.POST("/hooks", h.ReceiveHook)

	srv := &http.Server{
		Addr:              ":" + c.Port,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
	}

	g, gctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		slog.Info("http listening", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	})

	g.Go(func() error {
		slog.Info("asynq worker starting")
		if err := asynqServer.Run(worker.Handler()); err != nil {
			return err
		}
		return nil
	})

	g.Go(func() error {
		<-gctx.Done()
		slog.Info("shutdown: asynq")
		asynqServer.Shutdown()
		return nil
	})

	g.Go(func() error {
		<-gctx.Done()
		slog.Info("shutdown: http")
		shutCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return srv.Shutdown(shutCtx)
	})

	return g.Wait()
}
