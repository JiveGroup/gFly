package api

import (
	"gfly/internal/constants"
	"gfly/internal/http/response"
	"gfly/pkg/http"
	"gfly/pkg/modules/auth/request"
	"gfly/pkg/modules/auth/services"
	"github.com/gflydev/core"
)

// ====================================================================
// ======================== Controller Creation =======================
// ====================================================================

// NewForgotPWApi As a constructor to get forgot password API.
// Related with NewResetPWApi
func NewForgotPWApi() *ForgotPWApi {
	return &ForgotPWApi{}
}

// ForgotPWApi API struct.
type ForgotPWApi struct {
	core.Api
}

// ====================================================================
// ======================== Request Validation ========================
// ====================================================================

// Validate Verify data from request.
func (h *ForgotPWApi) Validate(c *core.Ctx) error {
	return http.ProcessData[request.ForgotPassword](c)
}

// ====================================================================
// ========================= Request Handling =========================
// ====================================================================

// Handle method to forget password.
// @Summary Forgot password
// @Description Forgot password.
// @Tags Password
// @Accept json
// @Produce json
// @Param data body request.ForgotPassword true "Forgot password payload"
// @Success 204
// @Failure 400 {object} response.Error
// @Router /password/forgot [post]
func (h *ForgotPWApi) Handle(c *core.Ctx) error {
	requestData := c.GetData(constants.Request).(request.ForgotPassword)

	err := services.ForgotPassword(requestData.ToDto())
	if err != nil {
		return c.Error(response.Error{
			Message: err.Error(),
		})
	}

	return c.NoContent()
}
