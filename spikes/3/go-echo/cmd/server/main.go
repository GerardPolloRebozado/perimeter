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
	"github.com/hibiken/asynq"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/sync/errgroup"

	"perimeter-volume/spike-3/go-echo/internal/api"
	"perimeter-volume/spike-3/go-echo/internal/cfg"
	"perimeter-volume/spike-3/go-echo/internal/db"
	"perimeter-volume/spike-3/go-echo/internal/db/sqlc"
	"perimeter-volume/spike-3/go-echo/internal/jobs"
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

	// Echo middleware is echo.MiddlewareFunc (func(next echo.HandlerFunc) echo.HandlerFunc),
	// so we write native Echo middleware instead of adapting net/http middleware
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

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	h := api.New(c, q, asynqClient)

	// Echo-native JWT middleware: validates Bearer token using go-oidc verifier,
	// skipping /ping and /health.
	e.Group("/api/v1").Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Path() == "/ping" || c.Path() == "/health" {
				return next(c)
			}
			if verifier == nil {
				return next(c)
			}
			h := c.Request().Header.Get("Authorization")
			if !strings.HasPrefix(h, "Bearer ") {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "missing bearer token"})
			}
			raw := strings.TrimPrefix(h, "Bearer ")
			if _, err := verifier.Verify(c.Request().Context(), raw); err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid token"})
			}
			return next(c)
		}
	})

	e.GET("/api/v1/ping", h.Ping)
	e.POST("/api/v1/hooks", h.ReceiveHook)

	srv := &http.Server{
		Addr:              ":" + c.Port,
		Handler:           e,
		ReadHeaderTimeout: 5 * time.Second,
	}

	g, gctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		slog.Info("http listening", "addr", srv.Addr)
		if err := e.StartServer(srv); err != nil && !errors.Is(err, http.ErrServerClosed) {
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
		if err := e.Shutdown(shutCtx); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	})

	return g.Wait()
}
