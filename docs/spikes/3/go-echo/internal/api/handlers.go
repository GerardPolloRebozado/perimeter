package api

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/hibiken/asynq"
	openapi_types "github.com/oapi-codegen/runtime/types"
	"github.com/labstack/echo/v4"

	"surface-volume/spike-3/go-echo/internal/cfg"
	"surface-volume/spike-3/go-echo/internal/db/sqlc"
	"surface-volume/spike-3/go-echo/internal/jobs"
)

type Handler struct {
	cfg  *cfg.Config
	q    *sqlc.Queries
	jobs *asynq.Client
}

func New(c *cfg.Config, q *sqlc.Queries, j *asynq.Client) *Handler {
	return &Handler{cfg: c, q: q, jobs: j}
}

func (h *Handler) Ping(c echo.Context) error {
	count, err := h.q.CountHooks(c.Request().Context())
	if err != nil {
		slog.Error("count hooks", "err", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "db error"})
	}
	return c.JSON(http.StatusOK, map[string]any{
		"status":     "pong",
		"service":    h.cfg.ServiceName,
		"version":    h.cfg.Version,
		"hook_count": count,
	})
}

func (h *Handler) ReceiveHook(c echo.Context) error {
	var body struct {
		Source  string         `json:"source"`
		Payload map[string]any `json:"payload"`
	}
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid json"})
	}
	if body.Source == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "source required"})
	}

	payload := []byte("{}")
	if body.Payload != nil {
		b, err := json.Marshal(body.Payload)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "marshal payload"})
		}
		payload = b
	}

	hook, err := h.q.InsertHook(c.Request().Context(), payload)
	if err != nil {
		slog.Error("insert hook", "err", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "db error"})
	}

	uid := openapi_types.UUID(hook.ID.Bytes)
	if _, err := jobs.EnqueueHook(c.Request().Context(), h.jobs, jobs.HookPayload{
		HookID: uid.String(),
		Source: body.Source,
	}); err != nil {
		slog.Error("enqueue", "err", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "queue error"})
	}

	return c.JSON(http.StatusAccepted, map[string]any{
		"id":     uid,
		"queued": true,
	})
}
