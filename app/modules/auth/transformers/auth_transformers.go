package transformers

import (
	"gfly/app/modules/auth"
	"gfly/app/modules/auth/response"
)

// ToSignInResponse function JWTTokens struct to SignIn response object.
func ToSignInResponse(tokens *auth.Token) response.SignIn {
	return response.SignIn{
		Access:  tokens.Access,
		Refresh: tokens.Refresh,
	}
}
