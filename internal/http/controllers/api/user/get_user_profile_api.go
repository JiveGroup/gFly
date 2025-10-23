package user

import (
	"gfly/internal/domain/models"
	_ "gfly/internal/http/response" // Used for Swagger documentation
	"gfly/internal/http/transformers"
	"github.com/gflydev/core"
	"github.com/gflydev/http"
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
// @Failure 400 {object} http.Error
// @Security ApiKeyAuth
// @Router /users/profile [get]
func (h *GetUserProfileApi) Handle(c *core.Ctx) error {
	if c.GetData(http.UserKey) == nil {
		return c.Error(http.Error{
			Message: "Unauthorized",
		}, core.StatusUnauthorized)
	}

	user := c.GetData(http.UserKey).(models.User)

	// Transform to response data
	var userRes = transformers.ToUserResponse(user)

	return c.Success(userRes)
}
