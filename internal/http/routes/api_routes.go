package routes

import (
	"fmt"
	"gfly/internal/domain/models/types"
	"gfly/internal/http/controllers/api"
	"gfly/internal/http/controllers/api/user"
	"gfly/internal/http/middleware"
	authRoute "gfly/pkg/modules/auth/routes"

	"github.com/gflydev/core"
	"github.com/gflydev/core/utils"
)

// ApiRoutes func for describe a group of API routes.
func ApiRoutes(r core.IFly) {
	prefixAPI := fmt.Sprintf(
		"/%s/%s",
		utils.Getenv("API_PREFIX", "api"),
		utils.Getenv("API_VERSION", "v1"),
	)

	// API Routers
	r.Group(prefixAPI, func(apiRouter *core.Group) {
		// curl -v -X GET http://localhost:7789/api/v1/info | jq
		apiRouter.GET("/info", api.NewDefaultApi())

		/* ============================ Auth Group ============================ */
		authRoute.RegisterApi(apiRouter)

		/* ============================ User Group ============================ */
		apiRouter.Group("/users", func(userRouter *core.Group) {
			// Allow admin permission to access `/users/*` API
			userRouter.Use(middleware.CheckRolesMiddleware(
				[]types.Role{types.RoleAdmin},
				prefixAPI+"/users/profile",
			))

			preventUpdateYourSelfFunc := r.Apply(middleware.PreventUpdateYourSelf)

			userRouter.GET("", user.NewListUsersApi())
			userRouter.POST("", user.NewCreateUserApi())
			userRouter.PUT("/{id}/status", preventUpdateYourSelfFunc(user.NewUpdateUserStatusApi()))
			userRouter.PUT("/{id}", preventUpdateYourSelfFunc(user.NewUpdateUserApi()))
			userRouter.DELETE("/{id}", preventUpdateYourSelfFunc(user.NewDeleteUserApi()))
			userRouter.GET("/{id}", user.NewGetUserByIdApi())
			userRouter.GET("/profile", user.NewGetUserProfileApi())
		})
	})
}
