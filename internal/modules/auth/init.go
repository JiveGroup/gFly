package auth

type Type string

const (
	SessionUsername = "username"

	// ========== Auth Type ==========

	TypeAPI Type = "api"
	TypeWeb Type = "web"

	// ========== JWT configurations ==========

	TtlOverDays    = "JWT_TTL_OVER_DAYS"
	Blacklist      = "JWT_BLACKLIST"
	CheckBlacklist = "JWT_CHECK_BLACKLIST"
	TtlMinutes     = "JWT_TTL_MINUTES"
	SecretKey      = "JWT_SECRET_KEY"
	RefreshKey     = "JWT_REFRESH_KEY"
)

// Token struct to describe tokens object.
type Token struct {
	Access  string
	Refresh string
}
