package dto

// SignUp struct to describe register a new user.
type SignUp struct {
	Email    string `json:"email" example:"john@jivecode.com" validate:"required,email,max=255" doc:"The email address of the user, must be a valid email address and is required"`
	Password string `json:"password" example:"M1PassW@s" validate:"required,max=255" doc:"The password for the new user account, up to 255 characters and is required"` // #nosec G117 -- Legitimate password input field for user registration
	Fullname string `json:"fullname" example:"John Doe" validate:"required,max=255" doc:"The full name of the user, up to 255 characters and is required"`
	Phone    string `json:"phone" example:"0989831911" validate:"required,max=20" doc:"The phone number of the user, up to 20 characters and is required"`
	Avatar   string `json:"avatar" example:"https://i.pravatar.cc/32" validate:"max=255" doc:"The avatar URL for the user, up to 255 characters and optional"`
	Status   string `json:"status" example:"pending" validate:"omitempty" doc:"The status of the user, optional field"`
}

// SignIn struct to describe sign in user
type SignIn struct {
	Username string `json:"username" example:"admin@gfly.dev" validate:"required,email,max=255" doc:"The email address or username used for signing in, must be a valid email and is required"`
	Password string `json:"password" example:"P@seWor9" validate:"required,max=255" doc:"The password for the user account, up to 255 characters and is required"` // #nosec G117 -- Legitimate password input field for user authentication
}

// RefreshToken struct to refresh JWT token.
type RefreshToken struct {
	Token string `json:"token" example:"d1a4216a226cbf75eaefc9107c2c64b6b2c0f18cd8634e3a6f495146c38e1324.1747914602" validate:"required,max=255" doc:"The refresh token for obtaining a new access token, up to 255 characters and is required"`
}
