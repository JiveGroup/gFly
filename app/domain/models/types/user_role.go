package types

// ====================================================================
// ============================ Data Types ============================
// ====================================================================

// ----------------------------- Role slug ----------------------------

type Role string

// User property types
const (
	RoleNA        Role = "na"
	RoleAdmin     Role = "admin"
	RoleModerator Role = "moderator"
	RoleMember    Role = "member"
	RoleGuest     Role = "guest"
)

var RoleList = []Role{
	RoleNA,
	RoleAdmin,
	RoleModerator,
	RoleMember,
	RoleGuest,
}

// ----------------------------- Role ID ------------------------------

type RoleID uint

// Role ID and id match with the database.
const (
	RoleNAID RoleID = iota
	RoleAdminID
	RoleModeratorID
	RoleMemberID
	RoleGuestID
)

var RoleMap = map[Role]RoleID{
	RoleNA:        RoleNAID,
	RoleAdmin:     RoleAdminID,
	RoleModerator: RoleModeratorID,
	RoleMember:    RoleMemberID,
	RoleGuest:     RoleGuestID,
}

// ====================================================================
// ============================= Methods ==============================
// ====================================================================

func (e RoleID) Name() Role {
	return RoleList[e]
}

func (e Role) ID() RoleID {
	return RoleMap[e]
}

// roleCollection represents a collection of Role values.
type roleCollection []Role

// String converts a collection of Roles to their string representations.
// Returns a slice containing the string value for each Role.
func (e roleCollection) String() []string {
	result := make([]string, len(e))
	for i, v := range e {
		result[i] = string(v)
	}

	return result
}

// RoleArrStr converts variadic Role parameters into a slice of strings.
// Parameters:
//   - roles (...Role): Variable number of Role parameters
//
// Returns:
//   - []string: Slice containing string representations of the roles
func RoleArrStr(roles ...Role) []string {
	return append(roleCollection{}, roles...).String()
}
