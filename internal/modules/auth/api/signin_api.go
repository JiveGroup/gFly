package api

import (
	"gfly/internal/constants"
	"gfly/internal/http"
	httpResponse "gfly/internal/http/response"
	"gfly/internal/modules/auth"
	"gfly/internal/modules/auth/dto"
	"gfly/internal/modules/auth/request"
	_ "gfly/internal/modules/auth/response" // Used for Swagger documentation
	"gfly/internal/modules/auth/services"
	"gfly/internal/modules/auth/transformers"
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
// @Failure 400 {object} httpResponse.Error
// @Router /auth/signin [post]
func (h *SignInApi) Handle(c *core.Ctx) error {
	// Get valid data from context
	signInDto := c.GetData(constants.Request).(dto.SignIn)

	tokens, err := services.SignIn(signInDto)
	if err != nil {
		return c.Error(httpResponse.Error{
			Code:    core.StatusBadRequest,
			Message: err.Error(),
		})
	}

	if h.Type == auth.TypeWeb {
		c.SetSession(auth.SessionUsername, signInDto.Username)

		return c.NoContent()
	}

	return c.JSON(transformers.ToSignInResponse(tokens))
}
