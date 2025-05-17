package routes

import (
	"gfly/app/modules/auth/middleware"
	"github.com/gflydev/core"
)

// RegisterWeb func for describe a group of Web page routes.
func RegisterWeb(r core.IFly) {
	// Session Manipulation (NOTE: Put code right position in web router)
	r.Use(middleware.SessionManipulation)
}
