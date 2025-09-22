package request

import "gfly/internal/modules/auth/dto"

// ResetPassword struct to describe reset password.
type ResetPassword struct {
	dto.ResetPassword
}

// ToDto Convert to ForgotPassword DTO object.
func (r ResetPassword) ToDto() dto.ResetPassword {
	return r.ResetPassword
}
