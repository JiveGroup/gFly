package api

import (
	"gfly/app/constants"
	"gfly/app/http"
	"gfly/app/http/response"
	"gfly/app/modules/auth/dto"
	"gfly/app/modules/auth/request"
	"gfly/app/modules/auth/services"
	"gfly/app/modules/auth/transformers"
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
	return http.ProcessRequest[request.RefreshToken, dto.RefreshToken](c)
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
// @Failure 400 {object} response.Error
// @Failure 401 {object} response.Unauthorized
// @Success 200 {object} response.SignIn
// @Security ApiKeyAuth
// @Router /auth/refresh [put]
func (h *RefreshTokenApi) Handle(c *core.Ctx) error {
	refreshToken := c.GetData(constants.Request).(dto.RefreshToken)
	// Check valid refresh token
	if !services.IsValidRefreshToken(refreshToken.Token) {
		return c.Error(response.Error{
			Code:    core.StatusUnauthorized,
			Message: "Invalid JWT token",
		}, core.StatusUnauthorized)
	}

	jwtToken := services.ExtractToken(c)
	// Refresh new pairs of access token & refresh token
	tokens, err := services.RefreshToken(jwtToken, refreshToken.Token)
	if err != nil {
		return c.Error(response.Error{
			Code:    core.StatusUnauthorized,
			Message: err.Error(),
		}, core.StatusUnauthorized)
	}

	return c.JSON(transformers.ToSignInResponse(tokens))
}
