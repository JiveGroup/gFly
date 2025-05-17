package user

import (
	"gfly/app/constants"
	"gfly/app/domain/models"
	"gfly/app/http/response"
	"gfly/app/http/transformers"
	"github.com/gflydev/core"
)

// ====================================================================
// ======================== Controller Creation =======================
// ====================================================================

// NewGetUserProfileApi As a constructor to get user profile API.
func NewGetUserProfileApi() *GetUserProfileApi {
	return &GetUserProfileApi{}
}

// GetUserProfileApi API struct.
type GetUserProfileApi struct {
	core.Api
}

// ====================================================================
// ========================= Request Handling =========================
// ====================================================================

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
	if c.GetData(constants.User) == nil {
		return c.Error(response.Error{
			Code:    core.StatusUnauthorized,
			Message: "Unauthorized",
		}, core.StatusUnauthorized)
	}

	user := c.GetData(constants.User).(models.User)

	// Transform to response data
	var userRes = transformers.ToUserResponse(user)

	return c.Success(userRes)
}
