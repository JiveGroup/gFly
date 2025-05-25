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

// NewResetPWApi As a constructor to get reset password API.
// Related with NewForgotPWApi
func NewResetPWApi() *ResetPWApi {
	return &ResetPWApi{}
}

// ResetPWApi API struct.
type ResetPWApi struct {
	core.Api
}

// ====================================================================
// ======================== Request Validation ========================
// ====================================================================

// Validate Verify data from request.
func (h *ResetPWApi) Validate(c *core.Ctx) error {
	return http.ProcessRequest[request.ResetPassword, dto.ResetPassword](c)
}

// ====================================================================
// ========================= Request Handling =========================
// ====================================================================

// Handle method to reset password.
// @Summary Reset password
// @Description Reset password.
// @Tags Password
// @Accept json
// @Produce json
// @Param data body request.ResetPassword true "Reset password payload"
// @Success 204
// @Failure 400 {object} response.Error
// @Router /password/reset [post]
func (h *ResetPWApi) Handle(c *core.Ctx) error {
	data := c.GetData(constants.Request).(dto.ResetPassword)

	err := services.ChangePassword(data)
	if err != nil {
		return c.Error(response.Error{
			Message: err.Error(),
			Code:    core.StatusBadRequest,
		})
	}

	return c.NoContent()
}
