package routes

import (
	"fmt"
	"gfly/app/domain/models/types"
	"gfly/app/http/controllers/api"
	"gfly/app/http/controllers/api/user"
	"gfly/app/http/middleware"
	"gfly/app/modules/jwt"
	jwtApi "gfly/app/modules/jwt/api"
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

		/* ============================ Auth Middleware ============================ */
		apiRouter.Use(jwt.New(
			prefixAPI+"/auth/signin",
			prefixAPI+"/auth/signup",
			prefixAPI+"/auth/refresh",
			prefixAPI+"/forgot-password/request",
			prefixAPI+"/forgot-password/reset",
		))

		/* ============================ Auth Group ============================ */
		apiRouter.Group("/auth", func(authGroup *core.Group) {
			authGroup.POST("/signin", jwtApi.NewSignInApi())
			authGroup.DELETE("/signout", jwtApi.NewSignOutApi())
			authGroup.POST("/signup", jwtApi.NewSignUpApi())
			authGroup.PUT("/refresh", jwtApi.NewRefreshTokenApi())
		})

		/* ============================ User Group ============================ */
		apiRouter.Group("/users", func(userRouter *core.Group) {
			// Allow admin permission to access `/users/*` API
			userRouter.Use(middleware.CheckRolesMiddleware(
				[]types.Role{types.RoleAdmin},
				prefixAPI+"/users/profile",
			))

			userRouter.GET("", user.NewGetUsersApi())
			userRouter.POST("", user.NewCreateUserApi())
			userRouter.PUT("/{id}/status", r.Middleware(middleware.PreventUpdateYourSelf)(user.NewUpdateUserStatusApi()))
			userRouter.PUT("/{id}", r.Middleware(middleware.PreventUpdateYourSelf)(user.NewUpdateUserApi()))
			userRouter.DELETE("/{id}", r.Middleware(middleware.PreventUpdateYourSelf)(user.NewDeleteUserApi()))
			userRouter.GET("/{id}", user.NewGetUserByIdApi())
			userRouter.GET("/profile", user.NewGetUserProfileApi())
		})
	})
}
