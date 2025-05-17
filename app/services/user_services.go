package services

import (
	"fmt"
	"gfly/app/domain/models"
	"gfly/app/domain/models/types"
	"gfly/app/domain/repository"
	"gfly/app/dto"
	"github.com/gflydev/core"
	"github.com/gflydev/core/errors"
	"github.com/gflydev/core/log"
	coreUtils "github.com/gflydev/core/utils"
	mb "github.com/gflydev/db"
	dbNull "github.com/gflydev/db/null"
	qb "github.com/jivegroup/fluentsql"
	"slices"
	"strings"
	"time"
)

const (
	UploadAvatarDir = "avatars"
)

// ====================================================================
// ========================= Main functions ===========================
// ====================================================================

// FindUsers retrieves a list of users from the database based on the provided filter criteria.
// It supports searching by keyword, ordering by specified fields, and pagination.
//
// Parameters:
//   - filterDto (dto.Filter): The filter containing search criteria, order by field, page, and per-page details.
//
// Returns:
//
//	([]models.User, int, error): A list of user models, the total number of users, and any error encountered.
func FindUsers(filterDto dto.Filter) ([]models.User, int, error) {
	// DB Model instance
	dbInstance := mb.Instance()
	// Error variable
	var err error

	// Define User variable.
	var users []models.User
	var total int
	var offset = 0

	if filterDto.Page > 0 {
		offset = (filterDto.Page - 1) * filterDto.PerPage
	}

	builder := dbInstance.Select("DISTINCT users.id", "users.*").
		Join(qb.LeftJoin, models.TableUserRole, qb.Condition{
			Field: models.TableUserRole + ".user_id",
			Opt:   qb.Eq,
			Value: qb.ValueField(models.TableUser + ".id"),
		}).
		Join(qb.LeftJoin, models.TableRole, qb.Condition{
			Field: models.TableRole + ".id",
			Opt:   qb.Eq,
			Value: qb.ValueField(models.TableUserRole + ".role_id"),
		}).
		Where(models.TableUser+".deleted_at", qb.Null, nil).
		When(filterDto.Keyword != "", func(query qb.WhereBuilder) *qb.WhereBuilder {
			query.WhereGroup(func(queryGroup qb.WhereBuilder) *qb.WhereBuilder {
				queryGroup.Where(models.TableRole+".name", qb.Like, "%"+filterDto.Keyword+"%").
					WhereOr(models.TableRole+".slug", qb.Like, "%"+filterDto.Keyword+"%").
					WhereOr(models.TableUser+".email", qb.Like, "%"+filterDto.Keyword+"%").
					WhereOr(models.TableUser+".fullname", qb.Like, "%"+filterDto.Keyword+"%").
					WhereOr(models.TableUser+".phone", qb.Like, "%"+filterDto.Keyword+"%")

				if slices.Contains(types.UserStatusList, types.UserStatus(filterDto.Keyword)) {
					queryGroup.WhereOr(models.TableUser+".status", qb.Eq, filterDto.Keyword)
				}

				return &queryGroup
			})

			return &query
		}).
		Limit(filterDto.PerPage, offset)

	if filterDto.OrderBy != "" {
		// Default order by
		direction := qb.Asc
		orderKey := filterDto.OrderBy

		if strings.HasPrefix(filterDto.OrderBy, "-") {
			orderKey = filterDto.OrderBy[1:]
			direction = qb.Desc
		}

		var orderByFields = core.Data{
			"id":          fmt.Sprintf("%s.id", models.TableUser),
			"email":       fmt.Sprintf("%s.email", models.TableUser),
			"fullname":    fmt.Sprintf("%s.fullname", models.TableUser),
			"phone":       fmt.Sprintf("%s.phone", models.TableUser),
			"status":      fmt.Sprintf("%s.status", models.TableUser),
			"last_access": fmt.Sprintf("%s.last_access_at", models.TableUser),
		}

		if field, ok := orderByFields[orderKey]; ok {
			builder.OrderBy(field.(string), direction)
		}
	}

	// Query data
	total, err = builder.Find(&users)

	// Return query result.
	return users, total, err
}

