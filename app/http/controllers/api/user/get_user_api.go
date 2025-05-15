package user

import (
	"gfly/app/constants"
	"gfly/app/domain/models"
	"gfly/app/http"
	"gfly/app/http/response"
	"gfly/app/http/transformers"
	"github.com/gflydev/core"
	"github.com/gflydev/core/log"
	mb "github.com/gflydev/db"
)

type GetUserByIdApi struct {
	core.Api
}

func NewGetUserByIdApi() *GetUserByIdApi {
	return &GetUserByIdApi{}
}

func (h *GetUserByIdApi) Validate(c *core.Ctx) error {
	// Receive path parameter ID
	itemID, errData := http.PathID(c)
	if errData != nil {
		return c.Error(errData)
	}

	// Store data into context.
	c.SetData(constants.Data, itemID)

	return nil
}

// Handle function get user by given id. If user not exists, return not found status.
// @Description Function get user by given id. If user not exists, return not found status.
// @Summary Get user by given id
// @Tags Users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} response.User
// @Failure 401 {object} response.Unauthorized
// @Failure 404 {object} response.Error
// @Security ApiKeyAuth
// @Router /users/{id} [get]
func (h *GetUserByIdApi) Handle(c *core.Ctx) error {
	userID := c.GetData(constants.Data).(int)

	user, err := mb.GetModelByID[models.User](userID)
	if err != nil {
		log.Error(err)

		return c.Error(response.Error{
			Code:    core.StatusNotFound,
			Message: "User not found",
		}, core.StatusNotFound)
	}

	// Transform to response data
	userTransformer := transformers.ToSignUpResponse(*user)

	return c.Success(userTransformer)
}
