package notifications

import (
	"github.com/gflydev/core"
	notifyMail "github.com/gflydev/notification/mail"
	view "github.com/gflydev/view/pongo"
)

type ResetPassword struct {
	ID    int
	Email string
	Name  string
	Token string
}

func (n ResetPassword) ToEmail() notifyMail.Data {
	body := view.New().Parse("mails/forgot_password", core.Data{
		// For primary template
		"title":    "Reset password",
		"base_url": core.AppURL,
		"email":    n.Email,
		// For reset_password template
		"user_name": n.Name,
		"token":     n.Token,
	})

	return notifyMail.Data{
		To:      n.Email,
		Subject: "Reset password",
		Body:    body,
	}
}
