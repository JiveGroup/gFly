package user

import (
	"gfly/app/constants"
	"gfly/app/http"
	"gfly/app/http/response"
	"gfly/app/services"
	"github.com/gflydev/core"
)

type DeleteUserApi struct {
	core.Api
}

func NewDeleteUserApi() *DeleteUserApi {
	return &DeleteUserApi{}
}

func (h *DeleteUserApi) Validate(c *core.Ctx) error {
	// Receive path parameter ID
	itemID, errData := http.PathID(c)
	if errData != nil {
		return c.Error(errData)
	}

	// Store data into context.
	c.SetData(constants.Data, itemID)

	return nil
}

// Handle function hard-delete user with its roles by given userID.
// @Description Function hard-delete user with its roles by given userID.
// @Summary Delete user by given userID
// @Tags Users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} response.User
// @Failure 401 {object} response.Unauthorized
// @Failure 404 {object} response.Error
// @Security ApiKeyAuth
// @Router /users/{id} [delete]
func (h *DeleteUserApi) Handle(c *core.Ctx) error {
	userId := c.GetData(constants.Data).(int)

	err := services.DeleteUserByID(userId)
	if err != nil {
		return c.Error(response.Error{
			Code:    core.StatusNotFound,
			Message: err.Error(),
		}, core.StatusNotFound)
	}

	return c.NoContent()
}
