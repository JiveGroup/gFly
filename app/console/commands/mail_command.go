package commands

import (
	"gfly/app/notifications"
	"github.com/gflydev/console"
	"github.com/gflydev/core/log"
	"github.com/gflydev/notification"
	"time"
)

// ---------------------------------------------------------------
//                      Register command.
// ./artisan cmd:run mail-test
// ---------------------------------------------------------------

// Auto-register command.
func init() {
	console.RegisterCommand(&MailCommand{}, "mail-test")
}

// ---------------------------------------------------------------
//                      MailCommand struct.
// ---------------------------------------------------------------

// MailCommand struct for hello command.
type MailCommand struct {
	console.Command
}

// Handle Process command.
func (c *MailCommand) Handle() {
	// ============== Send mail ==============
	sendMail := notifications.SendMail{
		Email: "admin@gfly.dev",
	}

	if err := notification.Send(sendMail); err != nil {
		log.Error(err)
	}

	log.Infof("MailCommand :: Run at %s", time.Now().Format("2006-01-02 15:04:05"))
}
