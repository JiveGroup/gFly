package queues

import (
	"github.com/gflydev/console"
	"github.com/gflydev/core/errors"
	"github.com/gflydev/core/log"
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
func (t HelloTask) Dequeue(task *console.TaskPayload) error {
	// Decode task payload
	var payload HelloTaskPayload
	if err := task.BindPayload(&payload); err != nil {
		return errors.New("json.Unmarshal failed: %v: %s", err, task.GetType())
	}

	// Process payload
	log.Infof("Handle HelloTask with message %s", payload.Message)

	return nil
}
