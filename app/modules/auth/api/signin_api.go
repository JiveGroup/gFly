package api

import (
	"gfly/app/constants"
	"gfly/app/http"
	"gfly/app/http/response"
	"gfly/app/modules/auth"
	"gfly/app/modules/auth/dto"
	"gfly/app/modules/auth/request"
	"gfly/app/modules/auth/services"
	"gfly/app/modules/auth/transformers"
	"github.com/gflydev/core"
)

// ====================================================================
// ======================== Controller Creation =======================
// ====================================================================

type SignInApi struct {
	Type auth.Type
	core.Api
}

// NewSignInApi is a constructor
func NewSignInApi(authType auth.Type) *SignInApi {
	return &SignInApi{
		Type: authType,
	}
}

// ====================================================================
// ======================== Request Validation ========================
// ====================================================================

// Validate data from request
func (h *SignInApi) Validate(c *core.Ctx) error {
	return http.ProcessRequest[request.SignIn, dto.SignIn](c)
}

// ====================================================================
// ========================= Request Handling =========================
// ====================================================================

// Handle func handle sign in user then returns access token and refresh token
// @Description Authenticating user's credentials then return access and refresh token if valid. Otherwise, return an error message.
// @Summary authenticating user's credentials
// @Tags Auth
// @Accept json
// @Produce json
// @Param data body request.SignIn true "Signin payload"
// @Success 200 {object} response.SignIn
// @Failure 400 {object} response.Error
// @Router /auth/signin [post]
func (h *SignInApi) Handle(c *core.Ctx) error {
	// Get valid data from context
	signInDto := c.GetData(constants.Data).(dto.SignIn)

	tokens, err := services.SignIn(&signInDto)
	if err != nil {
		return c.Error(response.Error{
			Code:    core.StatusBadRequest,
			Message: err.Error(),
		})
	}

	if h.Type == auth.TypeWeb {
		c.SetSession(auth.SessionUsername, signInDto.Username)

		return c.NoContent()
	}

	// Transform to response object.
	signInResponse := transformers.ToSignInResponse(tokens)

	return c.JSONAny(signInResponse)
}
