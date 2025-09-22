package commands

import (
	"gfly/internal/console/queues"
	"github.com/gflydev/console"
	"github.com/gflydev/core/log"
	"time"
)

// ---------------------------------------------------------------
// 					Register command.
// ./artisan cmd:run hello-world
// ---------------------------------------------------------------

// Auto-register command.
func init() {
	console.RegisterCommand(&helloCommand{}, "hello-world")
}

// ---------------------------------------------------------------
// 					HelloCommand struct.
// ---------------------------------------------------------------

// HelloCommand struct for hello command.
type helloCommand struct {
	console.Command
}

// Handle Process command.
func (c *helloCommand) Handle() {
	// Dispatch a task into Queue.
	console.DispatchTask(queues.NewHelloTask("Hello"))

	log.Infof("HellCommand :: Run at %s", time.Now().Format("2006-01-02 15:04:05"))
}
