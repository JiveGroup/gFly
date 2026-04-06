package models

import "github.com/gflydev/core"

// ====================================================================
// ========================= User Searchable ==========================
// ====================================================================

// SearchIndex returns the table / Elasticsearch index name for User.
func (u User) SearchIndex() string { return TableUser }

// SearchKey returns the primary key used to address the document.
func (u User) SearchKey() any { return u.ID }

// SearchableFields lists the columns the DatabaseDriver will ILIKE against
// when a keyword is supplied in the search request.
// Field names are unqualified so they work in both SQL (single-table query)
// and as Elasticsearch document field names.
func (u User) SearchableFields() []string {
	return []string{
		"fullname",
		"email",
		"phone",
	}
}

// ToSearchDocument returns the flat map written to the Elasticsearch index
// when this user is indexed.
func (u User) ToSearchDocument() core.Data {
	return core.Data{
		"id":       u.ID,
		"fullname": u.Fullname,
		"email":    u.Email,
		"phone":    u.Phone,
		"status":   string(u.Status),
	}
}
