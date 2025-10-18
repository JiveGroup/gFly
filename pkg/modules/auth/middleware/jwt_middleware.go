package middleware

import (
	"gfly/internal/constants"
	"gfly/internal/domain/models"
	"gfly/internal/http/response"
	"gfly/pkg/modules/auth/services"
	"github.com/gflydev/core"
	"github.com/gflydev/core/log"
	mb "github.com/gflydev/db"
	"slices"
	"time"
)

// JWTAuth an HTTP middleware that process login via JWT token.
//
// Use:
//
//	app.Use(middleware.JWTAuth(
//		prefixAPI+"/info",
//		prefixAPI+"/auth/signin",
//		prefixAPI+"/auth/refresh",
//	))
func JWTAuth(excludes ...string) core.MiddlewareHandler {
	return func(c *core.Ctx) error {
		path := c.Path()
		if slices.Contains(excludes, path) {
			log.Tracef("Skip JWTAuth checking for '%v'", path)

			return nil
		}

		// Forge status code 401 (Unauthorized) instead of 500 (internal error)
		c.Status(core.StatusUnauthorized)

		jwtToken := services.ExtractToken(c)
		isBlocked, err := services.IsBlockedToken(jwtToken)
		if err != nil {
			log.Errorf("Check JWT error '%v'", err)

			return c.Error(response.Error{
				Code:    core.StatusUnauthorized,
				Message: "Invalid JWT token",
			}, core.StatusUnauthorized)
		}

		if isBlocked {
			return c.Error(response.Error{
				Code:    core.StatusUnauthorized,
				Message: "JWT token was blocked",
			}, core.StatusUnauthorized)
		}

		// Get claims from JWT.
		claims, err := services.ExtractTokenMetadata(jwtToken)
		if err != nil {
			log.Errorf("Parse JWT error '%v'", err)

			return c.Error(response.Error{
				Code:    core.StatusUnauthorized,
				Message: "Parse JWT error",
			}, core.StatusUnauthorized)
		}

		if claims.Expires < time.Now().Unix() {
			log.Errorf("JWT token expired '%v'", jwtToken)

			return c.Error(response.Error{
				Code:    core.StatusUnauthorized,
				Message: "JWT token expired",
			}, core.StatusUnauthorized)
		}

		// Get user by ID.
		user, err := mb.GetModelByID[models.User](claims.UserID)
		if err != nil || user == nil {
			log.Errorf("User not found '%v'", err)

			return c.Error(response.Error{
				Code:    core.StatusUnauthorized,
				Message: "User not found",
			}, core.StatusUnauthorized)
		}

		c.Status(core.StatusOK)
		c.SetData(constants.User, *user)

		return nil
	}
}
