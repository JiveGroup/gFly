package transformers

import (
	"gfly/app/modules/jwt"
	"gfly/app/modules/jwt/response"
)

// ToSignInResponse function JWTTokens struct to SignIn response object.
func ToSignInResponse(tokens *jwt.Tokens) response.SignIn {
	return response.SignIn{
		Access:  tokens.Access,
		Refresh: tokens.Refresh,
	}
}
