package routes

import (
	"fmt"
	"gfly/app/modules/auth"
	"gfly/app/modules/auth/api"
	"gfly/app/modules/auth/middleware"
	"github.com/gflydev/core"
	"github.com/gflydev/core/utils"
)

// RegisterApi func for describe a group of API routes.
func RegisterApi(apiRouter *core.Group) {
	prefixAPI := fmt.Sprintf(
		"/%s/%s",
		utils.Getenv("API_PREFIX", "api"),
		utils.Getenv("API_VERSION", "v1"),
	)

	apiRouter.Use(middleware.JWTAuth(
		// Frontend APIs
		prefixAPI+"/frontend/auth/signin",
		prefixAPI+"/frontend/auth/signout",

		// Backend APIs
		prefixAPI+"/auth/signin",
		prefixAPI+"/auth/signup",
		prefixAPI+"/auth/refresh",
		prefixAPI+"/password/forgot",
		prefixAPI+"/password/reset",
	))

	// Frontend APIs
	apiRouter.Group("/frontend", func(frontendRouter *core.Group) {
		// Auth APIs
		frontendRouter.Group("/auth", func(authGroup *core.Group) {
			authGroup.POST("/signin", api.NewSignInApi(auth.TypeWeb))
			authGroup.DELETE("/signout", api.NewSignOutApi(auth.TypeWeb))
		})
	})

	/* ============================ Auth Group ============================ */
	apiRouter.Group("/auth", func(authGroup *core.Group) {
		authGroup.POST("/signin", api.NewSignInApi(auth.TypeAPI))
		authGroup.DELETE("/signout", api.NewSignOutApi(auth.TypeAPI))
		authGroup.POST("/signup", api.NewSignUpApi())
		authGroup.PUT("/refresh", api.NewRefreshTokenApi())
	})

	/* ============================ Password Group ============================ */
	// Forgot password APIs
	apiRouter.Group("/password", func(passwordGroup *core.Group) {
		passwordGroup.POST("/forgot", api.NewForgotPWApi())
		passwordGroup.POST("/reset", api.NewResetPWApi())
	})
}
