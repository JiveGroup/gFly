package routes

import (
	"gfly/app/http/controllers/page"
	"gfly/app/http/controllers/page/auth"
	"gfly/app/http/controllers/page/user"
	"gfly/app/modules/auth/middleware"
	authRoute "gfly/app/modules/auth/routes"
	"github.com/gflydev/core"
)

// WebRoutes func for describe a group of Web page routes.
func WebRoutes(r core.IFly) {
	// Session Manipulation (NOTE: Put code right position in web router)
	r.Use(middleware.SessionManipulation)

	// Web Routers
	r.GET("/", page.NewHomePage())

	r.GET("/login", auth.NewLoginPage())
	r.GET("/profile", r.Middleware(middleware.SessionAuthPage)(user.NewProfilePage()))
	r.GET("/users", r.Middleware(middleware.SessionAuthPage)(user.NewListPage()))

	/* ============================ Auth Group ============================ */
	authRoute.RegisterWeb(r)
}
