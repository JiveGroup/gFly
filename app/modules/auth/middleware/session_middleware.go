package middleware

import (
	"fmt"
	"gfly/app/constants"
	"gfly/app/domain/repository"
	"gfly/app/modules/auth"
	"github.com/gflydev/core"
	"github.com/gflydev/core/log"
	"github.com/gflydev/core/try"
	"slices"
)

func processSession(c *core.Ctx) (err error) {
	try.Perform(func() {
		// Just get session to trigger updating value TTL.
		username := c.GetSession(auth.SessionUsername)

		// Check Logged-in data
		if username == nil || username.(string) == "" {
			try.Throw("no username in session")
		}

		// Put logged-in user to request data pool.
		user := repository.Pool.GetUserByEmail(username.(string))
		c.SetData(constants.User, *user)
	}).Catch(func(e try.E) {
		err = fmt.Errorf("error %v", e)
	})

	return
}

// SessionAuth an HTTP middleware that process login via Session/Cookie token for API or Page requests.
//
// Use:
//
//	apiRouter.Use(middleware.SessionAuth(
//		prefixAPI+"/frontend/auth/signin",
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

		if err != nil {
			// Check from `app/http/routes/web_routes.go`
			_ = c.Redirect("/login?redirect_url=" + c.OriginalURL())
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
		// Check from `app/http/routes/web_routes.go`
		_ = c.Redirect("/login?redirect_url=" + c.OriginalURL())
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
