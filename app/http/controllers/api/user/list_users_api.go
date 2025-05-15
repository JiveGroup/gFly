package user

import (
	"gfly/app/constants"
	"gfly/app/dto"
	"gfly/app/http"
	"gfly/app/http/response"
	"gfly/app/http/transformers"
	"gfly/app/services"
	"github.com/gflydev/core"
)

type GetUsersApi struct {
	core.Api
}

func NewGetUsersApi() *GetUsersApi {
	return &GetUsersApi{}
}

func (h *GetUsersApi) Validate(c *core.Ctx) error {
	return http.ValidateFilter(c)
}

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
// @Failure 400 {object} response.Error
// @Failure 401 {object} response.Unauthorized
// @Success 200 {object} response.ListUser
// @Security ApiKeyAuth
// @Router /users [get]
func (h *GetUsersApi) Handle(c *core.Ctx) error {
	filterDto := c.GetData(constants.Filter).(dto.Filter)
	users, total, err := services.FindUsers(filterDto)
	if err != nil {
		return err
	}

	// Pagination metadata
	metadata := dto.Meta{
		Page:    filterDto.Page,
		PerPage: filterDto.PerPage,
		Total:   total,
	}

	// Transform to response data
	data := transformers.ToListResponse(users, transformers.ToUserResponse)

	return c.Success(response.ListUser{
		Meta: metadata,
		Data: data,
	})
}
