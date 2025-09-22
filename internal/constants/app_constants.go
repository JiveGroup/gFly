package constants

const (
	// ====================================================================
	// ======================== Common Constants ==========================
	// ====================================================================

	// AppName Application name
	AppName string = "gFly"

	// ====================================================================
	// ===================== HTTP Context Constants =======================
	// ====================================================================

	// User key in Context's Data for User information
	User string = "__user__"
	// PathID key in Context's Data for ID extracted from path parameter
	PathID string = "__path_id__"
	// Request key in Context's Data for raw request data
	Request string = "__request__"
	// Data key in Context's Data for processed/transformed request data
	Data string = "__data__"
	// Filter key in Context's Data for filtering parameters
	Filter string = "__filter__"

	// ====================================================================
	// ========================= Auth Constants ===========================
	// ====================================================================

	AuthResetPasswordUri string = "AUTH_RESET_PASSWORD_URI"
	AuthLoginUri         string = "AUTH_LOGIN_URI"
)
