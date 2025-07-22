package page

import (
	"database/sql"
	"gfly/app/constants"
	"gfly/app/domain/models"
	"github.com/gflydev/core"
	"github.com/gflydev/core/utils"
	dbNull "github.com/gflydev/db/null"
	"time"
)

// ====================================================================
// ======================== Controller Creation =======================
// ====================================================================

type BasePage struct {
	core.Page
}

// ====================================================================
// ========================= Request Handling =========================
// ====================================================================

func (m *BasePage) View(c *core.Ctx, template string, data core.Data) error {
	// -------------- Append data --------------
	if _, exists := data["title_page"]; !exists {
		data["title_page"] = "gFly | Laravel inspired web framework written in Go"
	}

	// Auto-load login user session
	user := c.GetData(constants.User)
	if user != nil {
		data["account"] = user.(models.User)
	}

	// -------------- Append functions --------------
	data["isPaths"] = func(paths ...string) bool {
		return utils.IndexOfStr(paths, c.Path()) >= 0
	}
	data["formatTime"] = func(t time.Time) string {
		return t.Format("2006-01-02 15:04:05")
	}
	data["nullString"] = func(v sql.NullString) string {
		return *dbNull.StringNil(v)
	}

	return c.View(template, data)
}
