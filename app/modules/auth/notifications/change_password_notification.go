package notifications

import (
	"github.com/gflydev/core"
	notifyMail "github.com/gflydev/notification/mail"
	view "github.com/gflydev/view/pongo"
)

type ChangePassword struct {
	ID    int
	Email string
	Name  string
}

func (n ChangePassword) ToEmail() notifyMail.Data {
	body := view.New().Parse("mails/change_password", core.Data{
		// For primary template
		"title":    "Change password",
		"base_url": core.AppURL,
		"email":    n.Email,
		// For changed_password template
		"user_name": n.Name,
	})

	return notifyMail.Data{
		To:      n.Email,
		Subject: "Change password",
		Body:    body,
	}
}
