package api

import (
	"gfly/internal/constants"
	httpResponse "gfly/internal/http/response"
	"gfly/pkg/http"
	"gfly/pkg/modules/auth/request"
	_ "gfly/pkg/modules/auth/response" // Used for Swagger documentation
	"gfly/pkg/modules/auth/services"
	"gfly/pkg/modules/auth/transformers"
	"github.com/gflydev/core"
)

// ====================================================================
// ======================== Controller Creation =======================
// ====================================================================

// NewRefreshTokenApi As a constructor to create new API.
func NewRefreshTokenApi() *RefreshTokenApi {
	return &RefreshTokenApi{}
}

type RefreshTokenApi struct {
	core.Api
}

// ====================================================================
// ======================== Request Validation ========================
// ====================================================================

// Validate validates request refresh token
func (h *RefreshTokenApi) Validate(c *core.Ctx) error {
	return http.ProcessData[request.RefreshToken](c)
}

// ====================================================================
// ========================= Request Handling =========================
// ====================================================================

// Handle method to refresh user token.
// @Description Refresh user token
// @Summary refresh user token
// @Tags Auth
// @Accept json
// @Produce json
// @Param data body request.RefreshToken true "RefreshToken payload"
// @Failure 400 {object} httpResponse.Error
// @Failure 401 {object} httpResponse.Error
// @Success 200 {object} response.SignIn
// @Security ApiKeyAuth
// @Router /auth/refresh [put]
func (h *RefreshTokenApi) Handle(c *core.Ctx) error {
	requestData := c.GetData(constants.Request).(request.RefreshToken)

	// Check valid refresh token
	if !services.IsValidRefreshToken(requestData.ToDto().Token) {
		return c.Error(httpResponse.Error{
			Message: "Invalid JWT token",
		}, core.StatusUnauthorized)
	}

	jwtToken := services.ExtractToken(c)
	// Refresh new pairs of access token & refresh token
	tokens, err := services.RefreshToken(jwtToken, requestData.ToDto().Token)
	if err != nil {
		return c.Error(httpResponse.Error{
			Message: err.Error(),
		}, core.StatusUnauthorized)
	}

	return c.JSON(transformers.ToSignInResponse(tokens))
}
