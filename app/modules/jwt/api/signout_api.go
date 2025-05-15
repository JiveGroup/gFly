package api

import (
	"gfly/app/modules/jwt"
	"github.com/gflydev/core"
	"github.com/gflydev/core/errors"
)

// NewSignOutApi As a constructor to create new API.
func NewSignOutApi() *SignOutApi {
	return &SignOutApi{}
}

type SignOutApi struct {
	core.Api
}

// Handle method to invalidate users access token by adding them to a blacklist in Redis
// and delete refresh token from the Redis
// @Description De-authorize user and delete refresh token from Redis.
// @Summary de-authorize user and delete refresh token from Redis
// @Tags Auth
// @Accept json
// @Produce json
// @Failure 400 {object} response.Error
// @Failure 401 {object} response.Unauthorized
// @Success 204
// @Security ApiKeyAuth
// @Router /auth/signout [delete]
func (h *SignOutApi) Handle(c *core.Ctx) error {
	jwtToken := jwt.ExtractToken(c)

	err := jwt.SignOut(jwtToken)
	if err != nil {
		return c.Error(errors.New("Error %v", err))
	}

	return c.NoContent()
}
