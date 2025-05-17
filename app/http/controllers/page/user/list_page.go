package user

import (
	"gfly/app/domain/models"
	"gfly/app/http/controllers/page"
	"github.com/gflydev/core"
	mb "github.com/gflydev/db"
	qb "github.com/jivegroup/fluentsql"
)

// ====================================================================
// ======================== Controller Creation =======================
// ====================================================================

// NewListPage As a constructor to create a Home Page.
func NewListPage() *ListPage {
	return &ListPage{}
}

type ListPage struct {
	page.BasePage
}

// ====================================================================
// ========================= Request Handling =========================
// ====================================================================

func (m *ListPage) Handle(c *core.Ctx) error {
	users, total, err := mb.FindModels[models.User](1, 100, "id", qb.Desc)
	if err != nil {
		return c.Error(err)
	}

	return m.View(c, "user", core.Data{
		"users": users,
		"total": total,
	})
}
