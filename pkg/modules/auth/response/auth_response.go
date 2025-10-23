package response

import (
	"gfly/internal/domain/models/types"
	"time"
)

// SignIn struct to describe login response.
type SignIn struct {
	Access  string `json:"access" doc:"The access token for authentication"`
	Refresh string `json:"refresh" doc:"The refresh token for obtaining a new access token"`
}

// SignUp struct to describe User response.
// The instance should be created from models.User.ToResponse()
type SignUp struct {
	ID        int              `json:"id" doc:"The unique identifier for the user."`
	Email     string           `json:"email" doc:"The email address of the user."`
	Fullname  string           `json:"fullname" doc:"The full name of the user."`
	Phone     string           `json:"phone" doc:"The phone number of the user."`
	Token     *string          `json:"token" doc:"The authorization token of the user."`
	Status    types.UserStatus `json:"status" doc:"The status of the user account."`
	CreatedAt time.Time        `json:"created_at" doc:"The timestamp of when the user was created."`
	UpdatedAt time.Time        `json:"updated_at" doc:"The timestamp of when the user was last updated."`
}
