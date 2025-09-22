# Controllers

This directory contains controller classes that handle HTTP requests and return appropriate responses.

## Purpose

Controllers are responsible for:
- Receiving and processing HTTP requests
- Validating input data
- Interacting with services to perform business logic
- Returning appropriate HTTP responses
- Handling errors and exceptions

## Structure

- **api/**: Controllers for RESTful API endpoints
- **page/**: Controllers for web pages and server-side rendering

## Usage

Controllers should be organized by their domain or functionality:

```
// Example pseudocode for a controller
// Note: This is not actual implementation code

// In a controller file (user_controller.go):
package controllers

import (
    "gfly/internal/http/request"
    "gfly/internal/http/response"
    "gfly/internal/services"
)

// UserController handles user-related requests
type UserController struct {
    userService *services.UserService
}

// NewUserController creates a new user controller
func NewUserController(userService *services.UserService) *UserController {
    return &UserController{
        userService: userService,
    }
}

// GetUser handles GET /users/:id requests
func (c *UserController) GetUser(ctx *fasthttp.RequestCtx) {
    // Get user ID from request
    userID := ctx.UserValue("id").(string)

    // Get user from service
    user, err := c.userService.GetUserByID(userID)
    if err != nil {
        // Handle error
        response.Error(ctx, err)
        return
    }

    // Return success response
    response.Success(ctx, user)
}
```

## Best Practices

- Keep controllers thin, delegating business logic to services
- Use dependency injection for services and other dependencies
- Follow RESTful conventions for API controllers
- Handle errors appropriately and return consistent error responses
- Validate all input data using request validation
- Use transformers to format response data
- Group related controller methods together
- Document controller methods with clear descriptions
- Test controllers with both valid and invalid inputs
- Consider using controller structs for better organization and dependency management
