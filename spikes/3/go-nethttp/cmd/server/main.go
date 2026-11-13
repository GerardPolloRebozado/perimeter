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

	"github.com/hibiken/asynq"
	"golang.org/x/sync/errgroup"

	"surface-volume/spike-3/go-nethttp/internal/auth"
	"surface-volume/spike-3/go-nethttp/internal/cfg"
	"surface-volume/spike-3/go-nethttp/internal/db"
	"surface-volume/spike-3/go-nethttp/internal/db/sqlc"
	"surface-volume/spike-3/go-nethttp/internal/api"
	"surface-volume/spike-3/go-nethttp/internal/jobs"
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

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})
	mux.HandleFunc("GET /api/v1/ping", api.New(c, q, asynqClient).Ping)
	mux.HandleFunc("POST /api/v1/hooks", api.New(c, q, asynqClient).ReceiveHook)

	var root http.Handler = mux
	if authMW != nil {
		root = authMW(mux)
	}

	srv := &http.Server{
		Addr:              ":" + c.Port,
		Handler:           root,
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