// CreateUser creates a new user in the system.
//
// This function performs the following steps:
// 1. Verifies that no other user exists with the same email.
// 2. Processes and uploads the user's avatar if provided.
// 3. Creates a new user entity in the database.
// 4. Assigns roles to the user (default role is "user" if no roles are provided).
//
// Parameters:
//   - createUserDto (dto.CreateUser): The payload containing the user's details.
//
// Returns:
//   - (*models.User, error): The created user object or an error if any step fails.
func CreateUser(createUserDto dto.CreateUser) (*models.User, error) {
	// Check if the user with the given email already exist
	userEmail := repository.Pool.GetUserByEmail(createUserDto.Email)
	if userEmail != nil {
		return nil, errors.New("the user with given email already exist")
	}

	// Create a new user
	user := &models.User{
		Status:       types.UserStatusActive,
		Email:        createUserDto.Email,
		Password:     coreUtils.GeneratePassword(createUserDto.Password),
		Fullname:     createUserDto.Fullname,
		Phone:        createUserDto.Phone,
		Token:        dbNull.String(""),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		LastAccessAt: dbNull.NowTime(),
		Avatar:       dbNull.String(createUserDto.Avatar),
	}

	if createUserDto.Status != "" {
		user.Status = types.UserStatus(createUserDto.Status)
	}

	err := mb.CreateModel(user)
	if err != nil {
		log.Errorf("Error while creating new user %v", err)
		return nil, errors.New("error occurs while creating new user")
	}

	if err = repository.Pool.SyncRolesWithUser(user.ID, createUserDto.Roles...); err != nil {
		log.Errorf("Error while syncing roles to user %v", err)
		return nil, errors.New("error occurs while syncing user roles")
	}

	return user, nil
}

// UpdateUser updates an existing user in the system.
//
// This function fetches the user by their ID, updates the fields based on the given DTO,
// and synchronizes the user's roles if provided.
//
// Parameters:
//   - updateUserDto (dto.UpdateUser): The payload containing the user's updated details.
//
// Returns:
//   - (*models.User, error): The updated user object or an error if any step fails.
//
// Possible Errors:
//   - "User not found": Returned when no user is found for the provided ID.
//   - "Error occurs while updating user": Returned when an error occurs during the update process.
//   - "Error occurs while syncing user roles": Returned when an error occurs during the role synchronization process.
func UpdateUser(updateUserDto dto.UpdateUser) (*models.User, error) {
	user, err := mb.GetModelByID[models.User](updateUserDto.ID)
	if err != nil {
		return nil, errors.New("User not found")
	}

	// Update the fields that are provided in the updateUserDto
	updatedUser := updateUserFromDto(user, updateUserDto)

	if err = mb.UpdateModel(updatedUser); err != nil {
		log.Errorf("Error while updating user %v", err)
		return nil, errors.New("Error occurs while updating user")
	}

	// Sync user roles
	if len(updateUserDto.Roles) > 0 {
		if err = repository.Pool.SyncRolesWithUser(user.ID, updateUserDto.Roles...); err != nil {
			log.Errorf("Error while syncing user roles %v", err)
			return nil, errors.New("error occurs while syncing user roles")
		}
	}

	return user, nil
}

