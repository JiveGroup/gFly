package user

import (
	"gfly/internal/http/request"
	_ "gfly/internal/http/response" // Used for Swagger documentation
	"gfly/internal/http/transformers"
	"gfly/internal/services"
	"gfly/pkg/constants"
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
// @Failure 400 {object} http.Error
// @Failure 401 {object} http.Error
// @Security ApiKeyAuth
// @Router /users/{id} [put]
func (h *UpdateUserApi) Handle(c *core.Ctx) error {
	requestData := c.GetData(constants.Request).(request.UpdateUser)

	user, err := services.UpdateUser(requestData.ToDto())
	if err != nil {
		return c.Error(http.Error{
			Message: err.Error(),
		})
	}

	// Transform to response data
	userTransformer := transformers.ToUserResponse(*user)

	return c.Success(userTransformer)
}
