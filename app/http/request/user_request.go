package request

import "gfly/app/dto"

// ====================================================================
// ========================== Add Requests ============================
// ====================================================================

// ---------------------- Create User ------------------------

type CreateUser struct {
	dto.CreateUser
}

// ToDto Convert to CreateUser DTO object.
func (r CreateUser) ToDto() dto.CreateUser {
	return r.CreateUser
}

// ====================================================================
// ========================= Update Requests ==========================
// ====================================================================

// ---------------------- Update User ------------------------

type UpdateUser struct {
	dto.UpdateUser
}

// ToDto Convert to UpdateUser DTO object.
func (r UpdateUser) ToDto() dto.UpdateUser {
	return r.UpdateUser
}

func (r UpdateUser) SetID(id int) {
	r.ID = id
}

// ------------------ Update User's status --------------------

// UpdateUserStatus struct to describe update user's status
type UpdateUserStatus struct {
	dto.UpdateUserStatus
}

func (r UpdateUserStatus) SetID(id int) {
	r.ID = id
}

// ToDto convert struct to UpdateUserStatus DTO object
func (r UpdateUserStatus) ToDto() dto.UpdateUserStatus {
	return r.UpdateUserStatus
}
