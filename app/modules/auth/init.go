package auth

const (
	// ========== JWT configurations ==========

	TtlOverDays    = "JWT_TTL_OVER_DAYS"
	Blacklist      = "JWT_BLACKLIST"
	CheckBlacklist = "JWT_CHECK_BLACKLIST"
	TtlMinutes     = "JWT_TTL_MINUTES"
	SecretKey      = "JWT_SECRET_KEY"
	RefreshKey     = "JWT_REFRESH_KEY"
)

// Tokens struct to describe tokens object.
type Tokens struct {
	Access  string
	Refresh string
}
