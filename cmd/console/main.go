package main

import (
	_ "gfly/internal/console/commands"  // Autoload commands into pool.
	_ "gfly/internal/console/queues"    // Autoload tasks into queue.
	_ "gfly/internal/console/schedules" // Autoload jobs into schedule.
	"github.com/gflydev/cache"
	cacheRedis "github.com/gflydev/cache/redis"
	"github.com/gflydev/console"
	"github.com/gflydev/core"
	mb "github.com/gflydev/db"
	dbPSQL "github.com/gflydev/db/psql"
	notificationMail "github.com/gflydev/notification/mail"
	"os"
)

func main() {
	// Bootstrap
	core.Bootstrap()

	// Register mail notification
	notificationMail.AutoRegister()

	// Register Redis cache
	cache.Register(cacheRedis.New())

	// Register DB driver & Load Model builder
	mb.Register(dbPSQL.New())
	mb.Load()

	args := os.Args[1:] // Skip application name

	switch {
	case len(args) > 0 && args[0] == "schedule:run":
		/*---------------------------------------
						Scheduler
		----------------------------------------*/
		// Start scheduler
		console.StartScheduler()
	case len(args) > 0 && args[0] == "queue:run":
		/*---------------------------------------
						QueueJob
		----------------------------------------*/
		// Start queue worker
		console.StartQueueWorker()
	case len(args) > 0 && args[0] == "cmd:run":
		/*---------------------------------------
						Command
		----------------------------------------*/
		// Run command
		console.RunCommands(args[1:])
	}
}
