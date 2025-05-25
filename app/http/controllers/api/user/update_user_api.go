package user

import (
	"gfly/app/constants"
	"gfly/app/dto"
	"gfly/app/http"
	"gfly/app/http/request"
	"gfly/app/http/response"
	"gfly/app/http/transformers"
	"gfly/app/services"
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
	return http.ProcessUpdateRequest[request.UpdateUser, dto.UpdateUser](c)
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
	updateUserDto := c.GetData(constants.Request).(dto.UpdateUser)

	user, err := services.UpdateUser(updateUserDto)
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
