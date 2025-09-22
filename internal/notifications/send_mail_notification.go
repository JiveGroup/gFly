package notifications

import (
	"github.com/gflydev/core"
	notifyMail "github.com/gflydev/notification/mail"
	view "github.com/gflydev/view/pongo"
)

type SendMail struct {
	Email string
}

func (n SendMail) ToEmail() notifyMail.Data {
	body := view.New().Parse("mails/send_mail", core.Data{
		// For primary template
		"title":    "Reset password",
		"base_url": core.AppURL,
		"email":    n.Email,
		// For send_mail template
		"user_name": "John",
		"text":      "Hello World",
	})

	return notifyMail.Data{
		To:      n.Email,
		Subject: "Check Mail",
		Body:    body,
	}
}
