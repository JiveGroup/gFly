package commands

import (
	"gfly/app/domain/models"
	"github.com/gflydev/console"
	"github.com/gflydev/core/log"
	mb "github.com/gflydev/db"
	"time"
)

// ---------------------------------------------------------------
//                      Register command.
// ./artisan cmd:run db-test
// ---------------------------------------------------------------

// Auto-register command.
func init() {
	console.RegisterCommand(&dbCommand{}, "db-test")
}

// ---------------------------------------------------------------
//                      DBCommand struct.
// ---------------------------------------------------------------

// DBCommand struct for hello command.
type dbCommand struct {
	console.Command
}

// Handle Process command.
func (c *dbCommand) Handle() {
	user, err := mb.GetModelBy[models.User]("email", "admin@gfly.dev")
	if err != nil || user == nil {
		log.Panic(err)
	}
	log.Infof("User %v\n", user)

	log.Infof("DBCommand :: Run at %s", time.Now().Format("2006-01-02 15:04:05"))
}
