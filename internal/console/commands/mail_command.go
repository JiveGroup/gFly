package commands

import (
	"gfly/internal/notifications"
	"github.com/gflydev/console"
	"github.com/gflydev/core/log"
	"github.com/gflydev/notification"
	"time"
)

// ---------------------------------------------------------------
//                        Register command.
// ./artisan cmd:run mail-test
// ---------------------------------------------------------------

// Auto-register command.
func init() {
	console.RegisterCommand(&mailCommand{}, "mail-test")
}

// ---------------------------------------------------------------
//                        MailCommand struct.
// ---------------------------------------------------------------

// MailCommand struct for hello command.
type mailCommand struct {
	console.Command
}

// Handle Process command.
func (c *mailCommand) Handle() {
	// ============== Send mail ==============
	sendMail := notifications.SendMail{
		Email: "admin@gfly.dev",
	}

	if err := notification.Send(sendMail); err != nil {
		log.Error(err)
	}

	log.Infof("MailCommand :: Run at %s", time.Now().Format("2006-01-02 15:04:05"))
}
