package models

import (
	mb "github.com/gflydev/db"
	"time"
)

// ====================================================================
// ============================ Data Types ============================
// ====================================================================

// N/A

// ====================================================================
// ============================== Table ===============================
// ====================================================================

// TableUserRole Table name
const TableUserRole = "user_roles"

// UserRole struct to describe a user role object.
type UserRole struct {
	// Table meta data
	MetaData mb.MetaData `db:"-" model:"table:user_roles"`

	// Table fields
	ID        int       `db:"id" model:"name:id; type:int,primary"`
	RoleID    int       `db:"role_id" model:"name:role_id; type:int"`
	UserID    int       `db:"user_id" model:"name:user_id; type:int"`
	CreatedAt time.Time `db:"created_at" model:"name:created_at"`
}
