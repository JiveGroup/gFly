package user

import (
	"gfly/app/constants"
	"gfly/app/dto"
	"gfly/app/http"
	"gfly/app/http/request"
	"gfly/app/http/response"
	"gfly/app/http/transformers"
	"gfly/app/services"
	"gfly/app/utils"
	"github.com/gflydev/core"
)

type CreateUserApi struct {
	core.Api
}

func NewCreateUserApi() *CreateUserApi {
	return &CreateUserApi{}
}

func (h *CreateUserApi) Validate(c *core.Ctx) error {
	// Receive request data.
	var requestBody request.CreateUser
	if errData := http.Parse(c, &requestBody); errData != nil {
		return c.Error(errData)
	}

	// Convert to DTO
	requestDto := requestBody.ToDto()

	// Validate DTO
	if errData := utils.Validate(requestDto); errData != nil {
		return c.Error(errData)
	}

	// Store data into context.
	c.SetData(constants.Data, requestDto)

	return nil
}

// Handle function allows Administrator create a new user with specific roles
// @Description Function allows Administrator create a new user with specific roles
// @Summary Create a new user for Administrator
// @Tags Users
// @Accept json
// @Produce json
// @Param data body request.CreateUser true "CreateUser payload"
// @Success 201 {object} response.User
// @Failure 400 {object} response.Error
// @Failure 401 {object} response.Unauthorized
// @Security ApiKeyAuth
// @Router /users [post]
func (h *CreateUserApi) Handle(c *core.Ctx) error {
	createUserDto := c.GetData(constants.Data).(dto.CreateUser)

	user, err := services.CreateUser(createUserDto)
	if err != nil {
		return c.Error(response.Error{
			Code:    core.StatusBadRequest,
			Message: err.Error(),
		})
	}

	// Transform to response data
	userTransformer := transformers.ToUserResponse(*user)

	return c.
		Status(core.StatusCreated).
		JSON(userTransformer)
}
