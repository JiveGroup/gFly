# Response Data Guide

This guide explains how to structure and return HTTP responses in the gFly/ThietNgon application. It covers response types, transformers, helper functions, and best practices.

## Table of Contents

1. [Response Structure Overview](#response-structure-overview)
2. [Generic Response Types](#generic-response-types)
3. [Custom Response Types](#custom-response-types)
4. [Transformers](#transformers)
5. [Helper Functions](#helper-functions)
6. [Sending Responses](#sending-responses)
7. [Best Practices](#best-practices)
8. [Examples](#examples)

---

## Response Structure Overview

All HTTP responses in the application follow a consistent structure. The response system is built on these principles:

- **Type Safety**: Using Go generics for reusable response structures
- **Consistency**: All responses follow the same pattern across the API
- **Documentation**: Swagger annotations on all response structs
- **Separation of Concerns**: Domain models are never returned directly; they are transformed to response DTOs

**Locations**:
- **Generic/Reusable Responses**: `github.com/gflydev/http` (package `http`) - Contains `Meta`, `List[T]`, `Success`, `Error`
- **Application-Specific Responses**: `internal/http/response/` (package `response`) - Contains domain-specific response types like `User`, `ServerInfo`

---

## Generic Response Types

Generic response types are defined in `github.com/gflydev/http` package.

### 1. List Response (`http.List[T]`)

Used for paginated list endpoints. Contains metadata and a generic data array.

**Structure**:
```go
type List[T any] struct {
    Meta Meta `json:"meta"`
    Data []T  `json:"data"`
}
```

**Metadata Structure**:
```go
type Meta struct {
    Page    int `json:"page,omitempty"`     // Current page number (starts from 1)
    PerPage int `json:"per_page,omitempty"` // Number of items per page
    Total   int `json:"total"`              // Total number of records
}
```

**When to Use**:
- List/index endpoints
- Search results
- Any endpoint returning multiple records with pagination

**Example JSON Response**:
```json
{
  "meta": {
    "page": 1,
    "per_page": 10,
    "total": 1354
  },
  "data": [
    { "id": 1, "name": "Item 1" },
    { "id": 2, "name": "Item 2" }
  ]
}
```

---

### 2. Success Response (`http.Success`)

Used for operations that succeed with an optional data payload and message.

**Structure**:
```go
type Success struct {
    Message string    `json:"message"`
    Data    core.Data `json:"data"` // Optional, can be any type
}
```

**When to Use**:
- Delete operations
- Update operations (when not returning the updated entity)
- Any operation that needs a success message
- Operations with optional return data

**Example JSON Response**:
```json
{
  "message": "User deleted successfully",
  "data": {
    "deleted_id": 42
  }
}
```

---

### 3. Error Response (`http.Error`)

Used for all error responses across the application.

**Structure**:
```go
type Error struct {
    Code    string    `json:"code"`    // Error code (e.g., "BAD_REQUEST")
    Message string    `json:"message"` // Human-readable error description
    Data    core.Data `json:"data"`    // Optional, useful for validation errors
}
```

**When to Use**:
- Validation errors (with `Data` field containing validation details)
- Business logic errors
- Not found errors
- Unauthorized/Forbidden errors
- Any error condition

**Example JSON Response** (Validation Error):
```json
{
  "code": "VALIDATION_ERROR",
  "message": "Invalid input",
  "data": {
    "email": ["Email is required", "Email must be valid"],
    "password": ["Password must be at least 8 characters"]
  }
}
```

**Example JSON Response** (General Error):
```json
{
  "code": "NOT_FOUND",
  "message": "User not found"
}
```

---

## Custom Response Types

### User Response (`response.User`)

Defined in `internal/http/response/user_response.go:10-25`.

**Structure**:
```go
type User struct {
    ID           int              `json:"id"`
    Email        string           `json:"email"`
    Fullname     string           `json:"fullname"`
    Phone        string           `json:"phone"`
    Token        *string          `json:"token"`
    Status       types.UserStatus `json:"status"`
    CreatedAt    time.Time        `json:"created_at"`
    UpdatedAt    time.Time        `json:"updated_at"`
    VerifiedAt   *time.Time       `json:"verified_at"`
    BlockedAt    *time.Time       `json:"blocked_at"`
    DeletedAt    *time.Time       `json:"deleted_at"`
    LastAccessAt *time.Time       `json:"last_access_at"`
    Avatar       *string          `json:"avatar"`
    Roles        []Role           `json:"roles"`
}
```

**List Variant** (`user_response.go:34-37`):
```go
type ListUser struct {
    Meta http.Meta `json:"meta"`  // Using http.Meta from pkg/http
    Data []User     `json:"data"`
}
```

### Server Info Response (`response.ServerInfo`)

Used for system/health check endpoints.

**Structure** (`info_response.go:9-13`):
```go
type ServerInfo struct {
    Name   string `json:"name"`   // API name (e.g., "ThietNgon API")
    Prefix string `json:"prefix"` // API prefix (e.g., "/api/v1")
    Server string `json:"server"` // Server name (e.g., "ThietNgon-Go Server")
}
```

### Creating Custom Responses

When creating new custom response types:

1. **Define in `internal/http/response/`**:
   ```go
   package response

   type Product struct {
       ID          int       `json:"id" doc:"Product unique identifier"`
       Name        string    `json:"name" doc:"Product name"`
       Price       float64   `json:"price" doc:"Product price"`
       Description string    `json:"description" doc:"Product description"`
       CreatedAt   time.Time `json:"created_at" doc:"Creation timestamp"`
   }
   ```

2. **Add list variant if needed**:
   ```go
   import "github.com/gflydev/http"

   type ListProduct struct {
       Meta http.Meta `json:"meta" doc:"Pagination metadata"`
       Data []Product `json:"data" doc:"List of products"`
   }
   ```

3. **Use `doc` tags** for Swagger documentation on each field

---

## Transformers

Transformers convert domain models to response DTOs. This ensures separation between database entities and API responses.

**Locations**:
- **Generic Transformer**: `github.com/gflydev/http` package
- **Application-Specific Transformers**: `internal/http/transformers/` (package `transformers`)

### Generic Transformer

**ToListResponse** - Transform a list of models to response DTOs:

```go
func ToListResponse[T any, R any](records []T, transformerFn func(T) R) []R
```

**Usage Example** (`list_users_api.go:60`):
```go
import "github.com/gflydev/http"

// Transform list using generic function
data := http.ToListResponse(users, transformers.ToUserResponse)
```

### User Transformers (`user_transformer.go`)

**ToUserResponse** - Convert `models.User` to `response.User`:

```go
func ToUserResponse(user models.User) response.User
```

**Features**:
- Converts null database fields to Go pointers
- Transforms avatar paths to public URLs via `PublicAvatar()`
- Loads user roles via `roles()` helper
- Handles all timestamp conversions

**ToSignUpResponse** - Similar to `ToUserResponse`, used specifically for signup responses:

```go
func ToSignUpResponse(user models.User) response.User
```

**Helper Functions**:

- **PublicAvatar** (`user_transformer.go:20-33`): Converts avatar file path to public URL
  ```go
  func PublicAvatar(avatar string) *string
  ```

- **roles** (`user_transformer.go:57-65`): Fetches and transforms user roles
  ```go
  func roles(userID int) []response.Role
  ```

- **ToRoleResponse** (`user_transformer.go:42-48`): Transforms role model to response
  ```go
  func ToRoleResponse(model models.Role) response.Role
  ```

### Creating Custom Transformers

**Pattern**:
```go
package transformers

import (
    "gfly/internal/domain/models"
    "gfly/internal/http/response"
)

// ToProductResponse converts a Product model to a Product response object
//
// Parameters:
//   - product: models.Product - The product model to convert
//
// Returns:
//   - response.Product: The converted product response object
func ToProductResponse(product models.Product) response.Product {
    return response.Product{
        ID:          product.ID,
        Name:        product.Name,
        Price:       product.Price,
        Description: product.Description,
        CreatedAt:   product.CreatedAt,
    }
}
```

**Best Practices**:
- Always add GoDoc comments with Parameters and Returns sections
- Handle null database fields properly using `dbNull` package
- Transform internal IDs, paths, enums to public-facing formats
- Load related data if needed (like roles for users)
- Keep transformers pure (no business logic)

---

## Helper Functions

Located in `github.com/gflydev/http` package, these utilities simplify common controller tasks.

### 1. PathID - Extract ID from URL Path

**Signature**:
```go
func PathID(c *core.Ctx, idName ...string) (int, *http.Error)
```

**Usage**:
```go
// Extract "id" from path
id, errResp := http.PathID(c)
if errResp != nil {
    return c.Error(*errResp)
}

// Extract custom parameter name
categoryID, errResp := http.PathID(c, "category_id")
```

**Features**:
- Validates ID is positive integer
- Returns structured error response if invalid
- Supports custom parameter names

---

### 2. Parse - Parse Request Body

**Signature**:
```go
func Parse[T any](c *core.Ctx, structData *T) *http.Error
```

**Usage**:
```go
var req request.CreateUser
if errResp := http.Parse(c, &req); errResp != nil {
    return c.Error(*errResp)
}
```

**Features**:
- Generic type-safe parsing
- Returns structured error on parse failure

---

### 3. FilterData - Extract Pagination & Filter Parameters

**Signature**:
```go
func FilterData(c *core.Ctx) dto.Filter
```

**Extracted Query Parameters**:
- `page` (default: 1)
- `per_page` (default: 10)
- `keyword` (search keyword)
- `order_by` (sort field)

**Usage** (from controller):
```go
filterDto := http.FilterData(c)
users, total, err := services.FindUsers(filterDto)
```

**Returns**:
```go
type Filter struct {
    Keyword string
    OrderBy string
    Page    int
    PerPage int
}
```

---

### 4. Validate - Perform Input Validation

**Signature**:
```go
func Validate(structData any, msgForTagFunc ...validation.MsgForTagFunc) *http.Error
```

**Usage**:
```go
requestData := request.CreateUser{
    Email: "invalid-email",
}

if errResp := http.Validate(requestData); errResp != nil {
    return c.Error(*errResp)
    // Returns: {"message": "Invalid input", "data": {"email": ["Invalid email format"]}}
}
```

**Features**:
- Uses `gflydev/validation` for validation rules
- Returns validation errors in `http.Error.Data` field
- Supports custom validation messages

---

## Sending Responses

The `core.Ctx` object provides several methods for sending responses.

### 1. Success Response (`c.Success()`)

**Usage**:
```go
import "github.com/gflydev/http"

return c.Success(http.Success{
    Message: "User deleted successfully",
    Data: core.Data{"deleted_id": userID},
})
```

**HTTP Status**: 200 OK

---

### 2. JSON Response (`c.JSON()`)

Send any data as JSON with default 200 status.

**Usage**:
```go
return c.JSON(userResponse)
```

**HTTP Status**: 200 OK

---

### 3. JSON with Custom Status (`c.Status().JSON()`)

**Usage** (`create_user_api.go:62-64`):
```go
return c.
    Status(core.StatusCreated).
    JSON(userResponse)
```

**HTTP Status**: 201 Created (or any custom status)

---

### 4. Error Response (`c.Error()`)

**Usage**:
```go
import "github.com/gflydev/http"

return c.Error(http.Error{
    Code:    "NOT_FOUND",
    Message: "User not found",
})
```

**HTTP Status**: Default error status (usually 400)

**With Custom Status**:
```go
import "github.com/gflydev/http"

return c.
    Status(core.StatusNotFound).
    Error(http.Error{
        Code:    "NOT_FOUND",
        Message: "User not found",
    })
```

---

### 5. No Content Response (`c.NoContent()`)

For operations that don't return data.

**Usage** (`signin_api.go:68`):
```go
if h.Type == auth.TypeWeb {
    c.SetSession(auth.SessionUsername, requestData.ToDto().Username)
    return c.NoContent()
}
```

**HTTP Status**: 204 No Content

---

## Best Practices

### 1. Always Use Transformers

**Bad** (Never do this):
```go
// DON'T return domain models directly
return c.JSON(user) // models.User
```

**Good**:
```go
// DO transform to response DTO
userResponse := transformers.ToUserResponse(user)
return c.JSON(userResponse)
```

**Why**: Separation of concerns, API stability, security (no accidental exposure of internal fields)

---

### 2. Use Consistent Status Codes

| Operation | Status Code | Response Type |
|-----------|-------------|---------------|
| List/Get success | 200 OK | `http.List[T]` or custom DTO |
| Create success | 201 Created | Custom DTO |
| Update success | 200 OK | Custom DTO |
| Delete success | 200 OK | `http.Success` |
| No content | 204 No Content | None |
| Validation error | 400 Bad Request | `http.Error` |
| Unauthorized | 401 Unauthorized | `http.Error` |
| Forbidden | 403 Forbidden | `http.Error` |
| Not found | 404 Not Found | `http.Error` |
| Server error | 500 Internal Server Error | `http.Error` |

---

### 3. Include Swagger Annotations

**Example** (`create_user_api.go:37-48`):
```go
// Handle function allows Administrator create a new user with specific roles
// @Description Function allows Administrator create a new user with specific roles
// @Summary Create a new user for Administrator
// @Tags Users
// @Accept json
// @Produce json
// @Param data body request.CreateUser true "CreateUser payload"
// @Success 201 {object} response.User
// @Failure 400 {object} http.Error
// @Failure 401 {object} http.Error
// @Security ApiKeyAuth
// @Router /users [post]
func (h *CreateUserApi) Handle(c *core.Ctx) error {
    // ...
}
```

**Required Annotations**:
- `@Summary` - Brief description
- `@Description` - Detailed description
- `@Tags` - API grouping
- `@Success` - Success response type
- `@Failure` - All possible error responses
- `@Router` - Route path and method

---

### 4. Handle Errors Properly

**Always check and handle errors**:
```go
import "github.com/gflydev/http"

user, err := services.GetUser(id)
if err != nil {
    return c.Error(http.Error{
        Message: err.Error(),
    })
}
```

**With custom status codes**:
```go
import "github.com/gflydev/http"

user, err := services.GetUser(id)
if err != nil {
    return c.
        Status(core.StatusNotFound).
        Error(http.Error{
            Code:    "USER_NOT_FOUND",
            Message: "User not found",
        })
}
```

---

### 5. Use Generic List Response

For any list endpoint, use `http.List[T]` pattern:

```go
import "github.com/gflydev/http"

metadata := http.Meta{
    Page:    filterDto.Page,
    PerPage: filterDto.PerPage,
    Total:   total,
}

data := http.ToListResponse(products, transformers.ToProductResponse)

return c.Success(http.List[response.Product]{
    Meta: metadata,
    Data: data,
})
```

Or define a custom list type like `response.ListUser`:
```go
import "github.com/gflydev/http"

type ListProduct struct {
    Meta http.Meta `json:"meta"`
    Data []Product `json:"data"`
}
```

---

## Examples

### Example 1: List Endpoint

**Controller** (`list_users_api.go:46-67`):
```go
func (h *ListUsersApi) Handle(c *core.Ctx) error {
    // Get filter from context (set by middleware)
    filterDto := c.GetData(constants.Filter).(dto.Filter)

    // Fetch data from service
    users, total, err := services.FindUsers(filterDto)
    if err != nil {
        return err
    }

    // Build metadata
    metadata := response.Meta{
        Page:    filterDto.Page,
        PerPage: filterDto.PerPage,
        Total:   total,
    }

    // Transform to response data
    data := transformers.ToListResponse(users, transformers.ToUserResponse)

    // Return response
    return c.Success(response.ListUser{
        Meta: metadata,
        Data: data,
    })
}
```

**Response**:
```json
{
  "meta": {
    "page": 1,
    "per_page": 10,
    "total": 45
  },
  "data": [
    {
      "id": 1,
      "email": "user@example.com",
      "fullname": "John Doe",
      "phone": "1234567890",
      "status": "active",
      "roles": [
        {"id": 1, "name": "User", "slug": "user"}
      ]
    }
  ]
}
```

---

### Example 2: Create Endpoint

**Controller** (`create_user_api.go:49-65`):
```go
func (h *CreateUserApi) Handle(c *core.Ctx) error {
    // Get validated request data from context
    requestData := c.GetData(constants.Request).(request.CreateUser)

    // Create user via service
    user, err := services.CreateUser(requestData.ToDto())
    if err != nil {
        return c.Error(response.Error{
            Message: err.Error(),
        })
    }

    // Transform to response data
    userResponse := transformers.ToUserResponse(*user)

    // Return 201 Created
    return c.
        Status(core.StatusCreated).
        JSON(userResponse)
}
```

**Response** (201 Created):
```json
{
  "id": 42,
  "email": "newuser@example.com",
  "fullname": "Jane Smith",
  "phone": "9876543210",
  "status": "pending",
  "roles": [
    {"id": 2, "name": "Admin", "slug": "admin"}
  ],
  "created_at": "2025-10-23T10:30:00Z",
  "updated_at": "2025-10-23T10:30:00Z"
}
```

---

### Example 3: Delete Endpoint

**Controller**:
```go
import "github.com/gflydev/http"

func (h *DeleteUserApi) Handle(c *core.Ctx) error {
    // Extract ID from path
    id, errResp := http.PathID(c)
    if errResp != nil {
        return c.Error(*errResp)
    }

    // Delete via service
    err := services.DeleteUser(id)
    if err != nil {
        return c.Error(http.Error{
            Message: err.Error(),
        })
    }

    // Return success message
    return c.Success(http.Success{
        Message: "User deleted successfully",
        Data:    core.Data{"deleted_id": id},
    })
}
```

**Response** (200 OK):
```json
{
  "message": "User deleted successfully",
  "data": {
    "deleted_id": 42
  }
}
```

---

### Example 4: Validation Error

**Controller**:
```go
func (h *CreateUserApi) Validate(c *core.Ctx) error {
    return http.ProcessData[request.CreateUser](c)
}
```

**Invalid Request**:
```json
{
  "email": "invalid-email",
  "password": "123"
}
```

**Response** (400 Bad Request):
```json
{
  "code": "VALIDATION_ERROR",
  "message": "Invalid input",
  "data": {
    "email": ["Email must be a valid email address"],
    "password": ["Password must be at least 8 characters"],
    "fullname": ["Fullname is required"]
  }
}
```

---

### Example 5: Sign In Endpoint

**Controller** (`signin_api.go:54-72`):
```go
func (h *SignInApi) Handle(c *core.Ctx) error {
    // Get valid data from context
    requestData := c.GetData(constants.Request).(request.SignIn)

    // Authenticate user
    tokens, err := services.SignIn(requestData.ToDto())
    if err != nil {
        return c.Error(httpResponse.Error{
            Message: err.Error(),
        })
    }

    // For web-based auth, use session
    if h.Type == auth.TypeWeb {
        c.SetSession(auth.SessionUsername, requestData.ToDto().Username)
        return c.NoContent()
    }

    // For API auth, return tokens
    return c.JSON(transformers.ToSignInResponse(tokens))
}
```

**API Response** (200 OK):
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
  "expires_in": 3600
}
```

**Web Response** (204 No Content):
```
(empty body, session cookie set)
```

---

## Summary

### Key Takeaways

1. **Use generic response types** (`http.List[T]`, `http.Success`, `http.Error` from `github.com/gflydev/http`) for consistency
2. **Always transform domain models** using transformers before returning
3. **Include complete Swagger annotations** on all controller handlers
4. **Use helper functions** (`http.PathID`, `http.Validate`, `http.FilterData`) for common tasks
5. **Return appropriate HTTP status codes** based on operation type
6. **Handle errors properly** and return structured error responses
7. **Keep transformers pure** - no business logic, just data transformation
8. **Document everything** - GoDoc comments, Swagger annotations, `doc` tags

### Response Flow

```
Domain Model (models.User)
    ↓
Transformer (transformers.ToUserResponse)
    ↓
Response DTO (response.User)
    ↓
Controller Method (c.JSON/c.Success/c.Error)
    ↓
HTTP Response (JSON)
```

### Quick Reference

| Need | Use |
|------|-----|
| Return a list | `http.List[T]` or custom `ListX` type |
| Return single item | Custom response type + transformer |
| Return success message | `http.Success` |
| Return error | `http.Error` |
| Parse request body | `http.Parse[T]()` |
| Get path ID | `http.PathID()` |
| Validate input | `http.Validate()` |
| Get filter params | `http.FilterData()` |
| Transform model | `transformers.ToXResponse()` |
| Transform list | `http.ToListResponse()` |

---

For additional information, see:
- `github.com/gflydev/http` - Generic response types (Meta, List, Success, Error), transformers (ToListResponse), and helper functions (PathID, Parse, Validate, FilterData)
- `internal/http/response/` - Application-specific response type definitions (User, ServerInfo, etc.)
- `internal/http/transformers/` - Application-specific transformer implementations (ToUserResponse, etc.)
- `CLAUDE.md` - Project overview and development workflow
