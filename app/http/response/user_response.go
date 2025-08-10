package response

import (
	"gfly/app/domain/models/types"
	"gfly/app/dto"
	"time"
)

// User struct to describe User response.
// The instance should be created from models.User.ToResponse()
type User struct {
	ID           int              `json:"id" doc:"The unique identifier for the user."`
	Email        string           `json:"email" doc:"The email address of the user."`
	Fullname     string           `json:"fullname" doc:"The full name of the user."`
	Phone        string           `json:"phone" doc:"The phone number of the user."`
	Token        *string          `json:"token" doc:"The authorization token of the user."`
	Status       types.UserStatus `json:"status" doc:"The status of the user account."`
	CreatedAt    time.Time        `json:"created_at" doc:"The timestamp of when the user was created."`
	UpdatedAt    time.Time        `json:"updated_at" doc:"The timestamp of when the user was last updated."`
	VerifiedAt   *time.Time       `json:"verified_at" example:"2023-01-01T10:30:00Z" doc:"The timestamp of when the user was verified."`
	BlockedAt    *time.Time       `json:"blocked_at" example:"null" doc:"The timestamp of when the user was blocked."`
	DeletedAt    *time.Time       `json:"deleted_at" example:"null" doc:"The timestamp of when the user was deleted."`
	LastAccessAt *time.Time       `json:"last_access_at" example:"2023-01-01T12:00:00Z" doc:"The timestamp of the user's last access."`
	Avatar       *string          `json:"avatar" doc:"The URL of the user's avatar or profile picture."`
	Roles        []Role           `json:"roles" doc:"A list of roles assigned to the user."`
}

// Role struct to describe Role response.
type Role struct {
	ID   int        `json:"id" doc:"The unique identifier for the role."`
	Name string     `json:"name" doc:"The name of the role."`
	Slug types.Role `json:"slug" doc:"The slug (URL-friendly name) of the role."`
}

type ListUser struct {
	Meta dto.Meta `json:"meta" doc:"Pagination metadata for a list of users."`
	Data []User   `json:"data" doc:"A list of users matching the query criteria."`
}
