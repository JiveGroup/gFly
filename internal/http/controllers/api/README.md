# API Controllers

This directory contains controllers that handle API requests and return JSON responses.

## Purpose

API controllers are responsible for:
- Processing RESTful API requests
- Validating API input data
- Interacting with services to perform business logic
- Returning standardized JSON responses
- Handling API-specific errors and exceptions
- Implementing API versioning

## Structure

- **default_api.go**: Default API controller with basic endpoints

## Usage

API controllers should follow RESTful conventions and return consistent JSON responses:

```
// Example pseudocode for an API controller
// Note: This is not actual implementation code

// In an API controller file (user_api_controller.go):
package api

import (
    "gfly/internal/http/request"
    "gfly/internal/http/response"
    "gfly/internal/services"
    "gfly/internal/http/transformers"
)

// UserAPIController handles user-related API requests
type UserAPIController struct {
    userService *services.UserService
}

// NewUserAPIController creates a new user API controller
func NewUserAPIController(userService *services.UserService) *UserAPIController {
    return &UserAPIController{
        userService: userService,
    }
}

// List handles GET /api/users requests
func (c *UserAPIController) List(ctx *fasthttp.RequestCtx) {
    // Get query parameters
    page := ctx.QueryArgs().GetUintOrZero("page")
    limit := ctx.QueryArgs().GetUintOrZero("limit")

    // Get users from service
    users, pagination, err := c.userService.ListUsers(page, limit)
    if err != nil {
        response.Error(ctx, err)
        return
    }

    // Transform and return response
    data := transformers.TransformUsers(users)
    response.SuccessWithPagination(ctx, data, pagination)
}
```

## Best Practices

- Follow RESTful API conventions
- Use appropriate HTTP methods (GET, POST, PUT, DELETE)
- Return consistent JSON response structures
- Include proper HTTP status codes
- Implement API versioning
- Document API endpoints with Swagger or similar tools
- Validate all input data
- Use transformers to format response data
- Handle errors with appropriate status codes and messages
- Implement proper authentication and authorization
- Consider rate limiting for public APIs
- Use pagination for list endpoints
- Support filtering, sorting, and searching where appropriate
