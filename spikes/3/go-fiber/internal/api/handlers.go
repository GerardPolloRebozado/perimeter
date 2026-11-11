package api

import (
	"encoding/json"
	"log/slog"

	"github.com/gofiber/fiber/v2"
	"github.com/hibiken/asynq"
	openapi_types "github.com/oapi-codegen/runtime/types"

	"perimeter-volume/spike-3/go-fiber/internal/cfg"
	"perimeter-volume/spike-3/go-fiber/internal/db/sqlc"
	"perimeter-volume/spike-3/go-fiber/internal/jobs"
)

type Handler struct {
	cfg  *cfg.Config
	q    *sqlc.Queries
	jobs *asynq.Client
}

func New(c *cfg.Config, q *sqlc.Queries, j *asynq.Client) *Handler {
	return &Handler{cfg: c, q: q, jobs: j}
}

func (h *Handler) Ping(c *fiber.Ctx) error {
	count, err := h.q.CountHooks(c.Context())
	if err != nil {
		slog.Error("count hooks", "err", err)
		return c.Status(fiber.StatusInternalServerError).JSON(map[string]string{"error": "db error"})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":     "pong",
		"service":    h.cfg.ServiceName,
		"version":    h.cfg.Version,
		"hook_count": count,
	})
}

func (h *Handler) ReceiveHook(c *fiber.Ctx) error {
	var body struct {
		Source  string         `json:"source"`
		Payload map[string]any `json:"payload"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(map[string]string{"error": "invalid json"})
	}
	if body.Source == "" {
		return c.Status(fiber.StatusBadRequest).JSON(map[string]string{"error": "source required"})
	}

	payload := []byte("{}")
	if body.Payload != nil {
		b, err := json.Marshal(body.Payload)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(map[string]string{"error": "marshal payload"})
		}
		payload = b
	}

	hook, err := h.q.InsertHook(c.Context(), payload)
	if err != nil {
		slog.Error("insert hook", "err", err)
		return c.Status(fiber.StatusInternalServerError).JSON(map[string]string{"error": "db error"})
	}

	uid := openapi_types.UUID(hook.ID.Bytes)
	if _, err := jobs.EnqueueHook(c.Context(), h.jobs, jobs.HookPayload{
		HookID: uid.String(),
		Source: body.Source,
	}); err != nil {
		slog.Error("enqueue", "err", err)
		return c.Status(fiber.StatusInternalServerError).JSON(map[string]string{"error": "queue error"})
	}

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"id":     uid,
		"queued": true,
	})
}
