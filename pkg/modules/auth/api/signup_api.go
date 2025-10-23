package api

import (
	"gfly/pkg/constants"
	"gfly/pkg/http"
	"gfly/pkg/modules/auth/request"
	_ "gfly/pkg/modules/auth/response" // Used for Swagger documentation
	"gfly/pkg/modules/auth/services"
	"gfly/pkg/modules/auth/transformers"
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
	return http.ProcessData[request.SignUp](c)
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
// @Failure 400 {object} http.Error
// @Success 200 {object} response.SignUp
// @Router /auth/signup [post]
func (h *SignUp) Handle(c *core.Ctx) error {
	requestData := c.GetData(constants.Request).(request.SignUp)

	user, err := services.SignUp(requestData.ToDto())
	if err != nil {
		return c.Error(http.Error{
			Message: err.Error(),
		})
	}

	return c.JSON(transformers.ToSignUpResponse(*user))
}
