package api

import (
	"gfly/app/constants"
	"gfly/app/http"
	"gfly/app/http/response"
	"gfly/app/http/transformers"
	"gfly/app/modules/auth/dto"
	"gfly/app/modules/auth/request"
	"gfly/app/modules/auth/services"
	"github.com/gflydev/core"
)

// ====================================================================
// ======================== Controller Creation =======================
// ====================================================================

type SignUp struct {
	core.Api
}

func NewSignUpApi() *SignUp {
	return &SignUp{}
}

// ====================================================================
// ======================== Request Validation ========================
// ====================================================================

func (h *SignUp) Validate(c *core.Ctx) error {
	return http.ProcessRequest[request.SignUp, dto.SignUp](c)
}

// ====================================================================
// ========================= Request Handling =========================
// ====================================================================

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
	signUpDto := c.GetData(constants.Data).(dto.SignUp)

	user, err := services.SignUp(&signUpDto)
	if err != nil {
		return c.Error(response.Error{
			Code:    core.StatusBadRequest,
			Message: err.Error(),
		})
	}

	return c.JSON(transformers.ToSignUpResponse(*user))
}
