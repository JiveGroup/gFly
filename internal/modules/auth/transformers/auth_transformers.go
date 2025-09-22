package transformers

import (
	"gfly/internal/modules/auth"
	"gfly/internal/modules/auth/response"
)

// ToSignInResponse function JWTTokens struct to SignIn response object.
func ToSignInResponse(tokens *auth.Token) response.SignIn {
	return response.SignIn{
		Access:  tokens.Access,
		Refresh: tokens.Refresh,
	}
}
