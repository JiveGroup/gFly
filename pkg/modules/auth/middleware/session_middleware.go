package middleware

import (
	"fmt"
	"gfly/internal/constants"
	"gfly/internal/domain/repository"
	"gfly/internal/http/response"
	"gfly/pkg/modules/auth"
	"github.com/gflydev/core"
	"github.com/gflydev/core/errors"
	"github.com/gflydev/core/log"
	"github.com/gflydev/core/try"
	"github.com/gflydev/core/utils"
	"github.com/gflydev/utils/str"
	"slices"
)

func processSession(c *core.Ctx) (err error) {
	try.Perform(func() {
		// Just get session to trigger updating value TTL.
		username := c.GetSession(auth.SessionUsername)

		// Check Logged-in data
		if username == nil || username.(string) == "" {
			try.Throw("No username in session")
		}

		// Put logged-in user to request data pool.
		user := repository.Pool.GetUserByEmail(username.(string))
		c.SetData(constants.User, *user)
	}).Catch(func(e try.E) {
		err = errors.New("%v", e)
	})

	return
}

func loginUrl(c *core.Ctx) string {
	// Check from `internal/http/routes/web_routes.go` to update .env file
	authLoginUrl := utils.Getenv(constants.AuthLoginUri, "/login")

	return authLoginUrl + "?redirect_url=" + c.OriginalURL()
}

// SessionAuth an HTTP middleware that process login via Session/Cookie token for API or Page requests.
//
// Use:
//
//	apiRouter.Use(middleware.SessionAuth(
//		prefixAPI+"/frontend/auth/signin",
//		"/dang-nhap",
//	))
func SessionAuth(excludes ...string) core.MiddlewareHandler {
	return func(c *core.Ctx) (err error) {
		err = processSession(c)
		path := c.Path()

		if slices.Contains(excludes, path) {
			log.Tracef("Skip SessionAuth checking for '%v'", path)
			err = nil

			return
		}

		// Response for API request
		prefixAPI := fmt.Sprintf(
			"/%s/%s",
			utils.Getenv("API_PREFIX", "api"),
			utils.Getenv("API_VERSION", "v1"),
		)
		if err != nil && str.StartsWith(path, prefixAPI) {
			return c.Error(response.Error{
				Message: err.Error(),
				Code:    core.StatusUnauthorized,
			}, core.StatusUnauthorized)
		}

		// Response for Page request
		if err != nil {
			_ = c.Redirect(loginUrl(c))
		}

		return
	}
}

// SessionAuthPage an HTTP middleware that process login via Session/Cookie token.
//
// Note:
//
//   - SessionAuthPage and SessionManipulation are used together
//     if you want to have a full manual control user's session for a specific Handler.
//
// Use:
//
//	groupUsers.GET("/profile", f.Apply(middleware.SessionAuthPage)(user.NewAccountPage()))
func SessionAuthPage(c *core.Ctx) (err error) {
	if c.GetData(constants.User) == nil {
		_ = c.Redirect(loginUrl(c))
	}

	return
}

// SessionManipulation an HTTP middleware that tries to process Session.
//
// Note:
//
//   - Place first and before the webpage routers declarations
//   - SessionAuthPage and SessionManipulation are used together
//     if you want to have a full manual control user's session for a specific Handler.
//
// Use:
//
//	webRouter.Use(middleware.SessionManipulation)
func SessionManipulation(c *core.Ctx) (err error) {
	try.Perform(func() {
		_ = processSession(c)
	}).Catch(func(e try.E) {
		log.Errorf("error %v", e)
	})

	return
}
