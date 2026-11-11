package main

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/danielgtaylor/huma/v2"
	"github.com/hibiken/asynq"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"golang.org/x/sync/errgroup"

	"perimeter-volume/spike-3/go-huma/internal/cfg"
	"perimeter-volume/spike-3/go-huma/internal/db"
	"perimeter-volume/spike-3/go-huma/internal/db/sqlc"
	"perimeter-volume/spike-3/go-huma/internal/jobs"
)

// operations are registered in Go and Huma emits the
// OpenAPI spec from them

type PingInput struct{}

type PingOutput struct {
	Body struct {
		Status    string `json:"status"`
		Service   string `json:"service"`
		Version   string `json:"version"`
		HookCount int64  `json:"hook_count"`
	}
}

type HookInput struct {
	Body struct {
		Source  string         `json:"source"`
		Payload map[string]any `json:"payload"`
	}
}

type HookOutput struct {
	Body struct {
		Id     string `json:"id"`
		Queued bool   `json:"queued"`
	}
}

type Handler struct {
	cfg  *cfg.Config
	q    *sqlc.Queries
	jobs *asynq.Client
}

func New(c *cfg.Config, q *sqlc.Queries, j *asynq.Client) *Handler {
	return &Handler{cfg: c, q: q, jobs: j}
}

func (h *Handler) Ping(ctx context.Context, input *PingInput) (*PingOutput, error) {
	count, err := h.q.CountHooks(ctx)
	if err != nil {
		return nil, huma.NewError(http.StatusInternalServerError, "db error")
	}
	return &PingOutput{
		Body: struct {
			Status    string `json:"status"`
			Service   string `json:"service"`
			Version   string `json:"version"`
			HookCount int64  `json:"hook_count"`
		}{
			Status:    "pong",
			Service:   h.cfg.ServiceName,
			Version:   h.cfg.Version,
			HookCount: count,
		},
	}, nil
}

func (h *Handler) ReceiveHook(ctx context.Context, input *HookInput) (*HookOutput, error) {
	if input.Body.Source == "" {
		return nil, huma.NewError(http.StatusBadRequest, "source required")
	}

	payload := []byte("{}")
	if input.Body.Payload != nil {
		b, err := json.Marshal(input.Body.Payload)
		if err != nil {
			return nil, huma.NewError(http.StatusBadRequest, "marshal payload")
		}
		payload = b
	}

	hook, err := h.q.InsertHook(ctx, payload)
	if err != nil {
		return nil, huma.NewError(http.StatusInternalServerError, "db error")
	}

	uid := openapi_types.UUID(hook.ID.Bytes)
	if _, err := jobs.EnqueueHook(ctx, h.jobs, jobs.HookPayload{
		HookID: uid.String(),
		Source: input.Body.Source,
	}); err != nil {
		return nil, huma.NewError(http.StatusInternalServerError, "queue error")
	}

	return &HookOutput{
		Body: struct {
			Id     string `json:"id"`
			Queued bool   `json:"queued"`
		}{
			Id:     uid.String(),
			Queued: true,
		},
	}, nil
}

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

	h := New(c, q, asynqClient)

	// define operations in Go, Huma emits OpenAPI
	// NewAPI with nil adapter uses the built-in net/http adapter
	api := huma.NewAPI(huma.DefaultConfig("AREA API (spike: go-huma)", "0.0.1"), nil)

	huma.Get(api, "/api/v1/ping", h.Ping)
	huma.Post(api, "/api/v1/hooks", h.ReceiveHook)

	// api.Adapter() returns an Adapter with ServeHTTP, which is an http.Handler
	var root http.Handler = api.Adapter()
	if verifier != nil {
		root = func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/health" || r.URL.Path == "/api/v1/ping" {
					next.ServeHTTP(w, r)
					return
				}
				h := r.Header.Get("Authorization")
				if !strings.HasPrefix(h, "Bearer ") {
					http.Error(w, "missing bearer token", http.StatusUnauthorized)
					return
				}
				raw := strings.TrimPrefix(h, "Bearer ")
				if _, err := verifier.Verify(r.Context(), raw); err != nil {
					http.Error(w, "invalid token", http.StatusUnauthorized)
					return
				}
				next.ServeHTTP(w, r)
			})
		}(api.Adapter())
	}

	// Healthz endpoint
	health := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	mux := http.NewServeMux()
	mux.Handle("/health", health)
	mux.Handle("/", root)

	srv := &http.Server{
		Addr:              ":" + c.Port,
		Handler:           mux,
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
