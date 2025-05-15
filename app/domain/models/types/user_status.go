package types

// ====================================================================
// ============================ Data Types ============================
// ====================================================================

type UserStatus string

// User property types
const (
	UserStatusActive  UserStatus = "active"
	UserStatusPending UserStatus = "pending"
	UserStatusBlocked UserStatus = "blocked"
)

var UserStatusList = []UserStatus{
	UserStatusActive,
	UserStatusPending,
	UserStatusBlocked,
}

// ====================================================================
// ============================= Methods ==============================
// ====================================================================

type userStatusCollection []UserStatus

// String converts a collection of UserStatus to an array of strings.
//
// Returns:
//   - []string: Array of UserStatus values as strings (e.g. ["active", "pending", "blocked"])
func (e userStatusCollection) String() []string {
	result := make([]string, len(e))
	for i, v := range e {
		result[i] = string(v)
	}

	return result
}

// UserStatusArrStr converts variable number of UserStatus to array of strings.
//
// Parameters:
//   - userStatus: Variable number of UserStatus values
//
// Returns:
//   - []string: Array of UserStatus values as strings (e.g. ["active", "pending", "blocked"])
func UserStatusArrStr(userStatus ...UserStatus) []string {
	return append(userStatusCollection{}, userStatus...).String()
}
