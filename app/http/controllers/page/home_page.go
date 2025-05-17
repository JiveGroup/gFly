package page

import "github.com/gflydev/core"

// ====================================================================
// ======================== Controller Creation =======================
// ====================================================================

// NewHomePage As a constructor to create a Home Page.
func NewHomePage() *HomePage {
	return &HomePage{}
}

type HomePage struct {
	BasePage
}

// ====================================================================
// ========================= Request Handling =========================
// ====================================================================

func (m *HomePage) Handle(c *core.Ctx) error {
	return m.View(c, "home", core.Data{
		"hero_text": "gFly - Laravel inspired web framework written in Go",
	})
}