// UpdateUserStatus updates the status of an existing user in the system.
//
// This function performs the following steps:
// 1. Finds the user by their ID.
// 2. Validates the new status value against allowed user states.
// 3. Updates the user's status in the database.
//
// Parameters:
//   - updateUserStatusDto (dto.UpdateUserStatus): The payload containing the user ID and the new status.
//
// Returns:
//   - (*models.User, error): The updated user object or an error if any step fails.
//
// Possible Errors:
//   - "User not found": Returned when no user is found for the provided ID.
//   - "Error occurs while updating user's status": Returned when the provided status is invalid or the update process fails.
func UpdateUserStatus(updateUserStatusDto dto.UpdateUserStatus) (*models.User, error) {
	user, err := mb.GetModelByID[models.User](updateUserStatusDto.ID)
	if err != nil {
		return nil, errors.New("User not found")
	}

	// Check user's status
	if !slices.Contains(types.UserStatusList, types.UserStatus(updateUserStatusDto.Status)) {
		return nil, errors.New("Error occurs while updating user's status %v", updateUserStatusDto.Status)
	}

	// Set new status and update user
	user.Status = types.UserStatus(updateUserStatusDto.Status)

	if err = mb.UpdateModel(user); err != nil {
		log.Errorf("Error while updating user %v", err)

		return nil, errors.New("Error occurs while updating user")
	}

	return user, nil
}

// DeleteUserByID deletes a user and all associated roles from the system.
//
// This function performs the following steps:
// 1. Fetches the user by their ID.
// 2. Deletes all roles synchronized with the user.
// 3. Deletes the user from the database.
//
// Parameters:
//   - userID (int): The unique identifier of the user to be deleted.
//
// Returns:
//   - error: An error object if any step fails. Possible errors include:
//   - "User not found": Returned when no user is found for the provided ID.
//   - "Error occurs while deleting user roles": Returned when an error occurs while deleting roles synchronized with the user.
//   - "Error occurs while deleting user": Returned when an error occurs during the deletion of the user record.
func DeleteUserByID(userID int) error {
	user, err := mb.GetModelByID[models.User](userID)
	if err != nil {
		return errors.New("User not found")
	}

	// Delete roles that sync with user
	if err := repository.Pool.SyncRolesWithUser(userID, ""); err != nil {
		log.Errorf("Error while deleting user roles: %v", err)
		return errors.New("error occurs while deleting user roles")
	}

	// Delete user
	if err := mb.DeleteModel(user); err != nil {
		log.Errorf("Error while deleting user: %v", err)
		return errors.New("error occurs while deleting user")
	}

	return nil
}

// UserHasRole checks if a user has any of the specified roles.
//
// Parameters:
//   - userID (int): The ID of the user whose roles are being checked.
//   - roles ([]models.RoleType): A list of roles to check against.
//
// Returns:
//   - bool: True if the user has at least one of the specified roles, otherwise false.
func UserHasRole(userID int, roles []types.Role) bool {
	roleList := repository.Pool.GetRolesByUserID(userID)
	if len(roleList) == 0 {
		return false
	}

	for _, role := range roleList {
		if slices.Contains(roles, role.Slug) {
			return true
		}
	}

	return false
}

// ====================================================================
// ======================== Helper Functions ==========================
// ====================================================================

// updateUserFromDto updates an existing User model with data from UpdateUser DTO.
// Only updates fields that are provided in the DTO.
//
// Parameters:
//   - user (*models.User): The existing user model to update
//   - updateUserDto (dto.UpdateUser): The DTO containing update data
//
// Returns:
//   - *models.User: The updated User model
func updateUserFromDto(user *models.User, updateUserDto dto.UpdateUser) *models.User {
	if updateUserDto.Password != "" {
		user.Password = coreUtils.GeneratePassword(updateUserDto.Password)
	}
	if updateUserDto.Fullname != "" && updateUserDto.Fullname != user.Fullname {
		user.Fullname = updateUserDto.Fullname
	}
	if updateUserDto.Phone != "" && updateUserDto.Phone != user.Phone {
		user.Phone = updateUserDto.Phone
	}

	// Perform upload file if there is avatar
	if updateUserDto.Avatar != "" && updateUserDto.Avatar != user.Avatar.String {
		user.Avatar = dbNull.String(updateUserDto.Avatar)
	}

	user.UpdatedAt = time.Now()

	return user
}
