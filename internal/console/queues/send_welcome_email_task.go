package queues

import (
	"gfly/internal/notifications"

	"github.com/gflydev/console"
	"github.com/gflydev/core/errors"
	"github.com/gflydev/core/log"
	"github.com/gflydev/notification"
)

// ---------------------------------------------------------------
//                        Register task.
// ---------------------------------------------------------------

// Auto-register task into queue.
func init() {
	console.RegisterTask(&SendWelcomeEmailTask{}, "send-welcome-email")
}

// ---------------------------------------------------------------
//                        Task info.
// ---------------------------------------------------------------

// NewSendWelcomeEmailTask creates a new queued task payload for sending a welcome email.
//
// Parameters:
//   - email (string): The recipient's email address.
//   - fullname (string): The recipient's full name.
//
// Returns:
//   - (SendWelcomeEmailPayload, string): The task payload and the registered task name.
func NewSendWelcomeEmailTask(email, fullname string) (SendWelcomeEmailPayload, string) {
	return SendWelcomeEmailPayload{
		Email:    email,
		Fullname: fullname,
	}, "send-welcome-email"
}

// SendWelcomeEmailPayload holds the data required to send a welcome email.
type SendWelcomeEmailPayload struct {
	Email    string `json:"email"`
	Fullname string `json:"fullname"`
}

// SendWelcomeEmailTask processes the send-welcome-email queue task.
type SendWelcomeEmailTask struct {
	console.Task
}

// Dequeue handles the queued welcome email task.
//
// Parameters:
//   - task (*console.TaskPayload): The task payload from the queue.
//
// Returns:
//   - error: Non-nil if the task fails to process.
func (t SendWelcomeEmailTask) Dequeue(task *console.TaskPayload) error {
	var payload SendWelcomeEmailPayload
	if err := task.BindPayload(&payload); err != nil {
		return errors.New("SendWelcomeEmailTask: failed to bind payload: %v", err)
	}

	log.Infof("[Queue] SendWelcomeEmail: sending to %s (%s)", payload.Email, payload.Fullname)

	_ = notification.Send(notifications.SendMail{
		Email: payload.Email,
	})

	return nil
}
