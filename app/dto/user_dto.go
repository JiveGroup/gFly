package dto

import "gfly/app/domain/models/types"

// CreateUser struct to describe the request body to create a new user.
// CreateUser struct to describe the request body to create a new user.
// @Description Request payload for creating a new user.
// @Tags Users
type CreateUser struct {
	Email    string       `json:"email" example:"john@jivecode.com" validate:"required,email,max=255" doc:"User's email address (required, max length 255)"`
	Password string       `json:"password" example:"M1PassW@s" validate:"required,max=255" doc:"User's password (required, max length 255)"`
	Fullname string       `json:"fullname" example:"John Doe" validate:"required,max=255" doc:"User's full name (required, max length 255)"`
	Phone    string       `json:"phone" example:"0989831911" validate:"required,max=20" doc:"User's phone number (required, max length 20)"`
	Avatar   string       `json:"avatar" example:"https://i.pravatar.cc/32" validate:"omitempty,max=255" doc:"URL of the user's avatar (optional, max length 255)"`
	Status   string       `json:"status" example:"pending" validate:"omitempty" doc:"User's status (optional)"`
	Roles    []types.Role `json:"roles" example:"admin,user" validate:"omitempty" doc:"List of user's roles (optional)"`
}

// UpdateUser struct to partially update an existed user.
// @Description Request payload for updating an existing user.
// @Tags Users
type UpdateUser struct {
	ID       int          `json:"-" validate:"omitempty,gte=1" doc:"User ID (greater than or equal to 1)"`
	Password string       `json:"password" example:"M1PassW@s" validate:"omitempty,max=255" doc:"User's new password (optional, max length 255)"`
	Fullname string       `json:"fullname" example:"John Doe" validate:"max=255" doc:"User's updated full name (optional, max length 255)"`
	Phone    string       `json:"phone" example:"0989831911" validate:"max=20" doc:"User's updated phone number (optional, max length 20)"`
	Avatar   string       `json:"avatar" example:"https://i.pravatar.cc/32" validate:"max=255" doc:"Updated URL of the user's avatar (optional, max length 255)"`
	Roles    []types.Role `json:"roles" example:"admin,user" validate:"omitempty" doc:"Updated list of user's roles (optional)"`
}

// UpdateUserStatus struct allows update `status` field from an existing user.
// @Description Request payload for updating the status field of an existing user.
// @Tags Users
type UpdateUserStatus struct {
	ID     int              `json:"-" validate:"omitempty" doc:"User ID associated with the status update"`
	Status types.UserStatus `json:"status" example:"active" validate:"required,oneof=active pending blocked" doc:"New status of the user (required, one of: active, pending, blocked)"`
}
