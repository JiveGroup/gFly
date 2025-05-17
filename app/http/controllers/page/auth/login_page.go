package auth

import (
	"gfly/app/constants"
	"gfly/app/http/controllers/page"
	"github.com/gflydev/core"
)

// ====================================================================
// ======================== Controller Creation =======================
// ====================================================================

// NewLoginPage As a constructor to create a Home Page.
func NewLoginPage() *LoginPage {
	return &LoginPage{}
}

type LoginPage struct {
	page.BasePage
}

// ====================================================================
// ========================= Request Handling =========================
// ====================================================================

func (m *LoginPage) Handle(c *core.Ctx) error {
	if c.GetData(constants.User) != nil {
		return c.Redirect("/profile")
	}

	return m.View(c, "login", core.Data{})
}
