package dto

// ForgotPassword struct to describe forgot password.
type ForgotPassword struct {
	Username string `json:"username" example:"john@jivecode.com" validate:"required,email,lte=255"`
}

// ResetPassword struct to describe reset password.
type ResetPassword struct {
	Password string `json:"password" example:"M1PassW@s" validate:"required,gte=6"`
	Token    string `json:"token" example:"293r823or832eioj2eo9282o423" validate:"required,lte=255"`
}
