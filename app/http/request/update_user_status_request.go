package request

import "gfly/app/dto"

// UpdateUserStatus struct to describe update user's status
type UpdateUserStatus struct {
	dto.UpdateUserStatus
}

// ToDto convert struct to UpdateUserStatus DTO object
func (r UpdateUserStatus) ToDto() dto.UpdateUserStatus {
	return r.UpdateUserStatus
}
