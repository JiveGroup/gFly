package models

import (
	"database/sql"
	"gfly/internal/domain/models/types"
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

// TableUser Table name
const TableUser = "users"

// User struct to describe a user object.
type User struct {
	// Table meta data
	MetaData mb.MetaData `db:"-" model:"table:users"`

	// Table fields
	ID           int              `db:"id" model:"name:id; type:serial,primary"`
	Email        string           `db:"email" model:"name:email"`
	Password     string           `db:"password" model:"name:password"` // #nosec G117 -- Legitimate password field in user model (stored hashed)
	Fullname     string           `db:"fullname" model:"name:fullname"`
	Phone        string           `db:"phone" model:"name:phone"`
	Token        sql.NullString   `db:"token" model:"name:token"`
	Status       types.UserStatus `db:"status" model:"name:status"`
	CreatedAt    time.Time        `db:"created_at" model:"name:created_at"`
	Avatar       sql.NullString   `db:"avatar" model:"name:avatar"`
	UpdatedAt    time.Time        `db:"updated_at" model:"name:updated_at"`
	VerifiedAt   sql.NullTime     `db:"verified_at" model:"name:verified_at"`
	BlockedAt    sql.NullTime     `db:"blocked_at" model:"name:blocked_at"`
	DeletedAt    sql.NullTime     `db:"deleted_at" model:"name:deleted_at"`
	LastAccessAt sql.NullTime     `db:"last_access_at" model:"name:last_access_at"`
}
