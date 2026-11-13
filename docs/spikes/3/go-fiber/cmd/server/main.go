package main

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/hibiken/asynq"
	"golang.org/x/sync/errgroup"

	"perimeter-volume/spike-3/go-fiber/internal/api"
	"perimeter-volume/spike-3/go-fiber/internal/cfg"
	"perimeter-volume/spike-3/go-fiber/internal/db"
	"perimeter-volume/spike-3/go-fiber/internal/db/sqlc"
	"perimeter-volume/spike-3/go-fiber/internal/jobs"
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

	// ISSUE: Fiber is built on fasthttp, NOT net/http
	// The go-oidc JWKS middleware (and oapi-codegen std-http-server templates,
	// and most net/http-native Go libraries) take *http.Request and return
	// http.Handler. fiber.Ctx is not an *http.Request.
	var verifierConfigured bool
	if c.OIDCIssuerURL != "" {
		verifierConfigured = true
		slog.Info("OIDC issuer configured but JWKS middleware is a no-op in Fiber",
			"reason", "fasthttp is not net/http; auth middleware requires adaptor layer")
	} else {
		slog.Warn("auth disabled: OIDC_ISSUER_URL unset")
	}

	app := fiber.New(fiber.Config{
		ReadTimeout: 5 * time.Second,
	})

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	h := api.New(c, q, asynqClient)
	app.Get("/api/v1/ping", h.Ping)

	api := app.Group("/api/v1")
	if verifierConfigured {
		api.Use(func(c *fiber.Ctx) error {
			if c.Path() == "/ping" {
				return c.Next()
			}
			// Fiber cannot directly use net/http middleware (go-oidc, keyfunc,
			// oapi-codegen std-http-server, etc.) because fiber.Ctx != *http.Request.
			return c.Next()
		})
	}
	app.Post("/api/v1/hooks", h.ReceiveHook)

	g, gctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		slog.Info("http listening", "addr", ":"+c.Port)
		if err := app.Listen(":" + c.Port); err != nil && !errors.Is(err, context.Canceled) {
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
		return app.ShutdownWithContext(shutCtx)
	})

	return g.Wait()
}
