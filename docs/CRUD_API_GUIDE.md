# gFly Framework: CRUD API Development Guide

## Table of Contents
1. [CRUD Overview](#crud-overview)
2. [HTTP Helper Functions](#http-helper-functions)
3. [Create (POST) API](#create-post-api)
4. [Read (GET Single) API](#read-get-single-api)
5. [Update (PUT/PATCH) API](#update-putpatch-api)
6. [Delete (DELETE) API](#delete-delete-api)
7. [List (GET Collection) API](#list-get-collection-api)
8. [Advanced: Custom Filters with Multiple Parameters](#advanced-custom-filters-with-multiple-parameters)
9. [Response Patterns](#response-patterns)
10. [HTTP Status Codes](#http-status-codes)
11. [Error Handling](#error-handling)
12. [Complete CRUD Example](#complete-crud-example)
13. [Best Practices](#best-practices)

---

## CRUD Overview

CRUD stands for **Create, Read, Update, Delete** - the four basic operations for persistent storage. In gFly, each operation follows a consistent pattern:

```
HTTP Request → Controller.Validate() → Controller.Handle() → Service → Repository → Database
                                                             ↓
HTTP Response ← Transformer ← Domain Model ← Service ← Repository ← Database
```

### Standard CRUD Endpoints

| Operation | HTTP Method | Endpoint Pattern | Example |
|-----------|-------------|------------------|---------|
| **Create** | POST | `/resource` | `POST /api/v1/users` |
| **Read** | GET | `/resource/:id` | `GET /api/v1/users/123` |
| **Update** | PUT/PATCH | `/resource/:id` | `PUT /api/v1/users/123` |
| **Delete** | DELETE | `/resource/:id` | `DELETE /api/v1/users/123` |
| **List** | GET | `/resource` | `GET /api/v1/users?page=1` |

---

## HTTP Helper Functions

gFly provides a set of helper functions in `/pkg/http/` for common request processing tasks.

### Request Processing Helpers

```go
// File: pkg/http/request_helpers.go

// ProcessData validates and processes CREATE requests (POST)
// Type parameter T must implement AddData interface
func ProcessData[T AddData](c *core.Ctx) error

// ProcessUpdateData validates and processes UPDATE requests (PUT/PATCH)
// Type parameter T must implement UpdateData interface
func ProcessUpdateData[T UpdateData](c *core.Ctx) error

// ProcessFilter validates and processes LIST requests (GET with query params)
func ProcessFilter(c *core.Ctx) error

// ProcessPathID validates and processes READ/DELETE requests (GET/DELETE with path ID)
// Extracts ID from URL path, validates it, and stores in context
func ProcessPathID(c *core.Ctx) error
```

### Request Interfaces

```go
// AddData interface for create operations
type AddData[D any] interface {
    ToDto() D         // Convert request to DTO
}

// UpdateData interface for update operations
type UpdateData[D any] interface {
    ToDto() D         // Convert request to DTO
    SetID(id int)     // Set ID from URL path parameter
}
```

### Core Helpers

```go
// File: pkg/http/http_helpers.go

// Parse parses JSON request body into a struct
func Parse(c *core.Ctx, requestData any) *response.Error

// Validate performs validation using struct tags
func Validate(structData any, msgForTagFunc ...validation.MsgForTagFunc) *response.Error

// PathID extracts and validates ID from URL path parameter
// NOTE: For controllers, prefer using ProcessPathID() instead
func PathID(c *core.Ctx) (int, *response.Error)

// FilterData retrieves filter parameters: page, per_page, keyword, order_by
func FilterData(c *core.Ctx) dto.Filter
```

**When to use each helper:**

| Helper | Use Case | Example |
|--------|----------|---------|
| `ProcessData[T]` | CREATE operations (POST with JSON body) | User registration, creating products |
| `ProcessUpdateData[T]` | UPDATE operations (PUT/PATCH with path ID + JSON body) | Updating user profile, editing products |
| `ProcessFilter` | LIST operations (GET with query params) | Paginated lists, search with filters |
| `ProcessPathID` | READ/DELETE operations (GET/DELETE with path ID only) | Get user by ID, delete product |
| `PathID` | When you need just the ID value (rare, prefer `ProcessPathID`) | Custom validation logic |

### Usage in Controllers

```go
// CREATE operation
func (h *CreateApi) Validate(c *core.Ctx) error {
    return http.ProcessData[request.CreateUser](c) // Extracts Data, validates, stores in context
}

// UPDATE operation
func (h *UpdateApi) Validate(c *core.Ctx) error {
    return http.ProcessUpdateData[request.UpdateUser](c) // Extracts Data, validates, stores in context
}

// LIST operation
func (h *ListApi) Validate(c *core.Ctx) error {
    return http.ProcessFilter(c) // Extracts Filter, validates, stores in context
}

// READ operations (ID from path)
func (h *GetUserApi) Validate(c *core.Ctx) error {
    return http.ProcessPathID(c)  // Extracts ID, validates, stores in context
}

// DELETE operations (ID from path)
func (h *DeleteApi) Validate(c *core.Ctx) error {
    return http.ProcessPathID(c)  // Extracts ID, validates, stores in context
}
```

---

## Create (POST) API

### Complete Flow

```
POST /api/v1/users
{
  "email": "john@example.com",
  "password": "SecurePass123",
  "fullname": "John Doe"
}

↓ Controller.Validate()
  → ProcessData[request.CreateUser](c)
    → Parse JSON body
    → Sanitize input
    → Validate using DTO tags
    → Store in context: c.SetData(constants.Request, requestData)

↓ Controller.Handle()
  → Get request: c.GetData(constants.Request)
  → Convert to DTO: requestData.ToDto()
  → Call service: services.CreateUser(dto)
  → Transform model: transformers.ToUserResponse(user)
  → Return 201 Created

Response: 201 Created
{
  "id": 123,
  "email": "john@example.com",
  "fullname": "John Doe",
  "created_at": "2025-01-22T10:30:00Z"
}
```

### 1. Define DTO

```go
// File: internal/dto/user_dto.go
package dto

// CreateUser defines the data structure for user creation
type CreateUser struct {
    Email    string `json:"email" validate:"required,email,max=255" doc:"User's email address"`
    Password string `json:"password" validate:"required,min=8,max=255" doc:"User's password (min 8 chars)"`
    Fullname string `json:"fullname" validate:"required,max=255" doc:"User's full name"`
    Phone    string `json:"phone" validate:"required,max=20" doc:"User's phone number"`
    Avatar   string `json:"avatar" validate:"omitempty,url,max=255" doc:"Avatar URL (optional)"`
}
```

### 2. Define Request

```go
// File: internal/http/request/user_request.go
package request

import "gfly/internal/dto"

// CreateUser wraps the CreateUser DTO
type CreateUser struct {
    dto.CreateUser  // Embedded - inherits all fields and validation
}

// ToDto converts Request to DTO
func (r CreateUser) ToDto() dto.CreateUser {
    return r.CreateUser
}
```

### 3. Define Response

```go
// File: internal/http/response/user_response.go
package response

import "time"

// User represents the user response structure
type User struct {
    ID        int       `json:"id" doc:"Unique user identifier"`
    Email     string    `json:"email" doc:"User's email address"`
    Fullname  string    `json:"fullname" doc:"User's full name"`
    Phone     string    `json:"phone" doc:"User's phone number"`
    Avatar    *string   `json:"avatar" doc:"Avatar URL (nullable)"`
    CreatedAt time.Time `json:"created_at" doc:"Account creation timestamp"`
    UpdatedAt time.Time `json:"updated_at" doc:"Last update timestamp"`
}
```

### 4. Define Transformer

```go
// File: internal/http/transformers/user_transformer.go
package transformers

import (
    "gfly/internal/domain/models"
    "gfly/internal/http/response"
    dbNull "github.com/gflydev/db/null"
)

// ToUserResponse converts User model to User response
func ToUserResponse(user models.User) response.User {
    return response.User{
        ID:        user.ID,
        Email:     user.Email,
        Fullname:  user.Fullname,
        Phone:     user.Phone,
        Avatar:    dbNull.StringNil(user.Avatar),  // sql.NullString → *string
        CreatedAt: user.CreatedAt,
        UpdatedAt: user.UpdatedAt,
    }
}
```

### 5. Implement Controller

```go
// File: internal/http/controllers/api/user/create_user_api.go
package user

import (
    "gfly/internal/constants"
    "gfly/internal/http/request"
    "gfly/internal/http/response"
    "gfly/internal/http/transformers"
    "gfly/internal/services"
    "gfly/pkg/http"
    "github.com/gflydev/core"
)

type CreateUserApi struct {
    core.Api
}

func NewCreateUserApi() *CreateUserApi {
    return &CreateUserApi{}
}

// Validate handles request parsing and validation
func (h *CreateUserApi) Validate(c *core.Ctx) error {
    // ProcessData: Parse JSON, sanitize, validate, store in context
    return http.ProcessData[request.CreateUser](c)
}

// Handle processes the create user request
// @Summary Create a new user
// @Description Creates a new user account with the provided information
// @Tags Users
// @Accept json
// @Produce json
// @Param data body request.CreateUser true "User creation data"
// @Success 201 {object} response.User "User created successfully"
// @Failure 400 {object} response.Error "Validation error or business rule violation"
// @Failure 401 {object} response.Unauthorized "Authentication required"
// @Security ApiKeyAuth
// @Router /users [post]
func (h *CreateUserApi) Handle(c *core.Ctx) error {
    // 1. Retrieve validated request from context
    requestData := c.GetData(constants.Request).(request.CreateUser)

    // 2. Convert Request → DTO
    createUserDto := requestData.ToDto()

    // 3. Call service layer with DTO
    user, err := services.CreateUser(createUserDto)
    if err != nil {
        return c.Error(response.Error{
            Code:    core.StatusBadRequest,
            Message: err.Error(),
        })
    }

    // 4. Transform Model → Response
    userResponse := transformers.ToUserResponse(*user)

    // 5. Return 201 Created with response body
    return c.
        Status(core.StatusCreated).
        JSON(userResponse)
}
```

### 6. Implement Service

```go
// File: internal/services/user_services.go
package services

import (
    "gfly/internal/domain/models"
    "gfly/internal/domain/repository"
    "gfly/internal/dto"
    mb "github.com/gflydev/db"
    "github.com/gflydev/core/errors"
    "time"
)

// CreateUser creates a new user in the system
//
// Parameters:
//   - createUserDto: User creation data
//
// Returns:
//   - (*models.User, error): Created user or error
func CreateUser(createUserDto dto.CreateUser) (*models.User, error) {
    // Business rule: Check email uniqueness
    existingUser := repository.Pool.GetUserByEmail(createUserDto.Email)
    if existingUser != nil {
        return nil, errors.New("user with this email already exists")
    }

    // Create domain model
    user := &models.User{
        Email:     createUserDto.Email,
        Password:  utils.GeneratePassword(createUserDto.Password),
        Fullname:  createUserDto.Fullname,
        Phone:     createUserDto.Phone,
        Avatar:    dbNull.String(createUserDto.Avatar),
        Status:    types.UserStatusActive,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }

    // Persist to database
    err := mb.CreateModel(user)
    if err != nil {
        log.Errorf("Failed to create user: %v", err)
        return nil, errors.New("failed to create user")
    }

    return user, nil
}
```

### 7. Register Route

```go
// File: internal/http/routes/api_routes.go

usersGroup := apiRouter.Group("/users")

// POST /api/v1/users - Create user (admin only)
usersGroup.POST("", app.Apply(
    middleware.JWTAuth(),
    middleware.CheckRolesMiddleware(types.RoleAdmin),
)(user.NewCreateUserApi()))
```

---

## Read (GET Single) API

### Complete Flow

```
GET /api/v1/users/123

↓ Controller.Validate()
  → PathID(c) extracts ID from URL
  → Validates ID is a positive integer
  → Store in context: c.SetData(constants.ItemID, 123)

↓ Controller.Handle()
  → Get ID: c.GetData(constants.ItemID)
  → Call service: services.GetUserByID(id)
  → Transform model: transformers.ToUserResponse(user)
  → Return 200 OK

Response: 200 OK
{
  "id": 123,
  "email": "john@example.com",
  "fullname": "John Doe",
  "created_at": "2025-01-22T10:30:00Z"
}
```

### Controller Implementation

```go
// File: internal/http/controllers/api/user/get_user_api.go
package user

import (
    "gfly/internal/constants"
    "gfly/internal/http/response"
    "gfly/internal/http/transformers"
    "gfly/internal/services"
    "gfly/pkg/http"
    "github.com/gflydev/core"
)

type GetUserApi struct {
    core.Api
}

func NewGetUserApi() *GetUserApi {
    return &GetUserApi{}
}

// Validate extracts and validates the user ID from URL path
func (h *GetUserApi) Validate(c *core.Ctx) error {
    // ProcessPathID: Extracts ID from path, validates, and stores in context
    return http.ProcessPathID(c)
}

// Handle retrieves a single user by ID
// @Summary Get user by ID
// @Description Retrieves detailed information about a specific user
// @Tags Users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} response.User "User found"
// @Failure 400 {object} response.Error "Invalid ID format"
// @Failure 404 {object} response.NotFound "User not found"
// @Failure 401 {object} response.Unauthorized "Authentication required"
// @Security ApiKeyAuth
// @Router /users/{id} [get]
func (h *GetUserApi) Handle(c *core.Ctx) error {
    // 1. Get ID from context
    itemID := c.GetData(constants.ItemID).(int)

    // 2. Call service to retrieve user
    user, err := services.GetUserByID(itemID)
    if err != nil {
        return c.Error(response.NotFound{
            Code:    core.StatusNotFound,
            Message: "User not found",
        })
    }

    // 3. Transform Model → Response
    userResponse := transformers.ToUserResponse(*user)

    // 4. Return 200 OK with user data
    return c.JSON(userResponse)
}
```

### Service Implementation

```go
// File: internal/services/user_services.go

// GetUserByID retrieves a user by their ID
//
// Parameters:
//   - id: User ID
//
// Returns:
//   - (*models.User, error): User or error if not found
func GetUserByID(id int) (*models.User, error) {
    user, err := mb.GetModelByID[models.User](id)
    if err != nil {
        return nil, errors.New("user not found")
    }

    return user, nil
}
```

### Route Registration

```go
// GET /api/v1/users/:id - Get user by ID
usersGroup.GET("/:id", app.Apply(
    middleware.JWTAuth(),
)(user.NewGetUserApi()))
```

---

## Update (PUT/PATCH) API

### Complete Flow

```
PUT /api/v1/users/123
{
  "fullname": "Jane Doe Updated",
  "phone": "0123456789"
}

↓ Controller.Validate()
  → ProcessUpdateData[request.UpdateUser](c)
    → Extract ID from URL path (123)
    → Parse JSON body
    → Sanitize input
    → Call requestData.SetID(123)
    → Validate using DTO tags
    → Store in context

↓ Controller.Handle()
  → Get request from context
  → Convert to DTO: requestData.ToDto()
  → Call service: services.UpdateUser(dto)
  → Transform model: transformers.ToUserResponse(user)
  → Return 200 OK

Response: 200 OK
{
  "id": 123,
  "email": "john@example.com",
  "fullname": "Jane Doe Updated",
  "phone": "0123456789",
  "updated_at": "2025-01-22T11:00:00Z"
}
```

### 1. Define DTO

```go
// File: internal/dto/user_dto.go

// UpdateUser defines partial update structure (all fields optional)
type UpdateUser struct {
    ID       int    `json:"-" validate:"omitempty,gte=1" doc:"User ID (from URL path)"`
    Password string `json:"password" validate:"omitempty,min=8,max=255" doc:"New password (optional)"`
    Fullname string `json:"fullname" validate:"omitempty,max=255" doc:"Updated full name (optional)"`
    Phone    string `json:"phone" validate:"omitempty,max=20" doc:"Updated phone (optional)"`
    Avatar   string `json:"avatar" validate:"omitempty,url,max=255" doc:"Updated avatar URL (optional)"`
}
```

### 2. Define Request

```go
// File: internal/http/request/user_request.go

type UpdateUser struct {
    dto.UpdateUser
}

// ToDto converts Request to DTO
func (r UpdateUser) ToDto() dto.UpdateUser {
    return r.UpdateUser
}

// SetID populates ID from URL path parameter
// IMPORTANT: Must use pointer receiver to modify the struct!
func (r *UpdateUser) SetID(id int) {
    r.ID = id
}
```

### 3. Implement Controller

```go
// File: internal/http/controllers/api/user/update_user_api.go
package user

import (
    "gfly/internal/constants"
    "gfly/internal/http/request"
    "gfly/internal/http/response"
    "gfly/internal/http/transformers"
    "gfly/internal/services"
    "gfly/pkg/http"
    "github.com/gflydev/core"
)

type UpdateUserApi struct {
    core.Api
}

func NewUpdateUserApi() *UpdateUserApi {
    return &UpdateUserApi{}
}

// Validate handles request parsing, ID extraction, and validation
func (h *UpdateUserApi) Validate(c *core.Ctx) error {
    // ProcessUpdateData: Extract path ID, parse body, sanitize, validate
    return http.ProcessUpdateData[request.UpdateUser](c)
}

// Handle processes the update user request
// @Summary Update user
// @Description Updates an existing user's information
// @Tags Users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param data body request.UpdateUser true "User update data"
// @Success 200 {object} response.User "User updated successfully"
// @Failure 400 {object} response.Error "Validation error"
// @Failure 404 {object} response.NotFound "User not found"
// @Failure 401 {object} response.Unauthorized "Authentication required"
// @Security ApiKeyAuth
// @Router /users/{id} [put]
func (h *UpdateUserApi) Handle(c *core.Ctx) error {
    // 1. Get validated request from context
    requestData := c.GetData(constants.Request).(request.UpdateUser)

    // 2. Convert Request → DTO (DTO.ID is already set by SetID())
    updateUserDto := requestData.ToDto()

    // 3. Call service with DTO
    user, err := services.UpdateUser(updateUserDto)
    if err != nil {
        return c.Error(response.Error{
            Code:    core.StatusBadRequest,
            Message: err.Error(),
        })
    }

    // 4. Transform Model → Response
    userResponse := transformers.ToUserResponse(*user)

    // 5. Return 200 OK with updated user
    return c.JSON(userResponse)
}
```

### 4. Implement Service

```go
// File: internal/services/user_services.go

// UpdateUser updates an existing user
//
// Parameters:
//   - updateUserDto: User update data (includes ID)
//
// Returns:
//   - (*models.User, error): Updated user or error
func UpdateUser(updateUserDto dto.UpdateUser) (*models.User, error) {
    // Retrieve existing user
    user, err := mb.GetModelByID[models.User](updateUserDto.ID)
    if err != nil {
        return nil, errors.New("user not found")
    }

    // Apply updates (only non-empty fields)
    if updateUserDto.Fullname != "" {
        user.Fullname = updateUserDto.Fullname
    }
    if updateUserDto.Phone != "" {
        user.Phone = updateUserDto.Phone
    }
    if updateUserDto.Avatar != "" {
        user.Avatar = dbNull.String(updateUserDto.Avatar)
    }
    if updateUserDto.Password != "" {
        user.Password = utils.GeneratePassword(updateUserDto.Password)
    }

    user.UpdatedAt = time.Now()

    // Persist changes
    err = mb.UpdateModel(user)
    if err != nil {
        log.Errorf("Failed to update user: %v", err)
        return nil, errors.New("failed to update user")
    }

    return user, nil
}
```

### 5. Register Route

```go
// PUT /api/v1/users/:id - Update user
usersGroup.PUT("/:id", app.Apply(
    middleware.JWTAuth(),
    middleware.CheckRolesMiddleware(types.RoleAdmin),
)(user.NewUpdateUserApi()))
```

---

## Delete (DELETE) API

### Complete Flow

```
DELETE /api/v1/users/123

↓ Controller.Validate()
  → PathID(c) extracts ID from URL
  → Store in context

↓ Controller.Handle()
  → Get ID from context
  → Call service: services.DeleteUser(id)
  → Return 200 OK or 204 No Content

Response: 200 OK
{
  "message": "User deleted successfully"
}

OR

Response: 204 No Content
(empty body)
```

### Controller Implementation

```go
// File: internal/http/controllers/api/user/delete_user_api.go
package user

import (
    "gfly/internal/constants"
    "gfly/internal/http/response"
    "gfly/internal/services"
    "gfly/pkg/http"
    "github.com/gflydev/core"
)

type DeleteUserApi struct {
    core.Api
}

func NewDeleteUserApi() *DeleteUserApi {
    return &DeleteUserApi{}
}

// Validate extracts and validates user ID from URL
func (h *DeleteUserApi) Validate(c *core.Ctx) error {
    // ProcessPathID: Extracts ID from path, validates, and stores in context
    return http.ProcessPathID(c)
}

// Handle processes the delete user request
// @Summary Delete user
// @Description Soft deletes a user account
// @Tags Users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} response.Success "User deleted successfully"
// @Success 204 "No content (alternative success response)"
// @Failure 400 {object} response.Error "Invalid ID"
// @Failure 404 {object} response.NotFound "User not found"
// @Failure 401 {object} response.Unauthorized "Authentication required"
// @Security ApiKeyAuth
// @Router /users/{id} [delete]
func (h *DeleteUserApi) Handle(c *core.Ctx) error {
    // 1. Get ID from context
    itemID := c.GetData(constants.ItemID).(int)

    // 2. Call service to delete user
    err := services.DeleteUser(itemID)
    if err != nil {
        return c.Error(response.Error{
            Code:    core.StatusBadRequest,
            Message: err.Error(),
        })
    }

    // 3. Return success response
    // Option 1: 200 OK with message
    return c.JSON(response.Success{
        Message: "User deleted successfully",
    })

    // Option 2: 204 No Content (uncomment to use)
    // return c.Status(core.StatusNoContent).Send(nil)
}
```

### Service Implementation

```go
// File: internal/services/user_services.go

// DeleteUser soft deletes a user
//
// Parameters:
//   - id: User ID
//
// Returns:
//   - error: Error if user not found or deletion failed
func DeleteUser(id int) error {
    // Retrieve user
    user, err := mb.GetModelByID[models.User](id)
    if err != nil {
        return errors.New("user not found")
    }

    // Soft delete (set deleted_at timestamp)
    user.DeletedAt = dbNull.TimeNow()

    err = mb.UpdateModel(user)
    if err != nil {
        log.Errorf("Failed to delete user: %v", err)
        return errors.New("failed to delete user")
    }

    return nil
}

// HardDeleteUser permanently removes a user (use with caution)
func HardDeleteUser(id int) error {
    user, err := mb.GetModelByID[models.User](id)
    if err != nil {
        return errors.New("user not found")
    }

    // Permanent deletion
    err = mb.DeleteModel(user)
    if err != nil {
        log.Errorf("Failed to permanently delete user: %v", err)
        return errors.New("failed to delete user")
    }

    return nil
}
```

### Route Registration

```go
// DELETE /api/v1/users/:id - Delete user
usersGroup.DELETE("/:id", app.Apply(
    middleware.JWTAuth(),
    middleware.CheckRolesMiddleware(types.RoleAdmin),
)(user.NewDeleteUserApi()))
```

---

## List (GET Collection) API

### Complete Flow

```
GET /api/v1/users?keyword=john&page=2&per_page=20&order_by=created_at

↓ Controller.Validate()
  → ProcessFilter(c)
    → Parse query parameters: keyword, page, per_page, order_by
    → Create dto.Filter
    → Sanitize and validate
    → Store in context

↓ Controller.Handle()
  → Get filter from context
  → Call service: services.FindUsers(filterDto)
  → Service returns: ([]models.User, totalCount, error)
  → Transform list: transformers.ToListResponse(users, transformers.ToUserResponse)
  → Create response.ListUser with Meta + Data
  → Return 200 OK

Response: 200 OK
{
  "meta": {
    "page": 2,
    "per_page": 20,
    "total": 154
  },
  "data": [
    {
      "id": 41,
      "email": "john.smith@example.com",
      "fullname": "John Smith",
      ...
    },
    ...
  ]
}
```

### 1. Define Filter DTO

```go
// File: internal/dto/generic_dto.go

// Filter defines query parameters for list/search operations
type Filter struct {
    Keyword string `json:"keyword" validate:"omitempty,max=255" doc:"Search keyword"`
    OrderBy string `json:"order_by" validate:"omitempty" doc:"Sort field"`
    Page    int    `json:"page" validate:"omitempty,gte=1" doc:"Page number (default: 1)"`
    PerPage int    `json:"per_page" validate:"omitempty,gte=1,lte=100" doc:"Items per page (default: 10, max: 100)"`
}
```

### 2. Define List Response

```go
// File: internal/http/response/user_response.go

// ListUser response for paginated user lists
type ListUser struct {
    Meta Meta   `json:"meta" doc:"Pagination metadata"`
    Data []User `json:"data" doc:"List of user objects"`
}

// File: internal/http/response/generic_response.go

// Meta defines pagination metadata
type Meta struct {
    Page    int `json:"page,omitempty" example:"1" doc:"Current page number"`
    PerPage int `json:"per_page,omitempty" example:"10" doc:"Items per page"`
    Total   int `json:"total" example:"1354" doc:"Total number of records"`
}
```

### 3. Implement Controller

```go
// File: internal/http/controllers/api/user/list_users_api.go
package user

import (
    "gfly/internal/constants"
    "gfly/internal/dto"
    "gfly/internal/http/controllers/api"
    "gfly/internal/http/response"
    "gfly/internal/http/transformers"
    "gfly/internal/services"
    "github.com/gflydev/core"
)

type ListUsersApi struct {
    api.ListApi  // Base controller with Validate() that calls ProcessFilter()
}

func NewListUsersApi() *ListUsersApi {
    return &ListUsersApi{}
}

// Validate is inherited from api.ListApi, which calls ProcessFilter()

// Handle retrieves a paginated list of users
// @Summary List users
// @Description Retrieves a paginated list of users with optional filtering and sorting
// @Description <b>Keyword fields:</b> email, fullname, phone, status
// @Description <b>Order by fields:</b> email, fullname, created_at, last_access_at
// @Tags Users
// @Accept json
// @Produce json
// @Param keyword query string false "Search keyword"
// @Param order_by query string false "Sort field (e.g., 'email', 'created_at')"
// @Param page query int false "Page number (default: 1)"
// @Param per_page query int false "Items per page (default: 10, max: 100)"
// @Success 200 {object} response.ListUser "List of users"
// @Failure 400 {object} response.Error "Invalid query parameters"
// @Failure 401 {object} response.Unauthorized "Authentication required"
// @Security ApiKeyAuth
// @Router /users [get]
func (h *ListUsersApi) Handle(c *core.Ctx) error {
    // 1. Get validated filter from context
    filterDto := c.GetData(constants.Filter).(dto.Filter)

    // 2. Call service with filter DTO
    users, total, err := services.FindUsers(filterDto)
    if err != nil {
        return c.Error(response.Error{
            Code:    core.StatusInternalServerError,
            Message: "Failed to retrieve users",
        })
    }

    // 3. Create pagination metadata
    metadata := response.Meta{
        Page:    filterDto.Page,
        PerPage: filterDto.PerPage,
        Total:   total,
    }

    // 4. Transform list of models → list of responses
    data := transformers.ToListResponse(users, transformers.ToUserResponse)

    // 5. Return paginated response
    return c.JSON(response.ListUser{
        Meta: metadata,
        Data: data,
    })
}
```

### 4. Base ListApi Controller

```go
// File: internal/http/controllers/api/generic_api.go
package api

import (
    "gfly/pkg/http"
    "github.com/gflydev/core"
)

// ListApi base controller for list operations
type ListApi struct {
    core.Api
}

// Validate handles filter validation (inherited by all list controllers)
func (h *ListApi) Validate(c *core.Ctx) error {
    return http.ProcessFilter(c)
}
```

### 5. Implement Service

```go
// File: internal/services/user_services.go

// FindUsers retrieves a filtered and paginated list of users
//
// Parameters:
//   - filterDto: Filter criteria (keyword, ordering, pagination)
//
// Returns:
//   - ([]models.User, int, error): Users, total count, and error
func FindUsers(filterDto dto.Filter) ([]models.User, int, error) {
    dbInstance := mb.Instance()
    var users []models.User
    var total int
    var offset = 0

    // Calculate offset for pagination
    if filterDto.Page > 0 {
        offset = (filterDto.Page - 1) * filterDto.PerPage
    }

    // Build query
    builder := dbInstance.Select("DISTINCT users.id", "users.*").
        From(models.TableUser)

    // Apply keyword search (if provided)
    if filterDto.Keyword != "" {
        builder.Where(mb.Group(func(q mb.IQueryBuilder) {
            q.Where(mb.Condition{
                Field: "users.email",
                Opt:   mb.Like,
                Value: "%" + filterDto.Keyword + "%",
            }).OrWhere(mb.Condition{
                Field: "users.fullname",
                Opt:   mb.Like,
                Value: "%" + filterDto.Keyword + "%",
            }).OrWhere(mb.Condition{
                Field: "users.phone",
                Opt:   mb.Like,
                Value: "%" + filterDto.Keyword + "%",
            })
        }))
    }

    // Apply ordering (if provided)
    if filterDto.OrderBy != "" {
        builder.Order(filterDto.OrderBy, mb.Desc)
    } else {
        builder.Order("users.created_at", mb.Desc)  // Default ordering
    }

    // Apply pagination
    builder.Limit(filterDto.PerPage, offset)

    // Execute query
    total, err := builder.Find(&users)
    if err != nil {
        log.Errorf("Failed to retrieve users: %v", err)
        return nil, 0, errors.New("failed to retrieve users")
    }

    return users, total, nil
}
```

### 6. Register Route

```go
// GET /api/v1/users - List users
usersGroup.GET("", app.Apply(
    middleware.JWTAuth(),
)(user.NewListUsersApi()))
```

---

## Advanced: Custom Filters with Multiple Parameters

### When to Use Custom Filters

The basic `http.ProcessFilter(c)` approach works for simple list operations with keyword search, pagination, and ordering. However, when you need **resource-specific filters** (e.g., filtering products by price, brand, or status), you need to create a **custom filter DTO** and extraction function.

**Use custom filters when:**
- You need to filter by specific fields (price, status, category, date ranges, etc.)
- You have complex boolean filters (is_featured, is_active, etc.)
- You need to validate filter values
- You want type-safe filter parameters

### Pattern Overview

```
GET /api/v1/products?keyword=phone&is_featured=true&brand_id=5&price=99.99&page=1&per_page=20

↓ Controller.Validate()
  → Custom FilterProductDto(c) function
    → Get base filter (keyword, page, per_page, order_by) using http.FilterData(c)
    → Extract product-specific query params (is_featured, price, brand_id, etc.)
    → Create FilterProduct DTO with all parameters
    → Validate the custom filter DTO
    → Store in context

↓ Controller.Handle()
  → Get custom filter from context
  → Call service with custom filter DTO
  → Return paginated response
```

### 1. Define Custom Filter DTO

Create a filter DTO that **embeds the base `dto.Filter`** and adds resource-specific fields:

```go
// File: internal/dto/product_dto.go
package dto

// FilterProduct extends base Filter with product-specific filters
type FilterProduct struct {
    Filter                              // Embed base filter (keyword, page, per_page, order_by)
    IsFeatured *bool    `json:"is_featured" example:"true" validate:"omitempty" doc:"Filter by featured status"`
    Price      *float64 `json:"price" example:"99.99" validate:"omitempty,gte=0" doc:"Filter by price"`
    Currency   *string  `json:"currency" example:"USD" validate:"omitempty" doc:"Filter by currency"`
    VarietyID  *int     `json:"variety_id" example:"1" validate:"omitempty,min=1" doc:"Filter by variety ID"`
    BrandID    *int     `json:"brand_id" example:"1" validate:"omitempty,min=1" doc:"Filter by brand ID"`
    IsActive   *bool    `json:"is_active" example:"true" validate:"omitempty" doc:"Filter by active status"`
}
```

**Key points:**
- Use **pointer types** (`*bool`, `*int`, `*string`) to distinguish between "not provided" (nil) and "zero value"
- Embed `dto.Filter` to inherit base filter fields
- Add validation tags for each filter field
- Use `omitempty` to make all filters optional

### 2. Create Custom Filter Extraction Function

Create a helper function to extract query parameters and build the custom filter DTO:

```go
// File: internal/http/request/product_request.go
package request

import (
    "gfly/internal/dto"
    "gfly/pkg/http"
    "github.com/gflydev/core"
)

// FilterProductDto extracts and builds a FilterProduct from query parameters
func FilterProductDto(c *core.Ctx) dto.FilterProduct {
    // 1. Get base filter data (keyword, page, per_page, order_by)
    baseFilter := http.FilterData(c)

    // 2. Create product filter DTO with base filter
    filterDto := dto.FilterProduct{
        Filter: baseFilter,
    }

    // 3. Extract product-specific filters from query parameters
    if isFeatured, err := c.QueryBool("is_featured"); err == nil && isFeatured != false {
        filterDto.IsFeatured = &isFeatured
    }

    if price, err := c.QueryFloat("price"); err == nil && price >= 0 {
        filterDto.Price = &price
    }

    if currency := c.QueryStr("currency"); currency != "" {
        filterDto.Currency = &currency
    }

    if varietyID, err := c.QueryInt("variety_id"); err == nil && varietyID > 0 {
        filterDto.VarietyID = &varietyID
    }

    if brandID, err := c.QueryInt("brand_id"); err == nil && brandID > 0 {
        filterDto.BrandID = &brandID
    }

    if isActive, err := c.QueryBool("is_active"); err == nil && isActive != false {
        filterDto.IsActive = &isActive
    }

    // 4. Set default ordering if not provided
    if filterDto.OrderBy == "" {
        filterDto.OrderBy = "-created_at"  // Default: newest first
    }

    return filterDto
}
```

**Available context methods for extracting query parameters:**
- `c.QueryBool(key)` - Extract boolean (`?is_active=true`)
- `c.QueryInt(key)` - Extract integer (`?brand_id=5`)
- `c.QueryFloat(key)` - Extract float64 (`?price=99.99`)
- `c.QueryStr(key)` - Extract string (`?currency=USD`)

### 3. Update Controller Validation

Use the custom filter extraction function in the controller's `Validate()` method:

```go
// File: internal/http/controllers/api/product/list_products_api.go
package product

import (
    "gfly/internal/constants"
    "gfly/internal/dto"
    "gfly/internal/http/controllers/api"
    "gfly/internal/http/request"
    "gfly/internal/http/response"
    "gfly/internal/http/transformers"
    "gfly/internal/services"
    "github.com/gflydev/core"
)

type ListProductsApi struct {
    core.Api  // NOTE: Do NOT embed api.ListApi when using custom filters
}

func NewListProductsApi() *ListProductsApi {
    return &ListProductsApi{}
}

// Validate extracts and validates custom product filters
func (h *ListProductsApi) Validate(c *core.Ctx) error {
    // 1. Extract custom filter using helper function
    filterDto := request.FilterProductDto(c)

    // 2. Validate the custom filter DTO
    if errData := api.Validate(filterDto); errData != nil {
        return c.Error(errData)
    }

    // 3. Store validated filter in context
    c.SetData(constants.Filter, filterDto)

    return nil
}

// Handle retrieves a paginated list of products with custom filters
// @Summary List products with filters
// @Description Retrieves a paginated list of products with optional filtering and sorting
// @Description <b>Keyword fields:</b> name, description
// @Description <b>Order by fields:</b> name, price, created_at
// @Description <b>Custom filters:</b> is_featured, price, currency, variety_id, brand_id, is_active
// @Tags Products
// @Accept json
// @Produce json
// @Param keyword query string false "Search keyword"
// @Param order_by query string false "Sort field (e.g., 'price', '-created_at')"
// @Param page query int false "Page number (default: 1)"
// @Param per_page query int false "Items per page (default: 10, max: 100)"
// @Param is_featured query bool false "Filter by featured status"
// @Param price query number false "Filter by price"
// @Param currency query string false "Filter by currency"
// @Param variety_id query int false "Filter by variety ID"
// @Param brand_id query int false "Filter by brand ID"
// @Param is_active query bool false "Filter by active status"
// @Success 200 {object} response.ListProduct "List of products"
// @Failure 400 {object} response.Error "Invalid query parameters"
// @Failure 401 {object} response.Unauthorized "Authentication required"
// @Security ApiKeyAuth
// @Router /products [get]
func (h *ListProductsApi) Handle(c *core.Ctx) error {
    // 1. Get validated custom filter from context
    filterDto := c.GetData(constants.Filter).(dto.FilterProduct)

    // 2. Call service with custom filter DTO
    products, total, err := services.ListProducts(filterDto)
    if err != nil {
        return c.Error(response.Error{
            Code:    core.StatusBadRequest,
            Message: err.Error(),
        })
    }

    // 3. Create pagination metadata
    metadata := response.Meta{
        Page:    filterDto.Page,
        PerPage: filterDto.PerPage,
        Total:   total,
    }

    // 4. Transform list of models → list of responses
    data := transformers.ToListResponse(products, transformers.ToProductResponse)

    // 5. Return paginated response
    return c.JSON(response.ListProduct{
        Meta: metadata,
        Data: data,
    })
}
```

**Important notes:**
- Do NOT embed `api.ListApi` when using custom filters (it uses basic `ProcessFilter`)
- Use `api.Validate()` to validate the custom filter DTO
- Document all custom filter parameters in Swagger annotations

### 4. Update Service to Handle Custom Filters

Update the service layer to accept and process the custom filter DTO:

```go
// File: internal/services/product_services.go
package services

import (
    "gfly/internal/domain/models"
    "gfly/internal/dto"
    mb "github.com/gflydev/db"
    "github.com/gflydev/core/errors"
)

// ListProducts retrieves a filtered and paginated list of products
//
// Parameters:
//   - filterDto: Custom filter with product-specific criteria
//
// Returns:
//   - ([]models.Product, int, error): Products, total count, and error
func ListProducts(filterDto dto.FilterProduct) ([]models.Product, int, error) {
    dbInstance := mb.Instance()
    var products []models.Product
    var total int
    var offset = 0

    // Calculate offset for pagination
    if filterDto.Page > 0 {
        offset = (filterDto.Page - 1) * filterDto.PerPage
    }

    // Build query
    builder := dbInstance.Select("*").
        From(models.TableProduct).
        Where(mb.Condition{
            Field: "deleted_at",
            Opt:   mb.IsNull,
        })

    // Apply keyword search (base filter)
    if filterDto.Keyword != "" {
        builder.Where(mb.Group(func(q mb.IQueryBuilder) {
            q.Where(mb.Condition{
                Field: "name",
                Opt:   mb.Like,
                Value: "%" + filterDto.Keyword + "%",
            }).OrWhere(mb.Condition{
                Field: "description",
                Opt:   mb.Like,
                Value: "%" + filterDto.Keyword + "%",
            })
        }))
    }

    // Apply custom filters (product-specific)
    if filterDto.IsFeatured != nil {
        builder.Where(mb.Condition{
            Field: "is_featured",
            Opt:   mb.Eq,
            Value: *filterDto.IsFeatured,
        })
    }

    if filterDto.Price != nil {
        builder.Where(mb.Condition{
            Field: "price",
            Opt:   mb.Eq,
            Value: *filterDto.Price,
        })
    }

    if filterDto.Currency != nil {
        builder.Where(mb.Condition{
            Field: "currency",
            Opt:   mb.Eq,
            Value: *filterDto.Currency,
        })
    }

    if filterDto.VarietyID != nil {
        builder.Where(mb.Condition{
            Field: "variety_id",
            Opt:   mb.Eq,
            Value: *filterDto.VarietyID,
        })
    }

    if filterDto.BrandID != nil {
        builder.Where(mb.Condition{
            Field: "brand_id",
            Opt:   mb.Eq,
            Value: *filterDto.BrandID,
        })
    }

    if filterDto.IsActive != nil {
        builder.Where(mb.Condition{
            Field: "is_active",
            Opt:   mb.Eq,
            Value: *filterDto.IsActive,
        })
    }

    // Apply ordering (base filter)
    if filterDto.OrderBy != "" {
        builder.Order(filterDto.OrderBy, mb.Desc)
    } else {
        builder.Order("created_at", mb.Desc)
    }

    // Apply pagination (base filter)
    builder.Limit(filterDto.PerPage, offset)

    // Execute query
    total, err := builder.Find(&products)
    if err != nil {
        log.Errorf("Failed to retrieve products: %v", err)
        return nil, 0, errors.New("failed to retrieve products")
    }

    return products, total, nil
}
```

**Key patterns:**
- Check if pointer fields are `nil` before applying filters
- Dereference pointer values with `*filterDto.FieldName`
- Combine base filter logic (keyword, pagination) with custom filter logic

### 5. Complete Example: User List with Role Filter

Here's another complete example for filtering users by role:

```go
// ==================== User Filter DTO ====================
// File: internal/dto/user_dto.go

type FilterUser struct {
    Filter
    Status   *string `json:"status" example:"active" validate:"omitempty,oneof=active inactive suspended" doc:"Filter by user status"`
    Role     *string `json:"role" example:"admin" validate:"omitempty,oneof=admin user guest" doc:"Filter by user role"`
    Verified *bool   `json:"verified" example:"true" validate:"omitempty" doc:"Filter by email verification status"`
}

// ==================== User Filter Request ====================
// File: internal/http/request/user_request.go

func FilterUserDto(c *core.Ctx) dto.FilterUser {
    baseFilter := http.FilterData(c)

    filterDto := dto.FilterUser{
        Filter: baseFilter,
    }

    if status := c.QueryStr("status"); status != "" {
        filterDto.Status = &status
    }

    if role := c.QueryStr("role"); role != "" {
        filterDto.Role = &role
    }

    if verified, err := c.QueryBool("verified"); err == nil {
        filterDto.Verified = &verified
    }

    if filterDto.OrderBy == "" {
        filterDto.OrderBy = "-created_at"
    }

    return filterDto
}

// ==================== Controller ====================
// File: internal/http/controllers/api/user/list_users_api.go

func (h *ListUsersApi) Validate(c *core.Ctx) error {
    filterDto := request.FilterUserDto(c)

    if errData := api.Validate(filterDto); errData != nil {
        return c.Error(errData)
    }

    c.SetData(constants.Filter, filterDto)
    return nil
}

func (h *ListUsersApi) Handle(c *core.Ctx) error {
    filterDto := c.GetData(constants.Filter).(dto.FilterUser)

    users, total, err := services.ListUsers(filterDto)
    if err != nil {
        return c.Error(response.Error{
            Code:    core.StatusInternalServerError,
            Message: "Failed to retrieve users",
        })
    }

    metadata := response.Meta{
        Page:    filterDto.Page,
        PerPage: filterDto.PerPage,
        Total:   total,
    }

    data := transformers.ToListResponse(users, transformers.ToUserResponse)

    return c.JSON(response.ListUser{
        Meta: metadata,
        Data: data,
    })
}
```

### Best Practices for Custom Filters

✅ **DO:**
- Use pointer types for optional filter fields
- Embed `dto.Filter` for base filter functionality
- Validate custom filters using `api.Validate()`
- Check for `nil` before applying filters in the service layer
- Document all filter parameters in Swagger annotations
- Provide sensible default values (e.g., default ordering)
- Use validation tags to ensure filter values are valid

❌ **DON'T:**
- Use non-pointer types (can't distinguish between "not set" and "zero value")
- Forget to validate custom filters
- Apply filters without nil checks (will cause panics)
- Mix basic `ProcessFilter()` with custom filters (choose one approach)
- Embed `api.ListApi` when using custom filters
- Expose internal field names in query parameters

### When to Use Basic vs. Custom Filters

| Use Case | Approach | Example |
|----------|----------|---------|
| Simple list with keyword search | `http.ProcessFilter()` + `api.ListApi` | `GET /users?keyword=john&page=1` |
| Resource-specific filters | Custom filter DTO + extraction function | `GET /products?brand_id=5&is_featured=true` |
| Date range filters | Custom filter DTO | `GET /orders?start_date=2025-01-01&end_date=2025-01-31` |
| Enum filters | Custom filter DTO with validation | `GET /users?status=active&role=admin` |

---

## Response Patterns

### Success Response Patterns

#### 1. Single Resource Response

```go
// Pattern: Return the resource object directly
return c.JSON(response.User{
    ID:        123,
    Email:     "john@example.com",
    Fullname:  "John Doe",
})
```

**Output:**
```json
{
  "id": 123,
  "email": "john@example.com",
  "fullname": "John Doe"
}
```

#### 2. List/Collection Response

```go
// Pattern: Wrap in Meta + Data structure
return c.JSON(response.ListUser{
    Meta: response.Meta{
        Page:    1,
        PerPage: 10,
        Total:   154,
    },
    Data: []response.User{...},
})
```

**Output:**
```json
{
  "meta": {
    "page": 1,
    "per_page": 10,
    "total": 154
  },
  "data": [...]
}
```

#### 3. Generic Success Response

```go
// Pattern: Message + optional data
return c.JSON(response.Success{
    Message: "Operation completed successfully",
    Data:    core.Data{"id": 123},
})
```

**Output:**
```json
{
  "message": "Operation completed successfully",
  "data": {
    "id": 123
  }
}
```

#### 4. Empty Success Response

```go
// Pattern: 204 No Content (no response body)
return c.Status(core.StatusNoContent).Send(nil)
```

**Output:** (empty, status 204)

### Transformer Patterns

#### Single Resource Transformation

```go
// Transform single model → response
userResponse := transformers.ToUserResponse(user)
return c.JSON(userResponse)
```

#### List Transformation

```go
// Transform list of models → list of responses
data := transformers.ToListResponse(users, transformers.ToUserResponse)
return c.JSON(response.ListUser{
    Meta: metadata,
    Data: data,
})
```

#### Generic List Transformer

```go
// File: internal/http/transformers/generic_transformer.go

// ToListResponse transforms a list using the provided transformer function
func ToListResponse[T any, R any](records []T, transformerFn func(T) R) []R {
    return fn.TransformList(records, transformerFn)
}
```

**Usage:**
```go
// Users
userResponses := transformers.ToListResponse(users, transformers.ToUserResponse)

// Products
productResponses := transformers.ToListResponse(products, transformers.ToProductResponse)

// Any model
responses := transformers.ToListResponse(models, transformers.ToResponse)
```

---

## HTTP Status Codes

### Success Status Codes

| Code | Name | Usage | Example |
|------|------|-------|---------|
| **200** | OK | Successful GET, PUT, PATCH, DELETE | `GET /users/123` |
| **201** | Created | Successful POST (resource created) | `POST /users` |
| **204** | No Content | Successful DELETE (no response body) | `DELETE /users/123` |

### Client Error Status Codes

| Code | Name | Usage | Example |
|------|------|-------|---------|
| **400** | Bad Request | Validation error, malformed request | Invalid email format |
| **401** | Unauthorized | Authentication required/failed | Missing or invalid JWT |
| **403** | Forbidden | Authenticated but no permission | Non-admin accessing admin endpoint |
| **404** | Not Found | Resource doesn't exist | User ID 999 not found |
| **409** | Conflict | Resource conflict | Email already exists |
| **422** | Unprocessable Entity | Semantic validation error | Business rule violation |
| **429** | Too Many Requests | Rate limit exceeded | Too many login attempts |

### Server Error Status Codes

| Code | Name | Usage | Example |
|------|------|-------|---------|
| **500** | Internal Server Error | Unexpected server error | Database connection failed |
| **503** | Service Unavailable | Service temporarily unavailable | Maintenance mode |

### Status Code Usage in Controllers

#### Success Responses

```go
// 200 OK (default)
return c.JSON(response)

// 201 Created
return c.Status(core.StatusCreated).JSON(response)

// 204 No Content
return c.Status(core.StatusNoContent).Send(nil)
```

#### Error Responses

```go
// 400 Bad Request (validation error)
return c.Error(response.Error{
    Code:    core.StatusBadRequest,
    Message: "Invalid input",
    Data:    validationErrors,
})

// 401 Unauthorized
return c.Error(response.Unauthorized{
    Code:    core.StatusUnauthorized,
    Message: "Authentication required",
})

// 404 Not Found
return c.Error(response.NotFound{
    Code:    core.StatusNotFound,
    Message: "User not found",
})

// 409 Conflict
return c.Error(response.Conflict{
    Code:    core.StatusConflict,
    Message: "Email already exists",
})

// 500 Internal Server Error
return c.Error(response.Error{
    Code:    core.StatusInternalServerError,
    Message: "An unexpected error occurred",
})
```

---

## Error Handling

### Error Response Structure

```go
// File: internal/http/response/generic_response.go

// Error represents a generic error response
type Error struct {
    Code    int       `json:"code" example:"400"`
    Message string    `json:"message" example:"Bad request"`
    Data    core.Data `json:"data" doc:"Field-level validation errors"`
}

// Unauthorized represents 401 error
type Unauthorized struct {
    Code    int    `json:"code" example:"401"`
    Message string `json:"error" example:"Unauthorized access"`
}

// NotFound represents 404 error
type NotFound struct {
    Code    int    `json:"code" example:"404"`
    Message string `json:"error" example:"Resource not found"`
}

// Conflict represents 409 error
type Conflict struct {
    Code    int    `json:"code" example:"409"`
    Message string `json:"error" example:"Resource conflict"`
}
```

### Validation Error Response

When validation fails, the response includes field-level errors:

```json
{
  "code": 400,
  "message": "Invalid input",
  "data": {
    "email": ["email must be a valid email address"],
    "password": ["password must be at least 8 characters"],
    "phone": ["phone must be at most 20 characters"]
  }
}
```

### Error Handling in Controllers

```go
// Validation error (from ProcessData/ProcessUpdateData)
func (h *CreateApi) Validate(c *core.Ctx) error {
    return http.ProcessData[request.CreateUser](c)
    // Automatically returns 400 with field errors if validation fails
}

// Business logic error (from service)
func (h *CreateApi) Handle(c *core.Ctx) error {
    user, err := services.CreateUser(dto)
    if err != nil {
        // Check error type and return appropriate status
        if strings.Contains(err.Error(), "already exists") {
            return c.Error(response.Conflict{
                Code:    core.StatusConflict,
                Message: err.Error(),
            })
        }

        return c.Error(response.Error{
            Code:    core.StatusBadRequest,
            Message: err.Error(),
        })
    }

    return c.Status(core.StatusCreated).JSON(transformers.ToUserResponse(*user))
}
```

### Error Handling in Services

```go
// File: internal/services/user_services.go

import "github.com/gflydev/core/errors"

func CreateUser(dto dto.CreateUser) (*models.User, error) {
    // Business rule validation
    if existingUser := repository.Pool.GetUserByEmail(dto.Email); existingUser != nil {
        return nil, errors.New("user with this email already exists")
    }

    // Database operation error
    err := mb.CreateModel(user)
    if err != nil {
        log.Errorf("Database error: %v", err)
        return nil, errors.New("failed to create user")  // Generic error for client
    }

    return user, nil
}
```

**Important:** Use `github.com/gflydev/core/errors` package, NOT `fmt.Errorf`

---

## Complete CRUD Example

Let's create a complete CRUD API for a "Product" resource.

### 1. Domain Model

```go
// File: internal/domain/models/product_model.go
package models

import (
    "database/sql"
    mb "github.com/gflydev/db"
    "time"
)

const TableProduct = "products"

type Product struct {
    MetaData    mb.MetaData    `db:"-" model:"table:products"`
    ID          int            `db:"id" model:"name:id; type:serial,primary"`
    Name        string         `db:"name" model:"name:name"`
    Description sql.NullString `db:"description" model:"name:description"`
    Price       float64        `db:"price" model:"name:price"`
    Stock       int            `db:"stock" model:"name:stock"`
    ImageURL    sql.NullString `db:"image_url" model:"name:image_url"`
    Status      string         `db:"status" model:"name:status"`
    CreatedAt   time.Time      `db:"created_at" model:"name:created_at"`
    UpdatedAt   time.Time      `db:"updated_at" model:"name:updated_at"`
    DeletedAt   sql.NullTime   `db:"deleted_at" model:"name:deleted_at"`
}
```

### 2. DTOs

```go
// File: internal/dto/product_dto.go
package dto

type CreateProduct struct {
    Name        string  `json:"name" validate:"required,max=255" doc:"Product name"`
    Description string  `json:"description" validate:"omitempty,max=1000" doc:"Product description"`
    Price       float64 `json:"price" validate:"required,gte=0" doc:"Product price (must be >= 0)"`
    Stock       int     `json:"stock" validate:"required,gte=0" doc:"Stock quantity (must be >= 0)"`
    ImageURL    string  `json:"image_url" validate:"omitempty,url,max=500" doc:"Product image URL"`
    Status      string  `json:"status" validate:"omitempty,oneof=active inactive" doc:"Product status"`
}

type UpdateProduct struct {
    ID          int     `json:"-" validate:"omitempty,gte=1" doc:"Product ID"`
    Name        string  `json:"name" validate:"omitempty,max=255" doc:"Product name"`
    Description string  `json:"description" validate:"omitempty,max=1000" doc:"Product description"`
    Price       float64 `json:"price" validate:"omitempty,gte=0" doc:"Product price"`
    Stock       int     `json:"stock" validate:"omitempty,gte=0" doc:"Stock quantity"`
    ImageURL    string  `json:"image_url" validate:"omitempty,url,max=500" doc:"Product image URL"`
    Status      string  `json:"status" validate:"omitempty,oneof=active inactive" doc:"Product status"`
}
```

### 3. Requests

```go
// File: internal/http/request/product_request.go
package request

import "gfly/internal/dto"

type CreateProduct struct {
    dto.CreateProduct
}

func (r CreateProduct) ToDto() dto.CreateProduct {
    return r.CreateProduct
}

type UpdateProduct struct {
    dto.UpdateProduct
}

func (r UpdateProduct) ToDto() dto.UpdateProduct {
    return r.UpdateProduct
}

func (r *UpdateProduct) SetID(id int) {
    r.ID = id
}
```

### 4. Response

```go
// File: internal/http/response/product_response.go
package response

import "time"

type Product struct {
    ID          int        `json:"id" doc:"Product ID"`
    Name        string     `json:"name" doc:"Product name"`
    Description *string    `json:"description" doc:"Product description"`
    Price       float64    `json:"price" doc:"Product price"`
    Stock       int        `json:"stock" doc:"Stock quantity"`
    ImageURL    *string    `json:"image_url" doc:"Product image URL"`
    Status      string     `json:"status" doc:"Product status"`
    CreatedAt   time.Time  `json:"created_at" doc:"Creation timestamp"`
    UpdatedAt   time.Time  `json:"updated_at" doc:"Last update timestamp"`
}

type ListProduct struct {
    Meta Meta      `json:"meta" doc:"Pagination metadata"`
    Data []Product `json:"data" doc:"List of products"`
}
```

### 5. Transformer

```go
// File: internal/http/transformers/product_transformer.go
package transformers

import (
    "gfly/internal/domain/models"
    "gfly/internal/http/response"
    dbNull "github.com/gflydev/db/null"
)

func ToProductResponse(product models.Product) response.Product {
    return response.Product{
        ID:          product.ID,
        Name:        product.Name,
        Description: dbNull.StringNil(product.Description),
        Price:       product.Price,
        Stock:       product.Stock,
        ImageURL:    dbNull.StringNil(product.ImageURL),
        Status:      product.Status,
        CreatedAt:   product.CreatedAt,
        UpdatedAt:   product.UpdatedAt,
    }
}
```

### 6. Service

```go
// File: internal/services/product_services.go
package services

import (
    "gfly/internal/domain/models"
    "gfly/internal/dto"
    mb "github.com/gflydev/db"
    dbNull "github.com/gflydev/db/null"
    "github.com/gflydev/core/errors"
    "time"
)

func CreateProduct(dto dto.CreateProduct) (*models.Product, error) {
    product := &models.Product{
        Name:        dto.Name,
        Description: dbNull.String(dto.Description),
        Price:       dto.Price,
        Stock:       dto.Stock,
        ImageURL:    dbNull.String(dto.ImageURL),
        Status:      "active",
        CreatedAt:   time.Now(),
        UpdatedAt:   time.Now(),
    }

    if dto.Status != "" {
        product.Status = dto.Status
    }

    err := mb.CreateModel(product)
    if err != nil {
        return nil, errors.New("failed to create product")
    }

    return product, nil
}

func GetProductByID(id int) (*models.Product, error) {
    product, err := mb.GetModelByID[models.Product](id)
    if err != nil {
        return nil, errors.New("product not found")
    }
    return product, nil
}

func UpdateProduct(dto dto.UpdateProduct) (*models.Product, error) {
    product, err := mb.GetModelByID[models.Product](dto.ID)
    if err != nil {
        return nil, errors.New("product not found")
    }

    if dto.Name != "" {
        product.Name = dto.Name
    }
    if dto.Description != "" {
        product.Description = dbNull.String(dto.Description)
    }
    if dto.Price > 0 {
        product.Price = dto.Price
    }
    if dto.Stock >= 0 {
        product.Stock = dto.Stock
    }
    if dto.ImageURL != "" {
        product.ImageURL = dbNull.String(dto.ImageURL)
    }
    if dto.Status != "" {
        product.Status = dto.Status
    }

    product.UpdatedAt = time.Now()

    err = mb.UpdateModel(product)
    if err != nil {
        return nil, errors.New("failed to update product")
    }

    return product, nil
}

func DeleteProduct(id int) error {
    product, err := mb.GetModelByID[models.Product](id)
    if err != nil {
        return errors.New("product not found")
    }

    product.DeletedAt = dbNull.TimeNow()
    err = mb.UpdateModel(product)
    if err != nil {
        return errors.New("failed to delete product")
    }

    return nil
}

func FindProducts(filterDto dto.Filter) ([]models.Product, int, error) {
    dbInstance := mb.Instance()
    var products []models.Product
    var total int
    var offset = 0

    if filterDto.Page > 0 {
        offset = (filterDto.Page - 1) * filterDto.PerPage
    }

    builder := dbInstance.Select("*").
        From(models.TableProduct).
        Where(mb.Condition{
            Field: "deleted_at",
            Opt:   mb.IsNull,
        })

    if filterDto.Keyword != "" {
        builder.Where(mb.Condition{
            Field: "name",
            Opt:   mb.Like,
            Value: "%" + filterDto.Keyword + "%",
        })
    }

    if filterDto.OrderBy != "" {
        builder.Order(filterDto.OrderBy, mb.Desc)
    } else {
        builder.Order("created_at", mb.Desc)
    }

    builder.Limit(filterDto.PerPage, offset)

    total, err := builder.Find(&products)
    if err != nil {
        return nil, 0, errors.New("failed to retrieve products")
    }

    return products, total, nil
}
```

### 7. Controllers

```go
// File: internal/http/controllers/api/product/create_product_api.go
package product

import (
    "gfly/internal/constants"
    "gfly/internal/http/request"
    "gfly/internal/http/response"
    "gfly/internal/http/transformers"
    "gfly/internal/services"
    "gfly/pkg/http"
    "github.com/gflydev/core"
)

type CreateProductApi struct {
    core.Api
}

func NewCreateProductApi() *CreateProductApi {
    return &CreateProductApi{}
}

func (h *CreateProductApi) Validate(c *core.Ctx) error {
    return http.ProcessData[request.CreateProduct](c)
}

// @Summary Create product
// @Tags Products
// @Accept json
// @Produce json
// @Param data body request.CreateProduct true "Product data"
// @Success 201 {object} response.Product
// @Failure 400 {object} response.Error
// @Security ApiKeyAuth
// @Router /products [post]
func (h *CreateProductApi) Handle(c *core.Ctx) error {
    requestData := c.GetData(constants.Request).(request.CreateProduct)
    product, err := services.CreateProduct(requestData.ToDto())
    if err != nil {
        return c.Error(response.Error{
            Code:    core.StatusBadRequest,
            Message: err.Error(),
        })
    }

    return c.Status(core.StatusCreated).JSON(transformers.ToProductResponse(*product))
}

// File: internal/http/controllers/api/product/get_product_api.go
type GetProductApi struct {
    core.Api
}

func NewGetProductApi() *GetProductApi {
    return &GetProductApi{}
}

func (h *GetProductApi) Validate(c *core.Ctx) error {
    return http.ProcessPathID(c)
}

// @Summary Get product by ID
// @Tags Products
// @Param id path int true "Product ID"
// @Success 200 {object} response.Product
// @Failure 404 {object} response.NotFound
// @Security ApiKeyAuth
// @Router /products/{id} [get]
func (h *GetProductApi) Handle(c *core.Ctx) error {
    itemID := c.GetData(constants.ItemID).(int)
    product, err := services.GetProductByID(itemID)
    if err != nil {
        return c.Error(response.NotFound{
            Code:    core.StatusNotFound,
            Message: "Product not found",
        })
    }

    return c.JSON(transformers.ToProductResponse(*product))
}

// File: internal/http/controllers/api/product/update_product_api.go
type UpdateProductApi struct {
    core.Api
}

func NewUpdateProductApi() *UpdateProductApi {
    return &UpdateProductApi{}
}

func (h *UpdateProductApi) Validate(c *core.Ctx) error {
    return http.ProcessUpdateData[request.UpdateProduct](c)
}

// @Summary Update product
// @Tags Products
// @Param id path int true "Product ID"
// @Param data body request.UpdateProduct true "Product data"
// @Success 200 {object} response.Product
// @Failure 400 {object} response.Error
// @Security ApiKeyAuth
// @Router /products/{id} [put]
func (h *UpdateProductApi) Handle(c *core.Ctx) error {
    requestData := c.GetData(constants.Request).(request.UpdateProduct)
    product, err := services.UpdateProduct(requestData.ToDto())
    if err != nil {
        return c.Error(response.Error{
            Code:    core.StatusBadRequest,
            Message: err.Error(),
        })
    }

    return c.JSON(transformers.ToProductResponse(*product))
}

// File: internal/http/controllers/api/product/delete_product_api.go
type DeleteProductApi struct {
    core.Api
}

func NewDeleteProductApi() *DeleteProductApi {
    return &DeleteProductApi{}
}

func (h *DeleteProductApi) Validate(c *core.Ctx) error {
    return http.ProcessPathID(c)
}

// @Summary Delete product
// @Tags Products
// @Param id path int true "Product ID"
// @Success 200 {object} response.Success
// @Failure 404 {object} response.NotFound
// @Security ApiKeyAuth
// @Router /products/{id} [delete]
func (h *DeleteProductApi) Handle(c *core.Ctx) error {
    itemID := c.GetData(constants.ItemID).(int)
    err := services.DeleteProduct(itemID)
    if err != nil {
        return c.Error(response.Error{
            Code:    core.StatusBadRequest,
            Message: err.Error(),
        })
    }

    return c.JSON(response.Success{
        Message: "Product deleted successfully",
    })
}

// File: internal/http/controllers/api/product/list_products_api.go
type ListProductsApi struct {
    api.ListApi
}

func NewListProductsApi() *ListProductsApi {
    return &ListProductsApi{}
}

// @Summary List products
// @Tags Products
// @Param keyword query string false "Search keyword"
// @Param page query int false "Page number"
// @Param per_page query int false "Items per page"
// @Success 200 {object} response.ListProduct
// @Security ApiKeyAuth
// @Router /products [get]
func (h *ListProductsApi) Handle(c *core.Ctx) error {
    filterDto := c.GetData(constants.Filter).(dto.Filter)
    products, total, err := services.FindProducts(filterDto)
    if err != nil {
        return c.Error(response.Error{
            Code:    core.StatusInternalServerError,
            Message: "Failed to retrieve products",
        })
    }

    metadata := response.Meta{
        Page:    filterDto.Page,
        PerPage: filterDto.PerPage,
        Total:   total,
    }

    data := transformers.ToListResponse(products, transformers.ToProductResponse)

    return c.JSON(response.ListProduct{
        Meta: metadata,
        Data: data,
    })
}
```

### 8. Routes

```go
// File: internal/http/routes/api_routes.go

productsGroup := apiRouter.Group("/products")

// CRUD routes
productsGroup.POST("", app.Apply(
    middleware.JWTAuth(),
    middleware.CheckRolesMiddleware(types.RoleAdmin),
)(product.NewCreateProductApi()))

productsGroup.GET("/:id", app.Apply(
    middleware.JWTAuth(),
)(product.NewGetProductApi()))

productsGroup.PUT("/:id", app.Apply(
    middleware.JWTAuth(),
    middleware.CheckRolesMiddleware(types.RoleAdmin),
)(product.NewUpdateProductApi()))

productsGroup.DELETE("/:id", app.Apply(
    middleware.JWTAuth(),
    middleware.CheckRolesMiddleware(types.RoleAdmin),
)(product.NewDeleteProductApi()))

productsGroup.GET("", app.Apply(
    middleware.JWTAuth(),
)(product.NewListProductsApi()))
```

---

## Best Practices

### 1. HTTP Helper Usage

✅ **DO:**
- Use `ProcessData[T]` for CREATE operations (POST)
- Use `ProcessUpdateData[T]` for UPDATE operations (PUT/PATCH)
- Use `ProcessFilter` for LIST operations (GET with query params)
- Use `ProcessPathID` for READ/DELETE operations (GET/DELETE with path ID) - **RECOMMENDED**
- Let helpers handle parsing, sanitization, and validation
- Use one-line validation methods whenever possible

❌ **DON'T:**
- Manually parse JSON bodies (`c.BodyParser`)
- Manually extract and validate path IDs (use `ProcessPathID` instead)
- Skip sanitization (security risk)
- Skip validation (data integrity risk)
- Reinvent the wheel (helpers exist for a reason)
- Write multi-line validation logic when helpers exist

### 2. Request → DTO Conversion

✅ **DO:**
- Always implement `ToDto()` method on requests
- Use pointer receiver for `SetID()` in update requests
- Keep request objects thin (just embedding + conversion)
- Store validated requests in context

❌ **DON'T:**
- Pass requests to service layer (pass DTOs)
- Add business logic to requests
- Skip `ToDto()` and pass request directly

### 3. Model/DTO → Response Transformation

✅ **DO:**
- Always use transformers (never return models directly)
- Handle null values properly (`sql.NullX` → `*Type`)
- Use `ToListResponse` for collections
- Format URLs, dates, and other presentation concerns
- Enrich with related data when needed

❌ **DON'T:**
- Return domain models as JSON
- Expose internal fields (passwords, tokens, etc.)
- Transform in controllers (use transformer functions)
- Forget to handle nullable fields

### 4. Response Format

✅ **DO:**
- Return single resources directly for GET/POST/PUT
- Wrap lists in `{meta, data}` structure
- Use consistent field naming (camelCase or snake_case)
- Include pagination metadata for lists
- Document all response fields with `doc` tags

❌ **DON'T:**
- Wrap single resources in unnecessary objects
- Mix response formats across endpoints
- Return raw arrays (use meta + data)
- Skip documentation

### 5. HTTP Status Codes

✅ **DO:**
- Return 201 for POST (resource created)
- Return 200 for successful GET/PUT/PATCH
- Return 204 or 200 for DELETE
- Return 400 for validation errors
- Return 404 for missing resources
- Return 409 for conflicts (e.g., duplicate email)

❌ **DON'T:**
- Return 200 for all successes (be specific)
- Return 500 for validation errors (use 400)
- Return 200 for errors (use appropriate error codes)

### 6. Error Handling

✅ **DO:**
- Return structured error responses
- Include field-level validation errors in `data`
- Use appropriate status codes
- Log errors for debugging (don't expose to client)
- Use `github.com/gflydev/core/errors` package

❌ **DON'T:**
- Expose stack traces to clients
- Return vague error messages
- Use `fmt.Errorf` (use gFly's error package)
- Mix error response formats

### 7. Validation

✅ **DO:**
- Validate at HTTP boundary (controller's Validate())
- Use comprehensive validation rules
- Validate business rules in service layer
- Return clear validation error messages

❌ **DON'T:**
- Skip validation
- Validate in service layer (field validation belongs in DTO)
- Trust user input

### 8. Service Layer

✅ **DO:**
- Accept DTOs as input
- Return domain models as output
- Implement business logic
- Use repository pattern for data access
- Return meaningful errors

❌ **DON'T:**
- Accept HTTP context or requests
- Return response objects
- Access database directly
- Swallow errors

### 9. Swagger Documentation

✅ **DO:**
- Add complete Swagger annotations to all handlers
- Document all parameters (@Param)
- Document all responses (@Success, @Failure)
- Run `make doc` after changes
- Test endpoints in Swagger UI

❌ **DON'T:**
- Skip Swagger annotations
- Leave outdated documentation
- Forget to regenerate docs

### 10. Route Organization

✅ **DO:**
- Group related routes
- Apply middleware consistently
- Use semantic endpoint names
- Follow RESTful conventions

❌ **DON'T:**
- Mix authentication strategies
- Use inconsistent route patterns
- Create nested routes unnecessarily

---

## Summary

This guide provides a comprehensive reference for building CRUD APIs in the gFly framework:

1. **HTTP Helpers** simplify request processing with type-safe validation
2. **Request → DTO** conversion keeps HTTP and business layers separate
3. **Transformers** handle Model/DTO → Response conversion and formatting
4. **Response Patterns** provide consistent, well-documented API outputs
5. **Status Codes** communicate operation results clearly
6. **Error Handling** delivers actionable feedback to clients

By following these patterns, you'll create maintainable, secure, and well-documented REST APIs that align with gFly's Clean Architecture principles.

---

## Quick Reference

### CRUD Operation Checklist

- [ ] Define DTO in `/internal/dto/`
- [ ] Define Request in `/internal/http/request/`
- [ ] Define Response in `/internal/http/response/`
- [ ] Create Transformer in `/internal/http/transformers/`
- [ ] Implement Service in `/internal/services/`
- [ ] Create Controller in `/internal/http/controllers/api/`
- [ ] Register Route in `/internal/http/routes/api_routes.go`
- [ ] Add Swagger annotations to controller
- [ ] Run `make doc` to generate documentation
- [ ] Test all endpoints

### Common HTTP Helpers

```go
// CREATE
http.ProcessData[request.CreateX](c)

// UPDATE
http.ProcessUpdateData[request.UpdateX](c)

// LIST
http.ProcessFilter(c)

// GET/DELETE
http.PathID(c)
```

### Common Transformer Patterns

```go
// Single resource
transformers.ToXResponse(model)

// List
transformers.ToListResponse(models, transformers.ToXResponse)
```

### Common Response Patterns

```go
// 201 Created
c.Status(core.StatusCreated).JSON(response)

// 200 OK
c.JSON(response)

// 204 No Content
c.Status(core.StatusNoContent).Send(nil)

// 400 Bad Request
c.Error(response.Error{Code: 400, Message: "..."})

// 404 Not Found
c.Error(response.NotFound{Code: 404, Message: "..."})
```

---

**Last Updated:** 2025-01-22
**gFly Framework Version:** Compatible with github.com/gflydev/core v1.x

**Related Guides:**
- [Data Flow Guide](DATA_FLOW_GUIDE.md) - Complete data flow patterns
- [Model Builder Guide](MODEL_BUILDER_GUIDE.md) - Database query patterns
