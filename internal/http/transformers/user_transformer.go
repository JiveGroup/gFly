package transformers

import (
	"gfly/internal/domain/models"
	"gfly/internal/domain/repository"
	"gfly/internal/http/response"
	"github.com/gflydev/core"
	dbNull "github.com/gflydev/db/null"
	"github.com/gflydev/storage"
	"strings"
)

// PublicAvatar converts an avatar path to a public URL
//
// Parameters:
//   - avatar: The avatar file path or URL string
//
// Returns:
//   - *string: Pointer to public avatar URL string, nil if avatar is empty
func PublicAvatar(avatar string) *string {
	if avatar == "" {
		return nil
	}
	fs := storage.Instance()

	// Absolute URL
	if strings.HasPrefix(avatar, core.SchemaHTTP) {
		return &avatar
	}
	avatar = fs.Url(avatar)

	return &avatar
}

// ToRoleResponse converts a Role model to a Role response object
//
// Parameters:
//   - model: models.Role - The role model to convert
//
// Returns:
//   - response.Role: The converted role response object
func ToRoleResponse(model models.Role) response.Role {
	return response.Role{
		ID:   model.ID,
		Name: model.Name,
		Slug: model.Slug,
	}
}

// roles retrieve and convert a role list for a user
//
// Parameters:
//   - userID: The user ID to get roles for
//
// Returns:
//   - []response.Role: Array of role response objects
func roles(userID int) []response.Role {
	var roles []response.Role
	roleList := repository.Pool.GetRolesByUserID(userID)
	for _, role := range roleList {
		roles = append(roles, ToRoleResponse(role))
	}

	return roles
}

// ToUserResponse converts a User model to a User response object
// with all fields populated for general purpose response
//
// Parameters:
//   - user: models.User - The user model to convert
//
// Returns:
//   - response.User: The converted user response object
func ToUserResponse(user models.User) response.User {
	return response.User{
		ID:           user.ID,
		Email:        user.Email,
		Fullname:     user.Fullname,
		Phone:        user.Phone,
		Token:        dbNull.StringNil(user.Token),
		Status:       user.Status,
		Avatar:       PublicAvatar(user.Avatar.String),
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		VerifiedAt:   dbNull.TimeNil(user.VerifiedAt),
		BlockedAt:    dbNull.TimeNil(user.BlockedAt),
		DeletedAt:    dbNull.TimeNil(user.DeletedAt),
		LastAccessAt: dbNull.TimeNil(user.LastAccessAt),
		Roles:        roles(user.ID),
	}
}
