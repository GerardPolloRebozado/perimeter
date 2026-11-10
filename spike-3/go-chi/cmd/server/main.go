package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/hibiken/asynq"
	"golang.org/x/sync/errgroup"

	"perimeter-volume/spike-3/go-chi/internal/auth"
	"perimeter-volume/spike-3/go-chi/internal/cfg"
	"perimeter-volume/spike-3/go-chi/internal/db"
	"perimeter-volume/spike-3/go-chi/internal/db/sqlc"
	"perimeter-volume/spike-3/go-chi/internal/api"
	"perimeter-volume/spike-3/go-chi/internal/jobs"
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

	authMW, err := auth.NewJWTMiddleware(ctx, c.OIDCIssuerURL)
	if err != nil {
		return err
	}
	if authMW == nil {
		slog.Warn("auth disabled: OIDC_ISSUER_URL unset")
	}

	r := chi.NewRouter()
	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	h := api.New(c, q, asynqClient)
	r.Get("/api/v1/ping", h.Ping)

	// Auth middleware applies net/http-native JWKS validation directly (chi is http.Handler-native)
	r.Route("/api/v1", func(r chi.Router) {
		if authMW != nil {
			r.Use(authMW)
		}
		r.Post("/hooks", h.ReceiveHook)
	})

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
