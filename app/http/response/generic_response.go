package response

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

// ServerInfo struct to describe system information.
// @Description contains system metadata including name, server prefix, and server name.
// @Name Name is the name of the API.
// @Prefix Prefix is the API prefix including version.
// @Server Server is the name of the server application.
// @Tags Success Responses
type ServerInfo struct {
	Name   string `json:"name" example:"ThietNgon API" doc:"API name"`
	Prefix string `json:"prefix" example:"/api/v1" doc:"API prefix including version"`
	Server string `json:"server" example:"ThietNgon-Go Server" doc:"Server application name"`
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
	Code    int       `json:"code" example:"400"`            // HTTP status code
	Message string    `json:"message" example:"Bad request"` // Error message description
	Data    core.Data `json:"data"`                          // Useful for validation's errors
}

// Unauthorized clone from app.core.errors.Unauthorized
// @Description Unauthorized error response structure
// @Code Code is the HTTP status code for the error.
// @Message Message is a description of the error that occurred.
// @Tags Error Responses
type Unauthorized struct {
	Code    int    `json:"code" example:"401"`                  // HTTP status code
	Message string `json:"error" example:"Unauthorized access"` // Error message description
}

// NotFound handle not found any record
// @Description Not found error response structure
// @Code Code is the HTTP status code for the error.
// @Message Message is a description of the error that occurred.
// @Tags Error Responses
type NotFound struct {
	Code    int    `json:"code" example:"404"`                 // HTTP status code
	Message string `json:"error" example:"Resource not found"` // Error message description
}

// Conflict describes a conflict error
// @Description Conflict error response structure
// @Code Code is the HTTP status code for the error.
// @Message Message is a description of the error that occurred.
// @Tags Error Responses
type Conflict struct {
	Code    int    `json:"code" example:"409"`                // HTTP status code
	Message string `json:"error" example:"Resource conflict"` // Error message description
}
