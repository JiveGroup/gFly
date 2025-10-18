package user

import (
	"gfly/internal/constants"
	"gfly/internal/http/request"
	"gfly/internal/http/response"
	"gfly/internal/http/transformers"
	"gfly/internal/services"
	"gfly/pkg/http"
	"github.com/gflydev/core"
)

// ====================================================================
// ======================== Controller Creation =======================
// ====================================================================

type UpdateUserApi struct {
	core.Api
}

func NewUpdateUserApi() *UpdateUserApi {
	return &UpdateUserApi{}
}

// ====================================================================
// ======================== Request Validation ========================
// ====================================================================

func (h *UpdateUserApi) Validate(c *core.Ctx) error {
	return http.ProcessUpdateData[request.UpdateUser](c)
}

// ====================================================================
// ========================= Request Handling =========================
// ====================================================================

// Handle function allows Administrator update users table or authorize user roles.
// @Description Function allows Administrator update users table or authorize user roles.
// @Summary Function allows Administrator update an existing user
// @Tags Users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param data body request.UpdateUser true "UpdateUser payload"
// @Success 200 {object} response.User
// @Failure 400 {object} response.Error
// @Failure 401 {object} response.Unauthorized
// @Security ApiKeyAuth
// @Router /users/{id} [put]
func (h *UpdateUserApi) Handle(c *core.Ctx) error {
	requestData := c.GetData(constants.Request).(request.UpdateUser)

	user, err := services.UpdateUser(requestData.ToDto())
	if err != nil {
		return c.Error(response.Error{
			Code:    core.StatusBadRequest,
			Message: err.Error(),
		})
	}

	// Transform to response data
	userTransformer := transformers.ToUserResponse(*user)

	return c.Success(userTransformer)
}
