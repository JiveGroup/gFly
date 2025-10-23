package transformers

import (
	"gfly/internal/domain/models"
	"gfly/pkg/modules/auth"
	"gfly/pkg/modules/auth/response"
)

// ToSignInResponse function JWTTokens struct to SignIn response object.
func ToSignInResponse(tokens *auth.Token) response.SignIn {
	return response.SignIn{
		Access:  tokens.Access,
		Refresh: tokens.Refresh,
	}
}

// ToSignUpResponse converts a User model to a User response object
// with all fields populated for signup response
//
// Parameters:
//   - user: models.User - The user model to convert
//
// Returns:
//   - response.User: The converted user response object
func ToSignUpResponse(user models.User) response.SignUp {
	return response.SignUp{
		ID:        user.ID,
		Email:     user.Email,
		Fullname:  user.Fullname,
		Phone:     user.Phone,
		Status:    user.Status,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
