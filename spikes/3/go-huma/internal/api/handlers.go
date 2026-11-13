package api

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/hibiken/asynq"
	openapi_types "github.com/oapi-codegen/runtime/types"

	"surface-volume/spike-3/go-huma/internal/cfg"
	"surface-volume/spike-3/go-huma/internal/db/sqlc"
	"surface-volume/spike-3/go-huma/internal/jobs"
)

type Handler struct {
	cfg  *cfg.Config
	q    *sqlc.Queries
	jobs *asynq.Client
}

func New(c *cfg.Config, q *sqlc.Queries, j *asynq.Client) *Handler {
	return &Handler{cfg: c, q: q, jobs: j}
}

func (h *Handler) Ping(ctx context.Context) (map[string]any, error) {
	count, err := h.q.CountHooks(ctx)
	if err != nil {
		slog.Error("count hooks", "err", err)
		return nil, err
	}
	return map[string]any{
		"status":     "pong",
		"service":    h.cfg.ServiceName,
		"version":    h.cfg.Version,
		"hook_count": count,
	}, nil
}

func (h *Handler) ReceiveHook(ctx context.Context, body map[string]any) (map[string]any, error) {
	source, _ := body["source"].(string)
	if source == "" {
		return nil, errSourceRequired
	}

	payload := []byte("{}")
	if p, ok := body["payload"].(map[string]any); ok {
		b, err := json.Marshal(p)
		if err != nil {
			return nil, err
		}
		payload = b
	}

	hook, err := h.q.InsertHook(ctx, payload)
	if err != nil {
		slog.Error("insert hook", "err", err)
		return nil, err
	}

	uid := openapi_types.UUID(hook.ID.Bytes)
	if _, err := jobs.EnqueueHook(ctx, h.jobs, jobs.HookPayload{
		HookID: uid.String(),
		Source: source,
	}); err != nil {
		slog.Error("enqueue", "err", err)
		return nil, err
	}

	return map[string]any{
		"id":     uid,
		"queued": true,
	}, nil
}

var errSourceRequired = &humError{msg: "source required", code: 400}

type humError struct {
	msg  string
	code int
}

func (e *humError) Error() string { return e.msg }
