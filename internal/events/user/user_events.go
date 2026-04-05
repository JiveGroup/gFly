package user

import (
	"gfly/internal/domain/models"
)

// ---------------------------------------------------------------
//                      Event name constants
// ---------------------------------------------------------------

const (
	// EventUserRegistered fires when a new user has been successfully created.
	EventUserRegistered = "user.registered"

	// EventUserUpdated fires when an existing user's profile data has been modified.
	EventUserUpdated = "user.updated"

	// EventUserDeleted fires when a user has been removed from the system.
	EventUserDeleted = "user.deleted"
)

// ---------------------------------------------------------------
//                        User Events
// ---------------------------------------------------------------

// UserRegistered is dispatched after a new user account is created.
type UserRegistered struct {
	// User is the newly created user model.
	User *models.User
}

// EventName returns the unique event identifier.
func (e UserRegistered) EventName() string { return EventUserRegistered }

// UserUpdated is dispatched after a user's profile has been modified.
type UserUpdated struct {
	// User is the updated user model.
	User *models.User
}

// EventName returns the unique event identifier.
func (e UserUpdated) EventName() string { return EventUserUpdated }

// UserDeleted is dispatched after a user has been deleted from the system.
type UserDeleted struct {
	// UserID is the ID of the deleted user.
	UserID int
	// Email is the email address of the deleted user (for notifications / cleanup).
	Email string
}

// EventName returns the unique event identifier.
func (e UserDeleted) EventName() string { return EventUserDeleted }
