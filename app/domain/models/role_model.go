package models

import (
	"database/sql"
	"gfly/app/domain/models/types"
	mb "github.com/gflydev/db"
	"time"
)

// ====================================================================
// ============================ Data Types ============================
// ====================================================================

// TBD

// ====================================================================
// ============================== Table ===============================
// ====================================================================

// TableRole Table name
const TableRole = "roles"

// Role struct to describe a role object.
type Role struct {
	// Table meta data
	MetaData mb.MetaData `db:"-" model:"table:roles"`

	// Table fields
	ID        int          `db:"id" model:"name:id; type:serial,primary"`
	Name      string       `db:"name" model:"name:name"`
	Slug      types.Role   `db:"slug" model:"name:slug"`
	CreatedAt time.Time    `db:"created_at" model:"name:created_at"`
	UpdatedAt sql.NullTime `db:"updated_at" model:"name:updated_at"`
}
