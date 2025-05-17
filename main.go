package main

import (
	"gfly/app/http/routes"
	"gfly/docs"
	"github.com/gflydev/cache"
	cacheRedis "github.com/gflydev/cache/redis"
	"github.com/gflydev/core"
	"github.com/gflydev/core/log"
	mb "github.com/gflydev/db"
	dbPSQL "github.com/gflydev/db/psql"
	notificationMail "github.com/gflydev/notification/mail"
	"github.com/gflydev/session"
	sessionRedis "github.com/gflydev/session/redis"
	"github.com/gflydev/storage"
	storageLocal "github.com/gflydev/storage/local"
	"github.com/gflydev/view/pongo"
	"github.com/joho/godotenv"
)

// Main function
// @title API
// @version 1.0
// @description gFly API. <b>Note: Authorize value `Bearer <ACCESS_TOKEN>`</b>.
// @termsOfService https://swagger.io/terms/
// @contact.name API Support
// @contact.email admin@jivecode.com
// @license.name Apache 2.0
// @license.url https://mit-license.org/
// @BasePath /api/v1
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	// Swagger info
	docs.SwaggerInfo.Title = "Swagger Example API"
	docs.SwaggerInfo.Description = "This is a sample server gFly."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "gFly.dev"
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file %v", err)
	}

	// Register view
	core.RegisterView(pongo.New())

	// Register mail notification
	notificationMail.AutoRegister()

	// Register Local storage
	storage.Register(storageLocal.Type, storageLocal.New())

	// Setup session
	session.Register(sessionRedis.New())
	core.RegisterSession(session.New())

	// Register Redis cache
	cache.Register(cacheRedis.New())

	// Register DB driver & Load Model builder
	mb.Register(dbPSQL.New())
	mb.Load()

	// Initial application
	app := core.New()

	// Register router
	app.RegisterRouter(routes.Router)

	// Run application
	app.Run()
}
