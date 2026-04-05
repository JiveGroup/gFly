package user

import (
	"gfly/internal/notifications"

	"github.com/gflydev/core/log"
	"github.com/gflydev/notification"
)

// SendWelcomeEmailListener sends a welcome email when a new user registers.
// This listener runs synchronously in the same goroutine as the dispatcher.
type SendWelcomeEmailListener struct{}

// Handle processes the UserRegistered event.
//
// Parameters:
//   - event (events.UserRegistered): The concrete user-registered event.
//
// Returns:
//   - error: Non-nil if the listener encounters a critical failure.
func (l *SendWelcomeEmailListener) Handle(event UserRegistered) error {
	log.Infof("[Listener] SendWelcomeEmail: sending to %s", event.User.Email)

	_ = notification.Send(notifications.SendMail{
		Email: event.User.Email,
	})

	return nil
}
