package api

import (
	"gfly/pkg/modules/auth/request"
	"gfly/pkg/modules/auth/services"
	"github.com/gflydev/core"
	"github.com/gflydev/http"
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
	return http.ProcessData[request.ResetPassword](c)
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
// @Failure 400 {object} http.Error
// @Router /password/reset [post]
func (h *ResetPWApi) Handle(c *core.Ctx) error {
	requestData := c.GetData(http.RequestKey).(request.ResetPassword)

	err := services.ChangePassword(requestData.ToDto())
	if err != nil {
		return c.Error(http.Error{
			Message: err.Error(),
		})
	}

	return c.NoContent()
}
