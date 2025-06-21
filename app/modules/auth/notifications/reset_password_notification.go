package notifications

import (
	"gfly/app/constants"
	"github.com/gflydev/core"
	"github.com/gflydev/core/utils"
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
	resetPasswordURI := utils.Getenv(constants.AuthResetPasswordUri, "/reset-password")

	body := view.New().Parse("mails/forgot_password", core.Data{
		// For primary template
		"title":    "Reset password",
		"base_url": core.AppURL,
		"email":    n.Email,
		// For reset_password template
		"user_name":          n.Name,
		"token":              n.Token,
		"reset_password_uri": resetPasswordURI,
	})

	return notifyMail.Data{
		To:      n.Email,
		Subject: "Reset password",
		Body:    body,
	}
}
