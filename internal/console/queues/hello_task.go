package queues

import (
	"context"
	"encoding/json"
	"github.com/gflydev/console"
	"github.com/gflydev/core/errors"
	"github.com/gflydev/core/log"
	"github.com/hibiken/asynq"
)

// ---------------------------------------------------------------
// 					Register task.
// ---------------------------------------------------------------

// Auto-register task into queue.
func init() {
	console.RegisterTask(&HelloTask{}, "hello-world")
}

// ---------------------------------------------------------------
// 					Task info.
// ---------------------------------------------------------------

// NewHelloTask Constructor HelloTask.
func NewHelloTask(message string) (HelloTaskPayload, string) {
	return HelloTaskPayload{
		Message: message,
	}, "hello-world"
}

// HelloTaskPayload Task payload.
type HelloTaskPayload struct {
	Message string
}

// HelloTask Hello task.
type HelloTask struct {
	console.Task
}

// Dequeue Handle a task in queue.
func (t HelloTask) Dequeue(ctx context.Context, task *asynq.Task) error {
	// Decode task payload
	var payload HelloTaskPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return errors.New("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	// Process payload
	log.Infof("Handle HelloTask with message %s", payload.Message)

	return nil
}
