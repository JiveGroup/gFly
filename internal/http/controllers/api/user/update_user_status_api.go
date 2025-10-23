package user

import (
	"gfly/internal/http/request"
	_ "gfly/internal/http/response" // Used for Swagger documentation
	"gfly/internal/http/transformers"
	"gfly/internal/services"
	"github.com/gflydev/core"
	"github.com/gflydev/http"
)

// ====================================================================
// ======================== Controller Creation =======================
// ====================================================================

type UpdateUserStatusApi struct {
	core.Api
}

func NewUpdateUserStatusApi() *UpdateUserStatusApi {
	return &UpdateUserStatusApi{}
}

// ====================================================================
// ======================== Request Validation ========================
// ====================================================================

func (h UpdateUserStatusApi) Validate(c *core.Ctx) error {
	return http.ProcessUpdateData[request.UpdateUserStatus](c)
}

// ====================================================================
// ========================= Request Handling =========================
// ====================================================================

// Handle Process main logic for API.
// @Summary Update user's status by ID
// @Description Update user's status by ID. <b>Administrator privilege required</b>
// @Tags Users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param request body request.UpdateUserStatus true "Update user status data"
// @Failure 400 {object} http.Error
// @Failure 401 {object} http.Error
// @Success 200 {object} response.User
// @Security ApiKeyAuth
// @Router /users/{id}/status [put]
func (h UpdateUserStatusApi) Handle(c *core.Ctx) error {
	requestData := c.GetData(http.RequestKey).(request.UpdateUserStatus)

	// Bind data to service
	user, err := services.UpdateUserStatus(requestData.ToDto())
	if err != nil {
		c.Status(core.StatusBadRequest)
		return err
	}

	// Transform response data
	userResponse := transformers.ToUserResponse(*user)

	return c.Success(userResponse)
}
