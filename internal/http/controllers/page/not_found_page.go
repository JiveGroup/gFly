package page

import (
	"github.com/gflydev/core"
)

// ====================================================================
// ============================ 404 page ============================
// ====================================================================

// NewNotFoundPage As a constructor to create new Page.
func NewNotFoundPage() *NotFoundPage {
	return &NotFoundPage{}
}

type NotFoundPage struct {
	core.Page
}

func (m *NotFoundPage) Handle(c *core.Ctx) error {
	// Set 404 status code
	c.Status(core.StatusNotFound)

	// Render the 404 template
	return c.View("404", core.Data{})
}
