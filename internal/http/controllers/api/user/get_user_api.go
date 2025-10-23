package user

import (
	"gfly/internal/domain/models"
	_ "gfly/internal/http/response" // Used for Swagger documentation
	"gfly/internal/http/transformers"
	"gfly/pkg/constants"
	"gfly/pkg/http"
	"github.com/gflydev/core"
	"github.com/gflydev/core/log"
	mb "github.com/gflydev/db"
)

// ====================================================================
// ======================== Controller Creation =======================
// ====================================================================

type GetUserByIdApi struct {
	core.Api
}

func NewGetUserByIdApi() *GetUserByIdApi {
	return &GetUserByIdApi{}
}

// ====================================================================
// ======================== Request Validation ========================
// ====================================================================

func (h *GetUserByIdApi) Validate(c *core.Ctx) error {
	return http.ProcessPathID(c)
}

// ====================================================================
// ========================= Request Handling =========================
// ====================================================================

// Handle function get user by given id. If user not exists, return not found status.
// @Description Function get user by given id. If user not exists, return not found status.
// @Summary Get user by given id
// @Tags Users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} response.User
// @Failure 401 {object} http.Error
// @Failure 404 {object} http.Error
// @Security ApiKeyAuth
// @Router /users/{id} [get]
func (h *GetUserByIdApi) Handle(c *core.Ctx) error {
	userID := c.GetData(constants.PathID).(int)

	user, err := mb.GetModelByID[models.User](userID)
	if err != nil {
		log.Error(err)

		return c.Error(http.Error{
			Message: "User not found",
		}, core.StatusNotFound)
	}

	// Transform to response data
	userTransformer := transformers.ToUserResponse(*user)

	return c.Success(userTransformer)
}
