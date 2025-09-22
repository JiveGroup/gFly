package routes

import (
	"github.com/gflydev/core"
)

func Router(r core.IFly) {
	ApiRoutes(r) // Register API routes.
	WebRoutes(r) // Register Web routes.
}
