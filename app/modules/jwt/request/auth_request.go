package request

import (
	"gfly/app/modules/jwt/dto"
)

type SignIn struct {
	dto.SignIn
}

type SignUp struct {
	dto.SignUp
}

// RefreshToken struct to refresh JWT token.
type RefreshToken struct {
	dto.RefreshToken
}

// ToDto convert to SignIn DTO object.
func (r SignIn) ToDto() dto.SignIn {
	return r.SignIn
}

func (r RefreshToken) ToDto() dto.RefreshToken {
	return r.RefreshToken
}

// ToDto Convert to SignUp DTO object.
func (r SignUp) ToDto() dto.SignUp {
	return r.SignUp
}
