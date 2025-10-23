package user

import (
	"gfly/internal/dto"
	"gfly/internal/http/response"
	"gfly/internal/http/transformers"
	"gfly/internal/services"
	"gfly/pkg/constants"
	"gfly/pkg/http"
	"github.com/gflydev/core"
)

// ====================================================================
// ======================== Controller Creation =======================
// ====================================================================

type ListUsersApi struct {
	http.ListApi
}

func NewListUsersApi() *ListUsersApi {
	return &ListUsersApi{}
}

// ====================================================================
// ========================= Request Handling =========================
// ====================================================================

// Handle Process main logic for API.
// @Summary Function list all users data
// @Description Function list all users data
// @Description <b>Keyword fields:</b> roles.name, roles.slug, users.email, users.fullname, users.phone, user.status
// @Description <b>Order_by fields:</b> users.email, users.fullname, users.phone, users.status, users.last_access_at
// @Tags Users
// @Accept json
// @Produce json
// @Param keyword query string false "Keyword"
// @Param order_by query string false "Order By"
// @Param page query int false "Page"
// @Param per_page query int false "Items Per Page"
// @Failure 400 {object} http.Error
// @Failure 401 {object} http.Error
// @Success 200 {object} response.ListUser
// @Security ApiKeyAuth
// @Router /users [get]
func (h *ListUsersApi) Handle(c *core.Ctx) error {
	filterDto := c.GetData(constants.Filter).(dto.Filter)
	users, total, err := services.FindUsers(filterDto)
	if err != nil {
		return err
	}

	// Pagination metadata
	metadata := http.Meta{
		Page:    filterDto.Page,
		PerPage: filterDto.PerPage,
		Total:   total,
	}

	// Transform to response data
	data := http.ToListResponse(users, transformers.ToUserResponse)

	return c.Success(response.ListUser{
		Meta: metadata,
		Data: data,
	})
}
