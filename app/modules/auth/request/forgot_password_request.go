package request

import "gfly/app/modules/auth/dto"

// ForgotPassword struct to describe forgot password.
type ForgotPassword struct {
	dto.ForgotPassword
}

// ToDto Convert to ForgotPassword DTO object.
func (r ForgotPassword) ToDto() dto.ForgotPassword {
	return r.ForgotPassword
}
