package repository

import (
	"gfly/internal/domain/models"
	mb "github.com/gflydev/db" // Model builder
)

// ====================================================================
// ======================= Repository Interface =======================
// ====================================================================

// IUserRepository represents the interface for managing user-related data in the repository.
// It provides methods for CRUD operations and additional utility functions.
//
// Methods:
//   - GetUserByEmail(email string) *models.User: Retrieves a user by their email address.
//   - GetUserByToken(token string) *models.User: Retrieves a user by their authentication token.
//   - SelectUser(page, limit int) ([]*models.User, int, error): Retrieves a paginated list of users with the total count.
type IUserRepository interface {
	// GetUserByEmail retrieves a user by their email address.
	// Parameters:
	//   - email (string): The email address of the user to retrieve.
	//
	// Returns:
	//   - (*models.User): The user associated with the given email address, or nil if not found.
	GetUserByEmail(email string) *models.User

	// GetUserByToken retrieves a user based on their authentication token.
	// Parameters:
	//   - token (string): The authentication token of the user to retrieve.
	//
	// Returns:
	//   - (*models.User): The user associated with the given token, or nil if not found.
	GetUserByToken(token string) *models.User
}

// ====================================================================
// ====================== Repository Implement ========================
// ====================================================================

// userRepository is a repository type for accessing and managing user data.
type userRepository struct {
}

func (r *userRepository) getBy(field string, value any) *models.User {
	user, err := mb.GetModelBy[models.User](field, value)
	if err != nil {
		return nil
	}

	return user
}

// GetUserByEmail retrieves a user by their email address.
//
// Parameters:
//   - email (string): The email address of the user to retrieve.
//
// Returns:
//   - (*models.User): The user associated with the given email address, or nil if not found.
func (r *userRepository) GetUserByEmail(email string) *models.User {
	return r.getBy("email", email)
}

// GetUserByToken retrieves a user based on their authentication token.
//
// Parameters:
//   - token (string): The authentication token of the user to retrieve.
//
// Returns:
//   - (*models.User): The user associated with the given token, or nil if not found.
func (r *userRepository) GetUserByToken(token string) *models.User {
	return r.getBy("token", token)
}
