package http

import "github.com/gflydev/core"

// ====================================================================
// ======================== Success Responses =========================
// ====================================================================

// Meta struct to describe pagination metadata information.
// @Description Contains pagination metadata including current page, items per page, and total count
// @Page Page is the current page number (optional, starts from 1)
// @PerPage PerPage is the number of items displayed per page (optional)
// @Total Total is the total number of records available
// @Tags Info Responses
type Meta struct {
	Page    int `json:"page,omitempty" example:"1" doc:"Current page number"`
	PerPage int `json:"per_page,omitempty" example:"10" doc:"Number of items per page"`
	Total   int `json:"total" example:"1354" doc:"Total number of records"`
}

// List struct to describe a generic list response.
// @Description Generic list response structure
// @Meta Meta contains metadata information for pagination.
// @Data Data is a slice of type T, which can be any data type.
// @Tags Success Responses
type List[T any] struct {
	Meta Meta `json:"meta" example:"{\"page\":1,\"per_page\":10,\"total\":100}" doc:"Metadata information for pagination"`
	Data []T  `json:"data" example:"[]" doc:"List of category data"`
}

// Success struct to describe a generic success response.
// @Description Generic success response structure
// @Data Data is optional and can be used to return additional information related to the operation.
// @Message Message is a success message that describes the operation.
// @Tags Success Responses
type Success struct {
	Message string    `json:"message" example:"Operation completed successfully"`  // Success message description
	Data    core.Data `json:"data" doc:"Additional data related to the operation"` // Optional data related to the success operation
}

// ====================================================================
// ========================= Error Responses ==========================
// ====================================================================

// Error struct to describe login response.
// @Description Generic error response structure
// @Data Data is optional and can be used to return additional information related to the operation.
// @Code Code is the HTTP status code for the error.
// @Message Message is a description of the error that occurred.
// @Tags Error Responses
type Error struct {
	Code    string    `json:"code" example:"BAD_REQUEST"`    // Error code
	Message string    `json:"message" example:"Bad request"` // Error message description
	Data    core.Data `json:"data"`                          // Useful for validation's errors
}
