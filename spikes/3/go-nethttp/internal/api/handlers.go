package api

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/hibiken/asynq"
	openapi_types "github.com/oapi-codegen/runtime/types"

	"surface-volume/spike-3/go-nethttp/internal/cfg"
	"surface-volume/spike-3/go-nethttp/internal/db/sqlc"
	"surface-volume/spike-3/go-nethttp/internal/jobs"
)

type Handler struct {
	cfg  *cfg.Config
	q    *sqlc.Queries
	jobs *asynq.Client
}

func New(c *cfg.Config, q *sqlc.Queries, j *asynq.Client) *Handler {
	return &Handler{cfg: c, q: q, jobs: j}
}

func (h *Handler) Ping(w http.ResponseWriter, r *http.Request) {
	count, err := h.q.CountHooks(r.Context())
	if err != nil {
		slog.Error("count hooks", "err", err)
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"status":    "pong",
		"service":   h.cfg.ServiceName,
		"version":   h.cfg.Version,
		"hook_count": count,
	})
}

func (h *Handler) ReceiveHook(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Source  string         `json:"source"`
		Payload map[string]any `json:"payload"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if body.Source == "" {
		http.Error(w, "source required", http.StatusBadRequest)
		return
	}

	payload := []byte("{}")
	if body.Payload != nil {
		b, err := json.Marshal(body.Payload)
		if err != nil {
			http.Error(w, "marshal payload", http.StatusBadRequest)
			return
		}
		payload = b
	}

	hook, err := h.q.InsertHook(r.Context(), payload)
	if err != nil {
		slog.Error("insert hook", "err", err)
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	uid := openapi_types.UUID(hook.ID.Bytes)
	if _, err := jobs.EnqueueHook(r.Context(), h.jobs, jobs.HookPayload{
		HookID: uid.String(),
		Source: body.Source,
	}); err != nil {
		slog.Error("enqueue", "err", err)
		http.Error(w, "queue error", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusAccepted, map[string]any{
		"id":     uid,
		"queued": true,
	})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
