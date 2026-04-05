package user

import (
	"gfly/internal/console/queues"
	"github.com/gflydev/console"
	"github.com/gflydev/core/log"
)

// QueuedWelcomeEmailListener defers the welcome email to the queue worker
// instead of processing it inline. This keeps the HTTP request fast and
// delegates the heavy work to the background queue.
//
// Requires the queue worker to be running: ./build/artisan queue:run
type QueuedWelcomeEmailListener struct{}

// Handle processes the UserRegistered event by dispatching a queue task.
//
// Parameters:
//   - event (events.UserRegistered): The concrete user-registered event.
//
// Returns:
//   - error: Non-nil if task dispatch fails.
func (l *QueuedWelcomeEmailListener) Handle(event UserRegistered) error {
	log.Infof("[Listener] QueuedWelcomeEmail: queuing for %s", event.User.Email)

	console.DispatchTask(queues.NewSendWelcomeEmailTask(event.User.Email, event.User.Fullname))

	return nil
}
