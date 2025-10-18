package services

import (
	"errors"
	"fmt"
	"gfly/internal/domain/repository"
	"gfly/pkg/modules/auth/dto"
	"gfly/pkg/modules/auth/notifications"
	"github.com/gflydev/core/log"
	"github.com/gflydev/core/utils"
	mb "github.com/gflydev/db"
	dbNull "github.com/gflydev/db/null"
	"github.com/gflydev/notification"
	"time"
)

// ====================================================================
// ========================= Main functions ===========================
// ====================================================================

// ForgotPassword processes a forgot password request by generating a reset
// token and sending a password reset email.
//
// This function handles the forgot password flow by:
// 1. Looking up the user by their email address
// 2. Generating a secure reset token using SHA256 hash of email + timestamp
// 3. Saving the reset token to the user's record
// 4. Sending a password reset email with instructions
//
// Parameters:
//   - forgotPassword: dto.ForgotPassword struct containing the user's email address
//
// Returns:
//   - error: nil if successful, otherwise:
//   - "invalid input data" if user email is not found
//   - "service error" if database update or email notification fails
//
// Example usage:
//
//	err := ForgotPassword(dto.ForgotPassword{
//		Username: "user@example.com"
//	})
func ForgotPassword(forgotPassword dto.ForgotPassword) error {
	// Get user by ID.
	user := repository.Pool.GetUserByEmail(forgotPassword.Username)
	// Item not found error
	if user == nil {
		return errors.New("invalid input data")
	}

	// Create a new SHA256 hash.
	hash := utils.Sha256(user.Email, time.Now().Unix())

	user.Token = dbNull.String(interpolateToken(hash))
	user.UpdatedAt = time.Now()

	if err := mb.UpdateModel(user); err != nil {
		return errors.New("service error")
	}

	// Send notification via mail
	if err := notification.Send(notifications.ResetPassword{
		Email: user.Email,
	}); err != nil {
		log.Errorf("Service forgot password error '%v'", err)

		return errors.New("service error")
	}

	return nil
}

// ChangePassword processes a password reset request and updates the user's password.
//
// This function handles the password reset flow by:
// 1. Verifying the reset token and retrieving associated user
// 2. Clearing the reset token
// 3. Updating the password with a new hashed value
// 4. Sending email notification about password change
//
// Parameters:
//   - resetPassword: dto.ResetPassword struct containing the new password and reset token
//
// Returns:
//   - error: nil if successful, otherwise:
//   - "invalid input data" if token is invalid/expired
//   - "service error" if database update or email notification fails
//
// Example usage:
//
//	err := ChangePassword(dto.ResetPassword{
//		Password: "newpass123",
//		Token: "abc123token",
//	})
func ChangePassword(resetPassword dto.ResetPassword) error {
	// Get user by ID.
	user := repository.Pool.GetUserByToken(interpolateToken(resetPassword.Token))
	// Item not found error
	if user == nil {
		return errors.New("invalid input data")
	}

	user.Token = dbNull.String("")
	user.UpdatedAt = time.Now()
	user.Password = utils.GeneratePassword(resetPassword.Password)

	if err := mb.UpdateModel(user); err != nil {
		log.Errorf("Change password error '%v'", err)

		return errors.New("service error")
	}

	// Send notification via mail
	if err := notification.Send(notifications.ChangePassword{
		Email: user.Email,
		Name:  user.Fullname,
	}); err != nil {
		log.Errorf("Change password error '%v'", err)

		return errors.New("service error")
	}

	return nil
}

// ====================================================================
// ======================== Helper Functions ==========================
// ====================================================================

func interpolateToken(token string) string {
	return fmt.Sprintf("reset_password:%s", token)
}
