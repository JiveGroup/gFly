package page

import (
	"github.com/gflydev/core"
	"github.com/gflydev/core/log"
	"time"
)

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
	c.SetSession("time", time.Now())
	now := c.GetSession("time").(time.Time)

	log.Infof("Access at: %s", now.Format("2006-01-02 15:04:05"))

	return m.View(c, "home", core.Data{
		"hero_text": "gFly - Laravel inspired web framework written in Go. Access time: " + now.Format("2006-01-02 15:04:05"),
	})
}
