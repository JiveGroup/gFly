package user

import (
	"gfly/internal/http/controllers/page"
	"github.com/gflydev/core"
)

// ====================================================================
// ======================== Controller Creation =======================
// ====================================================================

// NewProfilePage As a constructor to create a Home Page.
func NewProfilePage() *ProfilePage {
	return &ProfilePage{}
}

type ProfilePage struct {
	page.BasePage
}

// ====================================================================
// ========================= Request Handling =========================
// ====================================================================

func (m *ProfilePage) Handle(c *core.Ctx) error {
	return m.View(c, "profile", core.Data{})
}
