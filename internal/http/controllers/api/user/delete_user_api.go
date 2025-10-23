package user

import (
	_ "gfly/internal/http/response" // Used for Swagger documentation
	"gfly/internal/services"
	"github.com/gflydev/core"
	"github.com/gflydev/http"
)

// ====================================================================
// ======================== Controller Creation =======================
// ====================================================================

type DeleteUserApi struct {
	core.Api
}

func NewDeleteUserApi() *DeleteUserApi {
	return &DeleteUserApi{}
}

// ====================================================================
// ======================== Request Validation ========================
// ====================================================================

func (h *DeleteUserApi) Validate(c *core.Ctx) error {
	return http.ProcessPathID(c)
}

// ====================================================================
// ========================= Request Handling =========================
// ====================================================================

// Handle function hard-delete user with its roles by given userID.
// @Description Function hard-delete user with its roles by given userID.
// @Summary Delete user by given userID
// @Tags Users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} response.User
// @Failure 401 {object} http.Error
// @Failure 404 {object} http.Error
// @Security ApiKeyAuth
// @Router /users/{id} [delete]
func (h *DeleteUserApi) Handle(c *core.Ctx) error {
	userId := c.GetData(http.PathIDKey).(int)

	err := services.DeleteUserByID(userId)
	if err != nil {
		return c.Error(http.Error{
			Message: err.Error(),
		}, core.StatusNotFound)
	}

	return c.NoContent()
}
