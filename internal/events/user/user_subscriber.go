package user

import (
	"github.com/gflydev/event"
)

// ---------------------------------------------------------------
//                        User Event Subscriber
// ---------------------------------------------------------------

// UserSubscriber groups all listeners for user-domain events.
// Register new user-related listeners here to keep event wiring centralised.
type UserSubscriber struct{}

// Subscribe UserSubscriber registers user event listeners on the given dispatcher.
//
// Registered mappings:
//   - user.registered → SendWelcomeEmailListener, QueuedWelcomeEmailListener
//   - user.deleted    → CleanupUserDataListener
func (s *UserSubscriber) Subscribe(d *event.Dispatcher) {
	event.ListenOn[UserRegistered](d, &SendWelcomeEmailListener{})
	event.ListenOn[UserRegistered](d, &QueuedWelcomeEmailListener{})
	event.ListenOn[UserDeleted](d, &CleanupUserDataListener{})
}
