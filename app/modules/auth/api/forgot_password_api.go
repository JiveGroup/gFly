package api

import (
	"gfly/app/constants"
	"gfly/app/http"
	"gfly/app/http/response"
	"gfly/app/modules/auth/dto"
	"gfly/app/modules/auth/request"
	"gfly/app/modules/auth/services"
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
	return http.ProcessRequest[request.ForgotPassword, dto.ForgotPassword](c)
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
	data := c.GetData(constants.Data).(dto.ForgotPassword)

	err := services.ForgotPassword(data)
	if err != nil {
		return c.Error(response.Error{
			Message: err.Error(),
			Code:    core.StatusBadRequest,
		})
	}

	return c.NoContent()
}
