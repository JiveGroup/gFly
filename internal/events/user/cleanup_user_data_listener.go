package user

import (
	"github.com/gflydev/core/log"
)

// CleanupUserDataListener removes user-related data (cache, files, sessions)
// after a user account has been deleted.
type CleanupUserDataListener struct{}

// Handle processes the UserDeleted event.
//
// Parameters:
//   - event (events.UserDeleted): The concrete user-deleted event.
//
// Returns:
//   - error: Non-nil if the cleanup encounters a critical failure.
func (l *CleanupUserDataListener) Handle(event UserDeleted) error {
	log.Infof("[Listener] CleanupUserData: cleaning up for user %d (%s)", event.UserID, event.Email)

	return nil
}
