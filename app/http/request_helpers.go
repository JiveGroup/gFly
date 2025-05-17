package http

import (
	"gfly/app/constants"
	"github.com/gflydev/core"
)

// ====================================================================
// ======================== Other Request Helpers =====================
// ====================================================================

// ProcessPathID is a generic function that extracts a path ID parameter and stores it in the context
// It handles the common pattern of validating a path ID parameter for API endpoints
func ProcessPathID(c *core.Ctx) error {
	// Receive path parameter ID
	itemID, errData := PathID(c)
	if errData != nil {
		return c.Error(errData)
	}

	// Store data into context
	c.SetData(constants.Data, itemID)

	return nil
}

// ====================================================================
// ======================= Update Request Helpers =====================
// ====================================================================

// UpdateRequest is an interface for types that can convert to a DTO
// It defines the contract for request types that need ID setting and DTO conversion capabilities
type UpdateRequest[D any] interface {
	// SetID sets the ID field of the request structure
	// Parameters:
	//   - id: Integer ID value to set
	SetID(int)

	// Request converts the request to its corresponding DTO type
	// Returns:
	//   - D: The converted DTO object of generic type D
	Request[D]
}

// ProcessUpdateRequest validates and processes update requests that require a path ID parameter
// It handles parsing the request body, setting the ID, converting to DTO, and validation
//
// Type Parameters:
//   - T: Request type that implements UpdateRequest interface
//   - D: Target DTO type that the request converts to
//
// Parameters:
//   - c: The context object containing the HTTP request/response data
//
// Returns:
//   - error: Returns nil if successful, otherwise returns an error response
//
// Example Usage:
//
//	func (h UpdateUserApi) Validate(c *core.Ctx) error {
//		return http.ProcessUpdateRequest[request.UpdateUser, dto.UpdateUser](c)
//	}
func ProcessUpdateRequest[T UpdateRequest[D], D any](c *core.Ctx) error {
	// Receive path parameter ID
	itemID, errData := PathID(c)
	if errData != nil {
		return c.Error(errData)
	}

	// Receive request data
	var requestBody T
	if errData := Parse(c, &requestBody); errData != nil {
		return c.Error(errData)
	}

	// Set ID on request body
	requestBody.SetID(itemID)

	// Convert to DTO
	requestDto := requestBody.ToDto()

	// Validate DTO
	if errData := Validate(requestDto); errData != nil {
		return c.Error(errData)
	}

	// Store data into context
	c.SetData(constants.Data, requestDto)

	return nil
}

// ====================================================================
// ======================== Add Request Helpers =======================
// ====================================================================

// Request is an interface for types that can convert to a DTO
// It defines the contract for request types that need DTO conversion capabilities
type Request[D any] interface {
	// ToDto converts the request to its corresponding DTO type
	// Returns:
	//   - D: The converted DTO object of generic type D
	ToDto() D
}

// ProcessRequest validates and processes create/add requests
// It handles parsing the request body, converting to DTO, and validation
//
// Type Parameters:
//   - T: Request type that implements Request interface
//   - D: Target DTO type that the request converts to
//
// Parameters:
//   - c: The context object containing the HTTP request/response data
//
// Returns:
//   - error: Returns nil if successful, otherwise returns an error response
//
// Example Usage:
//
//	func (h CreateUserApi) Validate(c *core.Ctx) error {
//		return http.ProcessRequest[request.CreateUser, dto.CreateUser](c)
//	}
func ProcessRequest[T Request[D], D any](c *core.Ctx) error {
	// Receive request data
	var requestBody T
	if errData := Parse(c, &requestBody); errData != nil {
		return c.Error(errData)
	}

	// Convert to DTO
	requestDto := requestBody.ToDto()

	// Validate DTO
	if errData := Validate(requestDto); errData != nil {
		return c.Error(errData)
	}

	// Store data into context
	c.SetData(constants.Data, requestDto)

	return nil
}
