package api

import (
	"gfly/app/http/transformers"
	"gfly/app/modules/jwt"
	"gfly/app/modules/jwt/dto"
	"gfly/app/modules/jwt/request"
	"github.com/gflydev/core"
	"github.com/gflydev/core/errors"
	"github.com/gflydev/validation"
)

type SignUp struct {
	core.Api
}

func NewSignUpApi() *SignUp {
	return &SignUp{}
}

func (h *SignUp) Validate(c *core.Ctx) error {
	var signUp request.SignUp
	err := c.ParseBody(&signUp)
	if err != nil {
		c.Status(core.StatusBadRequest)
		return err
	}

	signUpDto := signUp.ToDto()

	errorData, err := validation.Check(signUpDto)
	if err != nil {
		_ = c.Error(errorData)
		return err
	}

	c.SetData(data, signUpDto)

	return nil
}

// Handle function handle sign up user includes create user, create user's role.
// @Description Create a new user with `request.SignUp` body then add `role id` to table `user_roles` with current `user id`
// @Summary Sign up a new user
// @Tags Auth
// @Accept json
// @Produce json
// @Param data body request.SignUp true "Signup payload"
// @Failure 400 {object} response.Error
// @Success 200 {object} response.User
// @Router /auth/signup [post]
func (h *SignUp) Handle(c *core.Ctx) error {
	signUpDto := c.GetData(data).(dto.SignUp)

	user, err := jwt.SignUp(&signUpDto)
	if err != nil {
		return c.Error(errors.New("Error %v", err))
	}

	return c.JSON(transformers.ToSignUpResponse(*user))
}
