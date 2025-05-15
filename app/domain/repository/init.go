package repository

// ====================================================================
// ======================== Repository factory ========================
// ====================================================================

// Repositories struct for collect all app repositories.
type Repositories struct {
	IRoleRepository
	IUserRepository
}

// Pool a repository pool to store all
var Pool = &Repositories{
	&roleRepository{},
	&userRepository{},
}
