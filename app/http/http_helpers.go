package http

import (
	"fmt"
	"gfly/app/constants"
	"gfly/app/dto"
	"gfly/app/http/response"
	"gfly/app/utils"
	"github.com/gflydev/core"
	"strconv"
)

// ---------------------- Path data ------------------------

// PathID get ID from path request
func PathID(c *core.Ctx, idName ...string) (int, *response.Error) {
	// Path name
	name := "id"
	if len(idName) > 0 {
		name = idName[0]
	}

	// Parse path parameter
	id, err := strconv.Atoi(c.PathVal(name))
	if err != nil || id < 1 {
		return id, &response.Error{
			Code:    core.StatusBadRequest,
			Message: fmt.Sprintf("%s must be positive integer", name),
		}
	}

	return id, nil
}

// ---------------------- Parse data ------------------------

// Parse get body data from request
func Parse[T any](c *core.Ctx, structData *T) *response.Error {
	// Parse request body
	err := c.ParseBody(structData)
	if err != nil {
		return &response.Error{
			Code:    core.StatusBadRequest,
			Message: err.Error(),
		}
	}

	return nil
}

// ---------------------- Filters ------------------------

func FilterData(c *core.Ctx) dto.Filter {
	// Receive request parameters
	page, _ := c.QueryInt("page")
	limit, _ := c.QueryInt("per_page")

	// Set default values.
	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 10
	}

	// Create DTO
	filterDto := dto.Filter{}
	filterDto.Keyword = c.QueryStr("keyword")
	filterDto.OrderBy = c.QueryStr("order_by")
	filterDto.Page = page
	filterDto.PerPage = limit

	return filterDto
}

// ---------------------- Validations ------------------------

// ValidateFilter Verify data from request.
func ValidateFilter(c *core.Ctx) error {
	filterDto := FilterData(c)

	// Validate DTO
	if errData := utils.Validate(filterDto); errData != nil {
		return c.Error(errData)
	}

	// Store data into context.
	c.SetData(constants.Filter, filterDto)

	return nil
}
