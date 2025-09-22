package middleware

import (
	"gfly/internal/constants"
	"gfly/internal/domain/models"
	"gfly/internal/domain/models/types"
	"gfly/internal/http/response"
	"gfly/internal/services"
	"github.com/gflydev/core"
	"github.com/gflydev/core/log"
	"slices"
)

// CheckRolesMiddleware is a middleware that verifies if a user has the required roles
// to access a specific route. It ensures that users without proper permissions
// are denied access.
//
// Parameters:
//   - roles ([]types.Role): A list of roles required to access the route.
//   - excludes (...string): Optional paths to exclude from role checks.
//
// Returns:
//   - core.MiddlewareHandler: A middleware handler function.
//
// @Throws response.Error with code 401 when user is not authenticated
// @Throws response.Error with code 403 when user lacks required role permissions
func CheckRolesMiddleware(roles []types.Role, excludes ...string) core.MiddlewareHandler {
	return func(c *core.Ctx) error {
		path := c.Path()

		// Skip role checks for excluded paths
		if slices.Contains(excludes, path) {
			log.Tracef("skip check roles for %v", path)
			return nil
		}

		if c.GetData(constants.User) == nil {
			return c.Error(response.Error{
				Code:    core.StatusUnauthorized,
				Message: "Unauthorized",
			}, core.StatusUnauthorized)
		}

		// Retrieve user data from the context
		user := c.GetData(constants.User).(models.User)

		// Check if user has any of the required roles
		if !services.UserHasRole(user.ID, roles) {
			// Return error response
			return c.Error(response.Error{
				Code:    core.StatusForbidden,
				Message: "Permission denied",
			}, core.StatusForbidden)
		}

		c.Status(core.StatusOK)

		return nil
	}
}
