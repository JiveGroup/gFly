package user

import (
	"gfly/app/domain/models"
	"gfly/app/http/transformers"
	"gfly/app/modules/jwt"
	"github.com/gflydev/core"
)

// NewGetUserProfileApi As a constructor to get user profile API.
func NewGetUserProfileApi() *GetUserProfileApi {
	return &GetUserProfileApi{}
}

// GetUserProfileApi API struct.
type GetUserProfileApi struct {
	core.Api
}

// Handle Process main logic for API.
// @Summary Get user profile
// @Description Get user profile
// @Tags Users
// @Accept json
// @Produce json
// @Success 200 {object} response.User
// @Failure 400 {object} response.Error
// @Security ApiKeyAuth
// @Router /users/profile [get]
func (h *GetUserProfileApi) Handle(c *core.Ctx) error {
	user := c.GetData(jwt.User).(models.User)

	// ==================== Transformer ====================
	// Transform to response data
	var userRes = transformers.ToUserResponse(user)

	// ==================== Response data ====================
	return c.Success(userRes)
}
