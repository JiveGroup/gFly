package request

import (
	"gfly/internal/modules/auth/dto"
)

type SignIn struct {
	dto.SignIn
}

// ToDto convert to SignIn DTO object.
func (r SignIn) ToDto() dto.SignIn {
	return r.SignIn
}

type SignUp struct {
	dto.SignUp
}

// ToDto Convert to SignUp DTO object.
func (r SignUp) ToDto() dto.SignUp {
	return r.SignUp
}

// RefreshToken struct to refresh JWT token.
type RefreshToken struct {
	dto.RefreshToken
}

func (r RefreshToken) ToDto() dto.RefreshToken {
	return r.RefreshToken
}
