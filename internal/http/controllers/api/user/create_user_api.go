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

type CreateUserApi struct {
	core.Api
}

func NewCreateUserApi() *CreateUserApi {
	return &CreateUserApi{}
}

// ====================================================================
// ======================== Request Validation ========================
// ====================================================================

func (h *CreateUserApi) Validate(c *core.Ctx) error {
	return http.ProcessData[request.CreateUser](c)
}

// ====================================================================
// ========================= Request Handling =========================
// ====================================================================

// Handle function allows Administrator create a new user with specific roles
// @Description Function allows Administrator create a new user with specific roles
// @Summary Create a new user for Administrator
// @Tags Users
// @Accept json
// @Produce json
// @Param data body request.CreateUser true "CreateUser payload"
// @Success 201 {object} response.User
// @Failure 400 {object} http.Error
// @Failure 401 {object} http.Error
// @Security ApiKeyAuth
// @Router /users [post]
func (h *CreateUserApi) Handle(c *core.Ctx) error {
	requestData := c.GetData(constants.Request).(request.CreateUser)

	user, err := services.CreateUser(requestData.ToDto())
	if err != nil {
		return c.Error(http.Error{
			Message: err.Error(),
		})
	}

	// Transform to response data
	userResponse := transformers.ToUserResponse(*user)

	return c.
		Status(core.StatusCreated).
		JSON(userResponse)
}
