package routes

import (
	"fmt"
	"gfly/app/modules/auth/api"
	"gfly/app/modules/auth/middleware"
	"github.com/gflydev/core"
	"github.com/gflydev/core/utils"
)

// Register func for describe a group of API routes.
func Register(apiRouter *core.Group) {
	prefixAPI := fmt.Sprintf(
		"/%s/%s",
		utils.Getenv("API_PREFIX", "api"),
		utils.Getenv("API_VERSION", "v1"),
	)

	apiRouter.Use(middleware.New(
		prefixAPI+"/auth/signin",
		prefixAPI+"/auth/signup",
		prefixAPI+"/auth/refresh",
		prefixAPI+"/forgot-password/request",
		prefixAPI+"/forgot-password/reset",
	))

	/* ============================ Auth Group ============================ */
	apiRouter.Group("/auth", func(authGroup *core.Group) {
		authGroup.POST("/signin", api.NewSignInApi())
		authGroup.DELETE("/signout", api.NewSignOutApi())
		authGroup.POST("/signup", api.NewSignUpApi())
		authGroup.PUT("/refresh", api.NewRefreshTokenApi())
	})
}
