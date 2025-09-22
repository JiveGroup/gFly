package repository

import (
	"fmt"
	"gfly/internal/domain/models"
	"gfly/internal/domain/models/types"
	"github.com/gflydev/core/errors"
	"github.com/gflydev/core/log"
	"github.com/gflydev/core/try"
	"time"

	mb "github.com/gflydev/db" // Model builder
)

// ====================================================================
// ======================= Repository Interface =======================
// ====================================================================

// IRoleRepository defines the interface for managing user roles.
// It provides methods for retrieving and assigning roles to users.
//
// Methods:
//   - GetRolesByUserID(userID int) []models.Role: Retrieves roles of specific user.
//     Retrieves all roles associated with a specific user ID.
//   - GetRolesBySlug(roleSlugs ...types.Role) []models.Role:
//     Retrieves a list of roles that match the specified slug values.
//   - AddRoleForUserID(userID int, roleSlug types.Role) error: Add a role via slug for user.
//     Creates a new user-role mapping in the database.
//   - SyncRolesWithUser(userID int, roleSlugs ...types.Role) (err error):
//     Synchronizes roles for a given user, associating the user with the specified roles.
type IRoleRepository interface {
	// GetRolesByUserID retrieves all roles associated with the given user ID.
	// Returns a slice of Role models containing the role data.
	//
	// Parameters:
	//   - userID (int): The unique identifier of the user
	//
	// Returns:
	//   - []models.Role: Slice of Role models containing the role data
	GetRolesByUserID(userID int) []models.Role

	// GetRolesBySlug retrieves a list of roles that match the specified slug values.
	// Parameters:
	//   - roleSlugs (...types.Role): A variadic parameter representing the slugs to filter roles by.
	// Returns:
	//   - ([]models.Role): A slice of roles matching the provided slugs.
	GetRolesBySlug(roleSlugs ...types.Role) []models.Role

	// AddRoleForUserID assigns a role to a user by their ID.
	// Creates a new user-role mapping in the database.
	//
	// Parameters:
	//   - userID (int): The unique identifier of the user
	//   - roleSlug (types.Role): The slug name of the role to assign
	//
	// Returns:
	//   - error: Returns nil on success, error on failure
	AddRoleForUserID(userID int, roleSlug types.Role) error

	// SyncRolesWithUser synchronizes ro
	// les for a given user.
	// It associates the user with the specified roles by slug.
	// Parameters:
	//   - userID (int): The unique identifier of the user.
	//   - roleSlugs (...types.Role): A variadic parameter representing the slugs of roles to synchronize.
	// Returns:
	//   - (error): An error if the synchronization fails.
	SyncRolesWithUser(userID int, roleSlugs ...types.Role) (err error)
}

// ====================================================================
// ====================== Repository Implement ========================
// ====================================================================

// roleRepository struct for queries from a Role model.
// The struct is an implementation of interface IRoleRepository
type roleRepository struct{}

// GetRolesByUserID query for getting roles by given user ID.
func (q *roleRepository) GetRolesByUserID(userID int) []models.Role {
	// Define role variable.
	var roles []models.Role

	_, err := mb.Instance().Select(models.TableRole+".*").
		Join(mb.InnerJoin, models.TableUserRole, mb.Condition{
			Field: models.TableRole + ".id",
			Opt:   mb.Eq,
			Value: mb.ValueField(models.TableUserRole + ".role_id"),
		}).
		Where(models.TableUserRole+".user_id", mb.Eq, userID).
		OrderBy(models.TableRole+".name", mb.Asc).
		Find(&roles)

	if err != nil {
		log.Error(err)
	}

	// Return query result.
	return roles
}

// GetRolesBySlug retrieves a list of roles that match the specified slug values.
// Parameters:
//   - roleSlugs (...types.Role): A variadic parameter representing the slugs to filter roles by.
//
// Returns:
//   - ([]models.Role): A slice of roles matching the provided slugs.
func (q *roleRepository) GetRolesBySlug(roleSlugs ...types.Role) []models.Role {
	// Error variable
	var err error
	var roles []models.Role

	try.Perform(func() {
		_, err = mb.Instance().
			Where("slug", mb.In, types.RoleArrStr(roleSlugs...)).
			Limit(1000, 0).
			Find(&roles)
	}).Catch(func(e try.E) {
		log.Error(e)

		err = e.(error)
	})

	// For case empty list => return an empty list
	if roles == nil || err != nil {
		roles = []models.Role{}
	}

	return roles
}

func Convert[T fmt.Stringer](items []T) []string {
	result := make([]string, len(items))
	for i, item := range items {
		result[i] = item.String()
	}
	return result
}

// AddRoleForUserID query for adding role for given user ID.
func (q *roleRepository) AddRoleForUserID(userID int, roleSlug types.Role) error {
	// Get a role by slug
	role, err := mb.GetModel[models.Role](mb.Condition{
		Field: models.TableRole + ".slug",
		Opt:   mb.Eq,
		Value: roleSlug,
	})
	if err != nil || role == nil {
		log.Error(err)

		return errors.New("Role not found")
	}

	// Create a new user
	userRole := models.UserRole{
		ID:        userID,
		RoleID:    role.ID,
		UserID:    userID,
		CreatedAt: time.Now(),
	}

	return mb.CreateModel(&userRole)
}

// SyncRolesWithUser synchronizes roles for a given user by associating the user with the specified roles.
// Parameters:
//   - userID (int): The unique identifier of the user.
//   - roleSlugs (...types.Role): A variadic parameter representing the slugs of roles to synchronize.
//
// Returns:
//   - (error): An error if the synchronization fails.
func (q *roleRepository) SyncRolesWithUser(userID int, roleSlugs ...types.Role) (err error) {
	// DB Model instance
	db := mb.Instance()

	try.Perform(func() {
		db.Begin()

		// Remove old Roles associated with the user.
		if err := db.Where("user_id", mb.Eq, userID).Delete(&models.UserRole{}); err != nil {
			try.Throw(err)
		}

		roles := q.GetRolesBySlug(roleSlugs...)

		// Create a relationship between roles and user.
		for _, role := range roles {
			userRole := models.UserRole{
				RoleID:    role.ID,
				UserID:    userID,
				CreatedAt: time.Now(),
			}

			if err := db.Create(&userRole); err != nil {
				try.Throw(err)
			}
		}

		// Commit the transaction to the database.
		err = db.Commit()
	}).Catch(func(e try.E) {
		err = e.(error)
		_ = db.Rollback() // Rollback the transaction in case of failure.
	})

	return err
}
