package http

import (
	"gfly/internal/constants"
	"gfly/pkg/security"
	"github.com/gflydev/core"
)

// ====================================================================
// ======================== Other Request Helpers =====================
// ====================================================================

// ProcessPathID is a generic function that extracts a path ID parameter and stores it in the context.
// It handles the common pattern of validating a path ID parameter for API endpoints and putting it in Ctx's Data.
//
// Parameters:
//   - c: The context object containing the HTTP request/response data
//
// Returns:
//   - error: Returns nil if successful, otherwise returns an error response
//
// Example Usage:
//
//	func (h DeleteUserApi) Validate(c *core.Ctx) error {
//		return http.ProcessPathID(c)
//	}
func ProcessPathID(c *core.Ctx) error {
	// Receive path parameter ID
	itemID, errData := PathID(c)
	if errData != nil {
		return c.Error(errData)
	}

	// Store data into context
	c.SetData(constants.PathID, itemID)

	return nil
}

// ProcessFilter validates and processes filter requests
// It handles parsing the query parameters, converting to DTO, and validation and put to Ctx's Data
//
// Parameters:
//   - c: The context object containing the HTTP request/response data
//
// Returns:
//   - error: Returns nil if successful, otherwise returns an error response
//
// Example Usage:
//
//	func (h ListUserApi) Validate(c *core.Ctx) error {
//		return http.ProcessFilter(c)
//	}
func ProcessFilter(c *core.Ctx) error {
	filterDto := FilterData(c)

	// Validate DTO
	if errData := Validate(filterDto); errData != nil {
		return c.Error(errData)
	}

	// Store data into context.
	c.SetData(constants.Filter, filterDto)

	return nil
}

// ====================================================================
// ======================= Update Request Helpers =====================
// ====================================================================

// UpdateData is an interface for types that can be updated.
// IMPORTANT: SetID must be implemented with a pointer receiver (e.g., func (r *RequestType) SetID(id int))
// to ensure the ID is properly set on the original struct, not a copy
type UpdateData interface {
	// SetID sets the ID field of the request structure
	// Must be implemented with a pointer receiver to modify the struct in place
	// Parameters:
	//   - id: Integer ID value to set
	SetID(int)
}

// ProcessUpdateData validates and processes update requests.
// It handles parsing the request body, setting the ID, converting to DTO, and validation and put to Ctx's Data.
//
// Type Parameters:
//   - T: The type that implements the UpdateData interface.
//
// Parameters:
//   - c: The context object containing the HTTP request/response data.
//
// Returns:
//   - error: Returns nil if successful, otherwise returns an error response.
func ProcessUpdateData[T UpdateData](c *core.Ctx) error {
	// Receive path parameter ID
	itemID, errData := PathID(c)
	if errData != nil {
		return c.Error(errData)
	}

	// Receive request data
	var requestData T
	if errData := Parse(c, &requestData); errData != nil {
		return c.Error(errData)
	}

	// Sanitize request data
	security.SanitizeStruct(&requestData)

	// Set ID on the request body
	requestData.SetID(itemID)

	// Validate DTO
	if errData := Validate(requestData); errData != nil {
		return c.Error(errData)
	}

	// Store data into context
	c.SetData(constants.Request, requestData)

	return nil
}

// ====================================================================
// ======================== Add Request Helpers =======================
// ====================================================================

// AddData is an interface for types that can be added.
type AddData interface {
}

// ProcessData validates and processes create/add requests.
// It handles parsing the request body, converting to DTO, and validation and put to Ctx's Data.
//
// Type Parameters:
//   - T: The type that implements the AddData interface.
//
// Parameters:
//   - c: The context object containing the HTTP request/response data.
//
// Returns:
//   - error: Returns nil if successful, otherwise returns an error response.
func ProcessData[T AddData](c *core.Ctx) error {
	// Receive request data
	var requestData T
	if errData := Parse(c, &requestData); errData != nil {
		return c.Error(errData)
	}

	// Sanitize request data
	security.SanitizeStruct(&requestData)

	// Validate DTO
	if errData := Validate(requestData); errData != nil {
		return c.Error(errData)
	}

	// Store data into context
	c.SetData(constants.Request, requestData)

	return nil
}
