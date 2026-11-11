package jobs

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/hibiken/asynq"
)

const TaskHookReceive = "hook:receive"

type HookPayload struct {
	HookID string `json:"hook_id"`
	Source string `json:"source"`
}

func EnqueueHook(ctx context.Context, client *asynq.Client, p HookPayload) (*asynq.TaskInfo, error) {
	b, err := json.Marshal(p)
	if err != nil {
		return nil, fmt.Errorf("marshal: %w", err)
	}
	return client.EnqueueContext(ctx, asynq.NewTask(TaskHookReceive, b))
}

type Worker struct{ log *slog.Logger }

func NewWorker(log *slog.Logger) *Worker { return &Worker{log: log} }

func (w *Worker) Handler() *asynq.ServeMux {
	mux := asynq.NewServeMux()
	mux.HandleFunc(TaskHookReceive, w.handleHookReceive)
	return mux
}

func (w *Worker) handleHookReceive(_ context.Context, t *asynq.Task) error {
	var p HookPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}
	w.log.Info("hook received", "hook_id", p.HookID, "source", p.Source)
	return nil
}
