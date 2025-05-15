package middleware

import (
	"gfly/app/constants"
	"gfly/app/domain/models"
	"gfly/app/http/response"
	"github.com/gflydev/core"
	"strconv"
)

// PreventUpdateYourSelf is a middleware that prevents a user from updating their own data.
// It extracts the user ID from the request path and compares it with the current authenticated user’s ID.
// If the user attempts to update their own record, it returns a 403 Forbidden status.
//
// @Param id path string true "User ID to update"
// @Success 200 {object} nil "Success"
// @Failure 403 {object} response.Error "Forbidden: User cannot update their own record"
// @Failure 400 {object} response.Error "Bad Request: Invalid User ID"
func PreventUpdateYourSelf(c *core.Ctx) error {
	// Parse request parameter
	userId, err := strconv.Atoi(c.PathVal("id"))
	if err != nil {
		return err
	}

	if c.GetData(constants.User) == nil {
		return c.Error(response.Error{
			Code:    core.StatusUnauthorized,
			Message: "Unauthorized",
		}, core.StatusUnauthorized)
	}

	// Get user data
	user := c.GetData(constants.User).(models.User)

	// In-case the user updates itself
	if userId == user.ID {
		// Force update request header
		c.Root().Request.Header.SetContentType(core.MIMEApplicationJSONCharsetUTF8)

		return c.Error(response.Error{
			Code:    core.StatusForbidden,
			Message: "Don't allow update yourself",
		}, core.StatusForbidden)
	}

	c.Status(core.StatusOK)

	return nil
}
