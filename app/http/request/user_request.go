package request

import "gfly/app/dto"

type CreateUser struct {
	dto.CreateUser
}

type UpdateUser struct {
	dto.UpdateUser
}

// ToDto Convert to UpdateUser DTO object.
func (r UpdateUser) ToDto() dto.UpdateUser {
	return r.UpdateUser
}

// ToDto Convert to CreateUser DTO object.
func (r CreateUser) ToDto() dto.CreateUser {
	return r.CreateUser
}
