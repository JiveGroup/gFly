package auth

import (
	"gfly/internal/http/controllers/page"
	"github.com/gflydev/core"
	"github.com/gflydev/http"
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
	if c.GetData(http.UserKey) != nil {
		return c.Redirect("/profile")
	}

	return m.View(c, "login", core.Data{})
}
