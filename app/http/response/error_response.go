package response

import "github.com/gflydev/core"

// Error struct to describe login response.
// @Description Generic error response structure
// @Tags Error Responses
type Error struct {
	Code    int       `json:"code" example:"400"` // HTTP status code
	Message string    `json:"message"`            // Error message description
	Data    core.Data `json:"data"`               // Useful for validation's errors
}

// Unauthorized clone from app.core.errors.Unauthorized
// @Description Unauthorized error response structure
// @Tags Error Responses
type Unauthorized struct {
	Code    int    `json:"code" example:"401"` // HTTP status code
	Message string `json:"error"`              // Error message description
}

// NotFound handle not found any record
// @Description Not found error response structure
// @Tags Error Responses
type NotFound struct {
	Code    int    `json:"code" example:"404"` // HTTP status code
	Message string `json:"error"`              // Error message description
}

// Conflict describes a conflict error
// @Description Conflict error response structure
// @Tags Error Responses
type Conflict struct {
	Code    int    `json:"code" example:"409"` // HTTP status code
	Message string `json:"error"`              // Error message description
}
