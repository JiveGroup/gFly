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
		log.Debugf("processSession error '%v'", e)

		err = fmt.Errorf("error %v", e)
	})

	return
}

// SessionAuthPage an HTTP middleware that process login via Session/Cookie token.
//
// Use:
//
//	groupUsers.GET("/profile", f.Middleware(middleware.SessionAuthPage)(user.NewAccountPage()))
func SessionAuthPage(c *core.Ctx) (err error) {
	if err = processSession(c); err != nil {
		// Check from `app/http/routes/web_routes.go`
		_ = c.Redirect("/login?redirect_url=" + c.OriginalURL())
	}

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
		path := c.Path()

		if slices.Contains(excludes, path) {
			log.Tracef("Skip SessionAuth checking for '%v'", path)

			return nil
		}

		return processSession(c)
	}
}

// SessionManipulation an HTTP middleware that process updating Session/Cookie's TTL for each Page request.
//
// Use:
//
//	f.Use(middleware.SessionManipulation)
func SessionManipulation(c *core.Ctx) (err error) {
	try.Perform(func() {
		_ = processSession(c)
	}).Catch(func(e try.E) {
		err = fmt.Errorf("error %v", e)
	})

	return err
}
