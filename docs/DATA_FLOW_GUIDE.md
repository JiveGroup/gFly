# gFly Framework: Request, DTO, Response & Transformer Guide

## Table of Contents
1. [Architecture Overview](#architecture-overview)
2. [DTO (Data Transfer Objects)](#dto-data-transfer-objects)
3. [Request Objects](#request-objects)
4. [Response Objects](#response-objects)
5. [Transformer Objects](#transformer-objects)
6. [Validation System](#validation-system)
7. [Complete Flow Examples](#complete-flow-examples)
8. [Best Practices](#best-practices)

---

## Architecture Overview

The gFly framework follows a **layered architecture** with clear separation of concerns. Data flows through distinct layers, each with its own responsibility:

```
┌─────────────────────────────────────────────────────────────┐
│                 HTTP Request (JSON/FORM)                    │
└────────────────────────────┬────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────┐
│  REQUEST LAYER (internal/http/request/)                     │
│  • Parses HTTP request                                      │
│  • Embeds DTO for validation                                │
│  • Implements ToDto() conversion                            │
└────────────────────────────┬────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────┐
│  DTO / Data (internal/dto/)                                 │
│  • Core data structure                                      │
│  • Validation rules via tags                                │
│  • Documentation via tags                                   │
│  • Single source of truth                                   │
└────────────────────────────┬────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────┐
│  SERVICE LAYER (internal/services/)                         │
│  • Receives DTO from controller                             │
│  • Executes business logic                                  │
│  • Uses repository for data access                          │
│  • Returns Domain Model / DTO                               │
└────────────────────────────┬────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────┐
│  DOMAIN MODEL (internal/domain/models/)                     │
│  • Database entity representation                           │
│  • Pure data structure                                      │
│  • Uses db tags for ORM mapping                             │
└────────────────────────────┬────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────┐
│  TRANSFORMER LAYER (internal/http/transformers/)            │
│  • Converts Domain Model / DTO → Response                   │
│  • Handles null values                                      │
│  • Formats URLs, dates, etc.                                │
│  • Enriches with related data                               │
└────────────────────────────┬────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────┐
│  RESPONSE DATA (internal/http/response/)                    │
│  • Final JSON structure                                     │
│  • Documentation via tags                                   │
│  • Used in Swagger/OpenAPI                                  │
└────────────────────────────┬────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────┐
│                     HTTP Response (JSON)                    │
└─────────────────────────────────────────────────────────────┘
```

### Key Principles

1. **Separation of Concerns**: Each layer has a single, well-defined responsibility
2. **Data Flow Direction**: Request → DTO → Service → Model / DTO → Transformer → Response
3. **Validation at Boundary**: Input validation happens at the HTTP boundary (DTO layer)
4. **Business Logic Isolation**: Services operate on DTOs, not HTTP requests
5. **Presentation Separation**: Transformers handle all response formatting concerns

---

## DTO (Data Transfer Objects)

### Purpose

DTOs are the **single source of truth** for data validation and structure. They:
- Define the contract between HTTP layer and business logic
- Contain validation rules via struct tags
- Provide documentation for API consumers
- Serve as input to service layer functions

### Location

- **Central DTOs**: `/internal/dto/` (shared across application)
- **Module-specific DTOs**: `/pkg/modules/{module}/dto/`

### Structure Pattern

```go
package dto

import "gfly/internal/domain/models/types"

// CreateUser struct to describe the request body to create a new user.
// @Description Request payload for creating a new user.
// @Tags Users
type CreateUser struct {
    Email    string       `json:"email" example:"john@example.com" validate:"required,email,max=255" doc:"User's email address (required, max length 255)"`
    Password string       `json:"password" example:"SecureP@ss123" validate:"required,min=8,max=255" doc:"User's password (required, 8-255 chars)"`
    Fullname string       `json:"fullname" example:"John Doe" validate:"required,max=255" doc:"User's full name (required, max length 255)"`
    Phone    string       `json:"phone" example:"0989831911" validate:"required,max=20" doc:"User's phone number (required, max length 20)"`
    Avatar   string       `json:"avatar" example:"https://i.pravatar.cc/32" validate:"omitempty,max=255" doc:"URL of user's avatar (optional)"`
    Status   string       `json:"status" example:"active" validate:"omitempty,oneof=active pending blocked" doc:"User's status (optional)"`
    Roles    []types.Role `json:"roles" example:"admin,user" validate:"omitempty" doc:"List of user's roles (optional)"`
}

// UpdateUser struct to partially update an existing user.
// @Description Request payload for updating an existing user.
// @Tags Users
type UpdateUser struct {
    ID       int          `json:"-" validate:"omitempty,gte=1" doc:"User ID (auto-populated from URL path)"`
    Password string       `json:"password" example:"NewP@ss456" validate:"omitempty,min=8,max=255" doc:"New password (optional)"`
    Fullname string       `json:"fullname" example:"Jane Doe" validate:"omitempty,max=255" doc:"Updated full name (optional)"`
    Phone    string       `json:"phone" example:"0123456789" validate:"omitempty,max=20" doc:"Updated phone (optional)"`
    Avatar   string       `json:"avatar" validate:"omitempty,max=255" doc:"Updated avatar URL (optional)"`
    Roles    []types.Role `json:"roles" validate:"omitempty" doc:"Updated roles (optional)"`
}

// UpdateUserStatus struct for status-only updates.
type UpdateUserStatus struct {
    ID     int              `json:"-" validate:"omitempty" doc:"User ID (auto-populated from URL path)"`
    Status types.UserStatus `json:"status" example:"active" validate:"required,oneof=active pending blocked" doc:"New status (required)"`
}

// Filter struct for list/search operations.
type Filter struct {
    Keyword string `json:"keyword" validate:"omitempty,max=255" doc:"Search keyword (optional)"`
    OrderBy string `json:"order_by" validate:"omitempty" doc:"Sort field (optional)"`
    Page    int    `json:"page" validate:"omitempty,gte=1" doc:"Page number (default: 1)"`
    PerPage int    `json:"per_page" validate:"omitempty,gte=1,lte=100" doc:"Items per page (default: 10, max: 100)"`
}
```

### Struct Tags Reference

| Tag | Purpose | Examples |
|-----|---------|----------|
| `json:"field_name"` | JSON serialization/deserialization | `json:"email"`, `json:"full_name"`, `json:"-"` (exclude) |
| `validate:"rules"` | Validation rules (comma-separated) | `validate:"required,email,max=255"` |
| `example:"value"` | Swagger example value | `example:"john@example.com"` |
| `doc:"description"` | Field documentation | `doc:"User's email address"` |
| `@Description` | Struct-level Swagger description | (at comment level) |
| `@Tags` | Swagger tag grouping | (at comment level) |

### Common Validation Rules

| Rule | Description | Example |
|------|-------------|---------|
| `required` | Field must be present and non-zero | `validate:"required"` |
| `omitempty` | Field is optional | `validate:"omitempty"` |
| `email` | Must be valid email format | `validate:"email"` |
| `min=N` | Minimum length (string) or value (number) | `validate:"min=8"` |
| `max=N` | Maximum length (string) or value (number) | `validate:"max=255"` |
| `gte=N` | Greater than or equal to | `validate:"gte=1"` |
| `lte=N` | Less than or equal to | `validate:"lte=100"` |
| `oneof=v1 v2 v3` | Must be one of the listed values | `validate:"oneof=active pending blocked"` |
| `len=N` | Exact length | `validate:"len=10"` |
| `url` | Must be valid URL | `validate:"url"` |
| `uuid` | Must be valid UUID | `validate:"uuid"` |
| `alphanum` | Alphanumeric characters only | `validate:"alphanum"` |

Refer [validation rules](https://doc.gfly.dev/02-basic/02-01-11.validation/)

### Module-Specific DTO Example

```go
// File: pkg/modules/auth/dto/auth_dto.go
package dto

// SignUp struct to describe user registration.
type SignUp struct {
    Email    string `json:"email" example:"john@example.com" validate:"required,email,max=255" doc:"User's email address"`
    Password string `json:"password" example:"SecureP@ss123" validate:"required,min=8,max=255" doc:"User's password"`
    Fullname string `json:"fullname" example:"John Doe" validate:"required,max=255" doc:"User's full name"`
    Phone    string `json:"phone" example:"0989831911" validate:"required,max=20" doc:"User's phone number"`
    Avatar   string `json:"avatar" validate:"omitempty,max=255" doc:"Avatar URL (optional)"`
}

// SignIn struct to describe user login.
type SignIn struct {
    Username string `json:"username" example:"admin@gfly.dev" validate:"required,email,max=255" doc:"Email or username"`
    Password string `json:"password" example:"P@ssw0rd" validate:"required,max=255" doc:"User password"`
}

// RefreshToken struct for token refresh.
type RefreshToken struct {
    Token string `json:"token" validate:"required,max=255" doc:"Refresh token for obtaining new access token"`
}
```

---

## Request Objects

> **⚠️ CRITICAL REQUIREMENT**:
>
> **Every Request struct MUST implement `ToDto()` method** to satisfy the `AddData[D]` or `UpdateData[D]` interface. This is NOT optional!
>
> - **Single DTO**: `ToDto()` returns the embedded DTO
> - **Multiple DTOs**: `ToDto()` returns primary/default DTO + implement `ToDtoXXX()` methods for alternatives

### Purpose

Request objects are **thin wrappers** around DTOs that:
- Handle HTTP-specific concerns (parsing, binding)
- **Implement required interfaces** (`AddData[D]` or `UpdateData[D]`) for validation pipeline
- Provide conversion methods to DTOs via **required `ToDto()` method**
- Keep HTTP layer separate from business logic

### Location

- **Central Requests**: `/internal/http/request/`
- **Module-specific Requests**: `/pkg/modules/{module}/request/`

### Pattern: Embedding DTO

```go
// File: internal/http/request/user_request.go
package request

import "gfly/internal/dto"

// ====================================================================
// ========================== Add Requests ============================
// ====================================================================

// CreateUser request wraps the CreateUser DTO.
type CreateUser struct {
    dto.CreateUser  // Embedded - inherits all fields and validation
}

// ToDto converts Request to DTO for service layer.
func (r CreateUser) ToDto() dto.CreateUser {
    return r.CreateUser
}

// ====================================================================
// ========================= Update Requests ==========================
// ====================================================================

// UpdateUser request wraps the UpdateUser DTO.
type UpdateUser struct {
    dto.UpdateUser  // Embedded DTO
}

// ToDto converts Request to DTO.
func (r UpdateUser) ToDto() dto.UpdateUser {
    return r.UpdateUser
}

// SetID sets the ID from URL path parameter.
// IMPORTANT: Must use pointer receiver for SetID to modify the struct!
func (r *UpdateUser) SetID(id int) {
    r.ID = id
}

// ====================================================================
// ===================== Status Update Requests =======================
// ====================================================================

// UpdateUserStatus request for status changes only.
type UpdateUserStatus struct {
    dto.UpdateUserStatus
}

// ToDto converts to DTO.
func (r UpdateUserStatus) ToDto() dto.UpdateUserStatus {
    return r.UpdateUserStatus
}

// SetID populates ID from path parameter.
func (r *UpdateUserStatus) SetID(id int) {
    r.ID = id
}
```

### Required Interfaces

> **⚠️ CRITICAL**: Every Request struct **MUST** implement these interfaces to work with gFly's validation pipeline.

```go
// AddData interface for create operations (POST)
type AddData[D any] interface {
    ToDto() D  // ⚠️ REQUIRED: Must return corresponding DTO
}

// UpdateData interface for update operations (PUT/PATCH)
type UpdateData[D any] interface {
    ToDto() D         // ⚠️ REQUIRED: Must return corresponding DTO
    SetID(id int)     // ⚠️ REQUIRED: Must populate ID from URL path parameter
}
```

**Key Points:**
- `ToDto()` is **NOT optional** - it's required by the interface contract
- The return type `D` is the generic type parameter (the DTO type)
- For single DTO requests: `ToDto()` directly returns the embedded DTO
- For multiple DTO requests: `ToDto()` returns the primary/default DTO (and you add `ToDtoXXX()` methods for alternatives)
- `SetID()` must use a **pointer receiver** (`*Request`) to modify the struct

### Module Request Example

```go
// File: pkg/modules/auth/request/auth_request.go
package request

import "gfly/pkg/modules/auth/dto"

// SignIn request wraps SignIn DTO.
type SignIn struct {
    dto.SignIn
}

// ToDto converts to SignIn DTO.
func (r SignIn) ToDto() dto.SignIn {
    return r.SignIn
}

// SignUp request wraps SignUp DTO.
type SignUp struct {
    dto.SignUp
}

// ToDto converts to SignUp DTO.
func (r SignUp) ToDto() dto.SignUp {
    return r.SignUp
}

// RefreshToken request for token refresh.
type RefreshToken struct {
    dto.RefreshToken
}

// ToDto converts to RefreshToken DTO.
func (r RefreshToken) ToDto() dto.RefreshToken {
    return r.RefreshToken
}
```

### Pattern: Multiple DTOs in One Request

> **⚠️ CRITICAL REQUIREMENT**: All Request structs **MUST** implement the `ToDto()` method to satisfy the `AddData[D]` or `UpdateData[D]` interface. This is non-negotiable. When using multiple DTOs, implement `ToDto()` to return your primary/default DTO, and add additional `ToDtoXXX()` methods for alternative conversions.

A powerful pattern in gFly is when a single Request struct embeds **multiple DTOs**, each representing a different context or use case. This allows one Request to be transformed into different DTOs depending on the operation.

#### Pattern Overview

```go
// Request with multiple DTOs
type MyRequest struct {
    dto.PrimaryDTO     // Primary/default DTO
    dto.AlternativeDTO // Alternative DTO for different context
}

// ⚠️ REQUIRED: Implement ToDto() for the interface
func (r MyRequest) ToDto() dto.PrimaryDTO {
    return r.ToDtoPrimary()  // Delegate to primary conversion
}

// Primary conversion method
func (r MyRequest) ToDtoPrimary() dto.PrimaryDTO {
    return r.PrimaryDTO  // or construct it manually
}

// Alternative conversion method
func (r MyRequest) ToDtoAlternative() dto.AlternativeDTO {
    return r.AlternativeDTO  // or construct it manually
}
```

**The Rule:**
- `ToDto()` = REQUIRED by interface (returns default/primary DTO)
- `ToDtoXXX()` = OPTIONAL additional methods (return alternative DTOs)
- `ToDto()` usually delegates to one of your `ToDtoXXX()` methods

#### When to Use Multiple DTOs

Use this pattern when:
1. **Different contexts require different validation** (admin vs. user operations)
2. **Same HTTP endpoint serves multiple purposes** (different service methods)
3. **Partial updates with different field sets** (status update vs. full update)
4. **Role-based operations** (admin can set more fields than users)

#### Example 1: Admin vs. User User Creation

```go
// File: internal/dto/user_dto.go
package dto

// CreateUser - Standard user creation (limited fields)
type CreateUser struct {
    Email    string `json:"email" validate:"required,email,max=255" doc:"User's email"`
    Password string `json:"password" validate:"required,min=8,max=255" doc:"User's password"`
    Fullname string `json:"fullname" validate:"required,max=255" doc:"User's full name"`
    Phone    string `json:"phone" validate:"required,max=20" doc:"User's phone"`
}

// CreateUserByAdmin - Admin can set additional fields (status, roles)
type CreateUserByAdmin struct {
    Email    string       `json:"email" validate:"required,email,max=255" doc:"User's email"`
    Password string       `json:"password" validate:"required,min=8,max=255" doc:"User's password"`
    Fullname string       `json:"fullname" validate:"required,max=255" doc:"User's full name"`
    Phone    string       `json:"phone" validate:"required,max=20" doc:"User's phone"`
    Avatar   string       `json:"avatar" validate:"omitempty,url,max=255" doc:"Avatar URL"`
    Status   string       `json:"status" validate:"omitempty,oneof=active pending blocked" doc:"User status"`
    Roles    []types.Role `json:"roles" validate:"omitempty" doc:"Assigned roles"`
}
```

```go
// File: internal/http/request/user_request.go
package request

import "gfly/internal/dto"

// CreateUserRequest can be used in different contexts
type CreateUserRequest struct {
    dto.CreateUser         // Basic user creation fields
    dto.CreateUserByAdmin  // Admin-specific fields (status, roles)
}

// ⚠️ IMPORTANT: ToDto() is REQUIRED to satisfy the AddData interface
// This is the default/primary conversion method
func (r CreateUserRequest) ToDto() dto.CreateUser {
    return r.ToDtoUser()
}

// ToDtoUser converts to basic CreateUser DTO (for regular user registration)
func (r CreateUserRequest) ToDtoUser() dto.CreateUser {
    return dto.CreateUser{
        Email:    r.Email,
        Password: r.Password,
        Fullname: r.Fullname,
        Phone:    r.Phone,
    }
}

// ToDtoAdmin converts to CreateUserByAdmin DTO (for admin operations)
func (r CreateUserRequest) ToDtoAdmin() dto.CreateUserByAdmin {
    return dto.CreateUserByAdmin{
        Email:    r.Email,
        Password: r.Password,
        Fullname: r.Fullname,
        Phone:    r.Phone,
        Avatar:   r.Avatar,
        Status:   r.Status,
        Roles:    r.Roles,
    }
}
```

**Usage in Controllers:**

```go
// User registration endpoint (limited fields)
func (h *RegisterUserApi) Handle(c *core.Ctx) error {
    requestData := c.GetData(constants.Request).(request.CreateUserRequest)

    // Convert to basic DTO (ignores admin-only fields)
    createUserDto := requestData.ToDtoUser()

    user, err := services.RegisterUser(createUserDto)
    // ...
}

// Admin create user endpoint (all fields)
func (h *AdminCreateUserApi) Handle(c *core.Ctx) error {
    requestData := c.GetData(constants.Request).(request.CreateUserRequest)

    // Convert to admin DTO (includes status, roles)
    createUserDto := requestData.ToDtoAdmin()

    user, err := services.CreateUserByAdmin(createUserDto)
    // ...
}
```

#### Example 2: Product Update with Different Scopes

```go
// File: internal/dto/product_dto.go
package dto

// UpdateProductInfo - Update basic product information
type UpdateProductInfo struct {
    ID          int     `json:"-" validate:"omitempty,gte=1"`
    Name        string  `json:"name" validate:"omitempty,max=255" doc:"Product name"`
    Description string  `json:"description" validate:"omitempty,max=1000" doc:"Product description"`
    ImageURL    string  `json:"image_url" validate:"omitempty,url" doc:"Product image URL"`
}

// UpdateProductInventory - Update inventory-related fields
type UpdateProductInventory struct {
    ID    int     `json:"-" validate:"omitempty,gte=1"`
    Price float64 `json:"price" validate:"omitempty,gte=0" doc:"Product price"`
    Stock int     `json:"stock" validate:"omitempty,gte=0" doc:"Stock quantity"`
}

// UpdateProductStatus - Update product status only
type UpdateProductStatus struct {
    ID     int    `json:"-" validate:"omitempty,gte=1"`
    Status string `json:"status" validate:"required,oneof=active inactive discontinued" doc:"Product status"`
}
```

```go
// File: internal/http/request/product_request.go
package request

import "gfly/internal/dto"

// UpdateProductRequest supports multiple update operations
type UpdateProductRequest struct {
    dto.UpdateProductInfo      // Basic info (name, description, image)
    dto.UpdateProductInventory // Inventory (price, stock)
    dto.UpdateProductStatus    // Status
}

// ⚠️ IMPORTANT: ToDto() is REQUIRED to satisfy the UpdateData interface
// This returns the default/primary DTO (product info in this case)
func (r UpdateProductRequest) ToDto() dto.UpdateProductInfo {
    return r.ToDtoInfo()
}

// ToDtoInfo converts to UpdateProductInfo DTO
func (r UpdateProductRequest) ToDtoInfo() dto.UpdateProductInfo {
    return dto.UpdateProductInfo{
        ID:          r.ID,
        Name:        r.Name,
        Description: r.Description,
        ImageURL:    r.ImageURL,
    }
}

// ToDtoInventory converts to UpdateProductInventory DTO
func (r UpdateProductRequest) ToDtoInventory() dto.UpdateProductInventory {
    return dto.UpdateProductInventory{
        ID:    r.ID,
        Price: r.Price,
        Stock: r.Stock,
    }
}

// ToDtoStatus converts to UpdateProductStatus DTO
func (r UpdateProductRequest) ToDtoStatus() dto.UpdateProductStatus {
    return dto.UpdateProductStatus{
        ID:     r.ID,
        Status: r.Status,
    }
}

// SetID sets the ID from URL path parameter (required by UpdateData interface)
func (r *UpdateProductRequest) SetID(id int) {
    r.ID = id
}
```

**Usage in Controllers:**

```go
// Update product info only
func (h *UpdateProductInfoApi) Handle(c *core.Ctx) error {
    requestData := c.GetData(constants.Request).(request.UpdateProductRequest)

    // Convert to info DTO (name, description, image)
    updateDto := requestData.ToDtoInfo()

    product, err := services.UpdateProductInfo(updateDto)
    // ...
}

// Update inventory only (price/stock)
func (h *UpdateProductInventoryApi) Handle(c *core.Ctx) error {
    requestData := c.GetData(constants.Request).(request.UpdateProductRequest)

    // Convert to inventory DTO (price, stock)
    updateDto := requestData.ToDtoInventory()

    product, err := services.UpdateProductInventory(updateDto)
    // ...
}

// Update status only
func (h *UpdateProductStatusApi) Handle(c *core.Ctx) error {
    requestData := c.GetData(constants.Request).(request.UpdateProductRequest)

    // Convert to status DTO
    updateDto := requestData.ToDtoStatus()

    product, err := services.UpdateProductStatus(updateDto)
    // ...
}
```

#### Example 3: Order Operations with Different DTOs

```go
// File: internal/dto/order_dto.go
package dto

// CreateOrder - Create new order with items
type CreateOrder struct {
    UserID      int         `json:"user_id" validate:"required,gte=1" doc:"User ID"`
    Items       []OrderItem `json:"items" validate:"required,min=1" doc:"Order items"`
    AddressID   int         `json:"address_id" validate:"required,gte=1" doc:"Shipping address ID"`
    PaymentType string      `json:"payment_type" validate:"required,oneof=cod card" doc:"Payment type"`
}

// UpdateOrderStatus - Update order status (for admin/system)
type UpdateOrderStatus struct {
    ID     int    `json:"-" validate:"omitempty,gte=1"`
    Status string `json:"status" validate:"required,oneof=pending confirmed shipping delivered cancelled" doc:"Order status"`
    Note   string `json:"note" validate:"omitempty,max=500" doc:"Status change note"`
}

// CancelOrder - Cancel order (for customer)
type CancelOrder struct {
    ID     int    `json:"-" validate:"omitempty,gte=1"`
    Reason string `json:"reason" validate:"required,max=500" doc:"Cancellation reason"`
}
```

```go
// File: internal/http/request/order_request.go
package request

import "gfly/internal/dto"

// OrderRequest supports create, status update, and cancellation
type OrderRequest struct {
    dto.CreateOrder       // For creating new orders
    dto.UpdateOrderStatus // For admin status updates
    dto.CancelOrder       // For customer cancellations
}

// ⚠️ IMPORTANT: ToDto() is REQUIRED to satisfy the AddData/UpdateData interface
// This returns the default/primary DTO (CreateOrder for POST requests)
func (r OrderRequest) ToDto() dto.CreateOrder {
    return r.ToDtoCreate()
}

// ToDtoCreate converts to CreateOrder DTO
func (r OrderRequest) ToDtoCreate() dto.CreateOrder {
    return r.CreateOrder
}

// ToDtoUpdateStatus converts to UpdateOrderStatus DTO
func (r OrderRequest) ToDtoUpdateStatus() dto.UpdateOrderStatus {
    return dto.UpdateOrderStatus{
        ID:     r.ID,
        Status: r.Status,
        Note:   r.Note,
    }
}

// ToDtoCancel converts to CancelOrder DTO
func (r OrderRequest) ToDtoCancel() dto.CancelOrder {
    return dto.CancelOrder{
        ID:     r.ID,
        Reason: r.Reason,
    }
}

// SetID sets the ID from URL path parameter
func (r *OrderRequest) SetID(id int) {
    r.ID = id
}
```

**Usage in Controllers:**

```go
// Create order endpoint
func (h *CreateOrderApi) Handle(c *core.Ctx) error {
    requestData := c.GetData(constants.Request).(request.OrderRequest)
    order, err := services.CreateOrder(requestData.ToDtoCreate())
    // ...
}

// Admin update order status endpoint
func (h *UpdateOrderStatusApi) Handle(c *core.Ctx) error {
    requestData := c.GetData(constants.Request).(request.OrderRequest)
    order, err := services.UpdateOrderStatus(requestData.ToDtoUpdateStatus())
    // ...
}

// Customer cancel order endpoint
func (h *CancelOrderApi) Handle(c *core.Ctx) error {
    requestData := c.GetData(constants.Request).(request.OrderRequest)
    order, err := services.CancelOrder(requestData.ToDtoCancel())
    // ...
}
```

#### Benefits of Multiple DTOs Pattern

1. **Code Reusability**: One Request struct serves multiple endpoints
2. **Flexible Validation**: Different DTOs enforce different validation rules
3. **Clear Intent**: `ToDtoAdmin()`, `ToDtoUser()` make the context explicit
4. **Type Safety**: Each service method receives the exact DTO it expects
5. **Maintainability**: Changes to shared fields update all DTOs automatically
6. **Separation of Concerns**: HTTP layer (Request) separate from business logic (DTO)

#### Important Notes

⚠️ **CRITICAL: `ToDto()` is ALWAYS REQUIRED**
- **Every Request MUST implement `ToDto()`** to satisfy the `AddData` or `UpdateData` interface
- When embedding multiple DTOs, `ToDto()` should return the **default/primary DTO**
- `ToDto()` can delegate to one of your specific `ToDtoXXX()` methods (e.g., `return r.ToDtoUser()`)
- The specific `ToDtoXXX()` methods are **in addition to** `ToDto()`, not a replacement

**Additional Guidelines:**
- **Always implement specific `ToDtoXXX()` methods** for each embedded DTO to make intent explicit
- **Choose a sensible default** for `ToDto()` - typically the most common use case
- **Field conflicts**: If multiple DTOs have the same field name, Go embedding rules apply (first embedded struct wins)
- **Interface compatibility**: The `ToDto()` return type must match the interface's generic type parameter `D`

### Why Separate Request from DTO?

1. **Separation of Concerns**: HTTP layer vs. business logic layer
2. **Interface Implementation**: Requests implement HTTP-specific interfaces
3. **Flexibility**: Can add HTTP-specific methods without polluting DTOs
4. **Testing**: Services can be tested with DTOs without HTTP context
5. **Reusability**: Same DTO can be used in different contexts (API, CLI, queue jobs)

---

## Response Objects

### Purpose

Response objects define the **output structure** returned to HTTP clients. They:
- Define JSON serialization structure
- Provide documentation for API consumers
- Handle presentation concerns (formatting, nullability)
- Serve as Swagger/OpenAPI schema definitions

### Location

- **Central Responses**: `/internal/http/response/`
- **Module-specific Responses**: `/pkg/modules/{module}/response/`

### Structure Pattern

```go
// File: internal/http/response/user_response.go
package response

import (
    "gfly/internal/domain/models/types"
    "time"
)

// User struct to describe User response.
// Created via transformers from models.User
type User struct {
    ID           int              `json:"id" doc:"Unique user identifier"`
    Email        string           `json:"email" doc:"User's email address"`
    Fullname     string           `json:"fullname" doc:"User's full name"`
    Phone        string           `json:"phone" doc:"User's phone number"`
    Token        *string          `json:"token" doc:"Authentication token (nullable)"`
    Status       types.UserStatus `json:"status" doc:"User account status"`
    CreatedAt    time.Time        `json:"created_at" doc:"Account creation timestamp"`
    UpdatedAt    time.Time        `json:"updated_at" doc:"Last update timestamp"`
    VerifiedAt   *time.Time       `json:"verified_at" example:"2023-01-01T10:30:00Z" doc:"Email verification timestamp (nullable)"`
    BlockedAt    *time.Time       `json:"blocked_at" example:"null" doc:"Account block timestamp (nullable)"`
    DeletedAt    *time.Time       `json:"deleted_at" example:"null" doc:"Soft delete timestamp (nullable)"`
    LastAccessAt *time.Time       `json:"last_access_at" doc:"Last login timestamp (nullable)"`
    Avatar       *string          `json:"avatar" doc:"Avatar image URL (nullable)"`
    Roles        []Role           `json:"roles" doc:"Assigned user roles"`
}

// Role struct to describe Role response.
type Role struct {
    ID   int        `json:"id" doc:"Unique role identifier"`
    Name string     `json:"name" doc:"Display name of the role"`
    Slug types.Role `json:"slug" doc:"URL-friendly role identifier"`
}

// ListUser response for paginated user lists.
type ListUser struct {
    Meta Meta   `json:"meta" doc:"Pagination metadata"`
    Data []User `json:"data" doc:"List of user objects"`
}
```

### Generic Response Patterns

```go
// File: internal/http/response/generic_response.go
package response

import "github.com/gflydev/core"

// ====================================================================
// ======================== Success Responses =========================
// ====================================================================

// Meta struct for pagination metadata.
type Meta struct {
    Page    int `json:"page,omitempty" example:"1" doc:"Current page number"`
    PerPage int `json:"per_page,omitempty" example:"10" doc:"Items per page"`
    Total   int `json:"total" example:"1354" doc:"Total number of records"`
}

// List generic response for paginated lists.
type List[T any] struct {
    Meta Meta `json:"meta" doc:"Pagination metadata"`
    Data []T  `json:"data" doc:"List of items"`
}

// Success generic success response.
type Success struct {
    Message string    `json:"message" example:"Operation completed successfully"`
    Data    core.Data `json:"data" doc:"Additional operation data"`
}

// ServerInfo response for system information.
type ServerInfo struct {
    Name   string `json:"name" example:"ThietNgon API" doc:"API name"`
    Prefix string `json:"prefix" example:"/api/v1" doc:"API prefix including version"`
    Server string `json:"server" example:"ThietNgon-Go Server" doc:"Server application name"`
}

// ====================================================================
// ========================= Error Responses ==========================
// ====================================================================

// Error generic error response.
type Error struct {
    Code    int       `json:"code" example:"400"`
    Message string    `json:"message" example:"Bad request"`
    Data    core.Data `json:"data" doc:"Field-level validation errors (if applicable)"`
}

// Unauthorized 401 error response.
type Unauthorized struct {
    Code    int    `json:"code" example:"401"`
    Message string `json:"error" example:"Unauthorized access"`
}

// NotFound 404 error response.
type NotFound struct {
    Code    int    `json:"code" example:"404"`
    Message string `json:"error" example:"Resource not found"`
}

// Conflict 409 error response.
type Conflict struct {
    Code    int    `json:"code" example:"409"`
    Message string `json:"error" example:"Resource conflict"`
}
```

### Handling Nullable Fields

Response objects use **pointers** for nullable fields:

```go
type User struct {
    ID       int       `json:"id"`           // Always present
    Email    string    `json:"email"`        // Always present
    Avatar   *string   `json:"avatar"`       // Nullable - pointer type
    DeletedAt *time.Time `json:"deleted_at"` // Nullable - pointer type
}
```

**JSON Output:**
```json
{
  "id": 123,
  "email": "user@example.com",
  "avatar": "https://example.com/avatar.jpg",  // Present
  "deleted_at": null                            // Explicitly null
}
```

or

```json
{
  "id": 123,
  "email": "user@example.com",
  "avatar": null,        // User has no avatar
  "deleted_at": null     // Not deleted
}
```

### Module Response Example

```go
// File: pkg/modules/auth/response/auth_response.go
package response

// SignIn response containing JWT tokens.
type SignIn struct {
    Access  string `json:"access" doc:"JWT access token for authentication"`
    Refresh string `json:"refresh" doc:"JWT refresh token for obtaining new access token"`
}
```

---

## Transformer Objects

### Purpose

Transformers convert **Domain Models → Response objects**. They:
- Handle null value conversion (`sql.NullString` → `*string`)
- Format URLs (relative paths → absolute URLs)
- Enrich responses with related data (e.g., user roles)
- Apply business logic for presentation (e.g., hide sensitive fields)

### Location

- **Central Transformers**: `/internal/http/transformers/`
- **Module Transformers**: `/pkg/modules/{module}/transformers/`

### Transformer Pattern

```go
// File: internal/http/transformers/user_transformer.go
package transformers

import (
    "gfly/internal/domain/models"
    "gfly/internal/domain/repository"
    "gfly/internal/http/response"
    "github.com/gflydev/core"
    dbNull "github.com/gflydev/db/null"
    "github.com/gflydev/storage"
    "strings"
)

// PublicAvatar converts avatar path to public URL.
//
// Parameters:
//   - avatar: Avatar file path or URL
//
// Returns:
//   - *string: Public avatar URL, nil if empty
func PublicAvatar(avatar string) *string {
    if avatar == "" {
        return nil
    }

    fs := storage.Instance()

    // Already absolute URL
    if strings.HasPrefix(avatar, core.SchemaHTTP) {
        return &avatar
    }

    // Convert relative path to absolute URL
    avatar = fs.Url(avatar)
    return &avatar
}

// ToRoleResponse converts Role model to Role response.
func ToRoleResponse(model models.Role) response.Role {
    return response.Role{
        ID:   model.ID,
        Name: model.Name,
        Slug: model.Slug,
    }
}

// roles retrieves and converts roles for a user.
func roles(userID int) []response.Role {
    var roles []response.Role
    roleList := repository.Pool.GetRolesByUserID(userID)

    for _, role := range roleList {
        roles = append(roles, ToRoleResponse(role))
    }

    return roles
}

// ToUserResponse converts User model to User response.
//
// Parameters:
//   - user: Domain User model
//
// Returns:
//   - response.User: Formatted user response object
func ToUserResponse(user models.User) response.User {
    return response.User{
        ID:           user.ID,
        Email:        user.Email,
        Fullname:     user.Fullname,
        Phone:        user.Phone,
        Token:        dbNull.StringNil(user.Token),        // sql.NullString → *string
        Status:       user.Status,
        Avatar:       PublicAvatar(user.Avatar.String),     // Path → URL
        CreatedAt:    user.CreatedAt,
        UpdatedAt:    user.UpdatedAt,
        VerifiedAt:   dbNull.TimeNil(user.VerifiedAt),     // sql.NullTime → *time.Time
        BlockedAt:    dbNull.TimeNil(user.BlockedAt),
        DeletedAt:    dbNull.TimeNil(user.DeletedAt),
        LastAccessAt: dbNull.TimeNil(user.LastAccessAt),
        Roles:        roles(user.ID),                       // Enrich with related data
    }
}

// ToSignUpResponse is a specialized transformer for signup responses.
func ToSignUpResponse(user models.User) response.User {
    // Could apply different formatting logic for signup vs. general responses
    return ToUserResponse(user)
}
```

### Generic Transformer Utilities

```go
// File: internal/http/transformers/generic_transformer.go
package transformers

import "github.com/gflydev/utils/fn"

// ToListResponse transforms a list of models to a list of responses.
//
// Generic Parameters:
//   - T: Model type (e.g., models.User)
//   - R: Response type (e.g., response.User)
//
// Parameters:
//   - records: Slice of domain models
//   - transformerFn: Function to transform single model → response
//
// Returns:
//   - []R: Slice of response objects
func ToListResponse[T any, R any](records []T, transformerFn func(T) R) []R {
    return fn.TransformList(records, transformerFn)
}
```

**Usage Example:**

```go
// In controller
users, total, err := services.FindUsers(filterDto)
if err != nil {
    return err
}

// Transform list of models → list of responses
data := transformers.ToListResponse(users, transformers.ToUserResponse)

return c.Success(response.ListUser{
    Meta: response.Meta{Page: 1, PerPage: 10, Total: total},
    Data: data,
})
```

### Null Handling Helpers

The `github.com/gflydev/db/null` package provides conversion helpers:

```go
import dbNull "github.com/gflydev/db/null"

// sql.NullString → *string
func StringNil(ns sql.NullString) *string {
    if !ns.Valid {
        return nil
    }
    return &ns.String
}

// sql.NullTime → *time.Time
func TimeNil(nt sql.NullTime) *time.Time {
    if !nt.Valid {
        return nil
    }
    return &nt.Time
}
```

**Usage:**

```go
type User struct {
    Token sql.NullString  // Domain model uses sql.NullString
    Avatar sql.NullString
}

// In transformer
response.User{
    Token: dbNull.StringNil(user.Token),   // → *string
    Avatar: PublicAvatar(user.Avatar.String),  // Custom handling
}
```

### Module Transformer Example

```go
// File: pkg/modules/auth/transformers/auth_transformers.go
package transformers

import (
    "gfly/pkg/modules/auth"
    "gfly/pkg/modules/auth/response"
)

// ToSignInResponse converts JWT tokens to SignIn response.
//
// Parameters:
//   - tokens: JWT token pair (access + refresh)
//
// Returns:
//   - response.SignIn: SignIn response object
func ToSignInResponse(tokens *auth.Token) response.SignIn {
    return response.SignIn{
        Access:  tokens.Access,
        Refresh: tokens.Refresh,
    }
}
```

---

## Validation System

### Validation Flow

```
1. HTTP Request arrives
   ↓
2. Controller.Validate(c) invoked by gFly framework
   ↓
3. ProcessData[RequestType](c) or ProcessUpdateData[RequestType](c)
   ├─ Parse JSON body → Request struct
   ├─ Sanitize input (XSS prevention)
   └─ Validate using DTO's validate tags
   ↓
4. If validation passes:
   ├─ Store validated request in context: c.SetData(constants.Request, request)
   └─ Proceed to Controller.Handle(c)
   ↓
5. If validation fails:
   └─ Return error response with field-level errors
```

### Validation Implementation

```go
// From: github.com/gflydev/http

// Validate performs input validation using gflydev/validation.
//
// Parameters:
//   - structData: Any struct with validate tags
//   - msgForTagFunc: Optional custom error messages
//
// Returns:
//   - *response.Error: Validation error with field details, or nil if valid
func Validate(structData any, msgForTagFunc ...validation.MsgForTagFunc) *response.Error {
    errorData, err := validation.Check(structData, msgForTagFunc...)

    if err != nil {
        return &response.Error{
            Code:    core.StatusBadRequest,
            Message: "Invalid input",
            Data:    errorData,  // Field-level validation errors
        }
    }

    return nil
}
```

### Request Processing Helpers

```go
// From: github.com/gflydev/http

// ProcessData validates and processes create/add requests.
//
// Type Parameters:
//   - T: Request type implementing AddData interface
//
// Returns:
//   - error: Validation error or nil if successful
func ProcessData[T AddData](c *core.Ctx) error {
    var requestData T

    // 1. Parse JSON body
    if errData := Parse(c, &requestData); errData != nil {
        return c.Error(errData)
    }

    // 2. Sanitize input (XSS, HTML injection prevention)
    security.SanitizeStruct(&requestData)

    // 3. Validate using DTO rules
    if errData := Validate(requestData); errData != nil {
        return c.Error(errData)
    }

    // 4. Store validated request in context
    c.SetData(constants.Request, requestData)

    return nil
}

// ProcessUpdateData validates update requests (includes path ID).
//
// Type Parameters:
//   - T: Request type implementing UpdateData interface
//
// Returns:
//   - error: Validation error or nil if successful
func ProcessUpdateData[T UpdateData](c *core.Ctx) error {
    // 1. Extract ID from URL path
    itemID, errData := PathID(c)
    if errData != nil {
        return c.Error(errData)
    }

    // 2. Parse JSON body
    var requestData T
    if errData := Parse(c, &requestData); errData != nil {
        return c.Error(errData)
    }

    // 3. Sanitize
    security.SanitizeStruct(&requestData)

    // 4. Set ID from path parameter
    requestData.SetID(itemID)

    // 5. Validate
    if errData := Validate(requestData); errData != nil {
        return c.Error(errData)
    }

    // 6. Store in context
    c.SetData(constants.Request, requestData)

    return nil
}

// ProcessFilter validates list/filter requests from query parameters.
func ProcessFilter(c *core.Ctx) error {
    var filterDto dto.Filter

    // Parse query parameters
    filterDto.Keyword = c.Query("keyword")
    filterDto.OrderBy = c.Query("order_by")
    filterDto.Page = c.QueryInt("page", 1)
    filterDto.PerPage = c.QueryInt("per_page", 10)

    // Sanitize
    security.SanitizeStruct(&filterDto)

    // Validate
    if errData := Validate(filterDto); errData != nil {
        return c.Error(errData)
    }

    // Store in context
    c.SetData(constants.Filter, filterDto)

    return nil
}
```

### Validation Error Response

When validation fails, the response contains detailed field errors:

```json
{
  "code": 400,
  "message": "Invalid input",
  "data": {
    "email": ["email must be a valid email address"],
    "password": ["password must be at least 8 characters"],
    "status": ["status must be one of: active pending blocked"]
  }
}
```

### Custom Validation Messages

```go
// Define custom error messages for specific tags
func customMessages() validation.MsgForTagFunc {
    return func(tag string) string {
        switch tag {
        case "required":
            return "This field is required"
        case "email":
            return "Please provide a valid email address"
        case "min":
            return "Value is too short"
        case "max":
            return "Value is too long"
        default:
            return ""
        }
    }
}

// Use in validation
if errData := Validate(requestData, customMessages()); errData != nil {
    return c.Error(errData)
}
```

---

## Complete Flow Examples

### Example 1: Create User (POST Request)

#### 1. Define DTO

```go
// File: internal/dto/user_dto.go
package dto

type CreateUser struct {
    Email    string `json:"email" validate:"required,email,max=255" doc:"User's email"`
    Password string `json:"password" validate:"required,min=8,max=255" doc:"User's password"`
    Fullname string `json:"fullname" validate:"required,max=255" doc:"User's full name"`
}
```

#### 2. Define Request

```go
// File: internal/http/request/user_request.go
package request

import "gfly/internal/dto"

type CreateUser struct {
    dto.CreateUser
}

func (r CreateUser) ToDto() dto.CreateUser {
    return r.CreateUser
}
```

#### 3. Define Response

```go
// File: internal/http/response/user_response.go
package response

import "time"

type User struct {
    ID        int       `json:"id" doc:"User ID"`
    Email     string    `json:"email" doc:"User's email"`
    Fullname  string    `json:"fullname" doc:"User's name"`
    CreatedAt time.Time `json:"created_at" doc:"Creation timestamp"`
}
```

#### 4. Define Transformer

```go
// File: internal/http/transformers/user_transformer.go
package transformers

import (
    "gfly/internal/domain/models"
    "gfly/internal/http/response"
)

func ToUserResponse(user models.User) response.User {
    return response.User{
        ID:        user.ID,
        Email:     user.Email,
        Fullname:  user.Fullname,
        CreatedAt: user.CreatedAt,
    }
}
```

#### 5. Implement Service

```go
// File: internal/services/user_services.go
package services

import (
    "gfly/internal/domain/models"
    "gfly/internal/dto"
    mb "github.com/gflydev/db"
    "github.com/gflydev/core/errors"
    "time"
)

// CreateUser creates a new user.
//
// Parameters:
//   - createUserDto: User creation data
//
// Returns:
//   - (*models.User, error): Created user or error
func CreateUser(createUserDto dto.CreateUser) (*models.User, error) {
    // Business logic: Check email uniqueness
    existing := repository.Pool.GetUserByEmail(createUserDto.Email)
    if existing != nil {
        return nil, errors.New("user with this email already exists")
    }

    // Create domain model
    user := &models.User{
        Email:     createUserDto.Email,
        Password:  utils.GeneratePassword(createUserDto.Password),
        Fullname:  createUserDto.Fullname,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }

    // Persist to database
    err := mb.CreateModel(user)
    if err != nil {
        return nil, errors.New("failed to create user: %w", err)
    }

    return user, nil
}
```

#### 6. Implement Controller

```go
// File: internal/http/controllers/api/user/create_user_api.go
package user

import (
    "gfly/internal/http/request"
    "gfly/internal/http/response"
    "gfly/internal/http/transformers"
    "gfly/internal/services"
    "github.com/gflydev/http"
    "github.com/gflydev/core"
)

type CreateUserApi struct {
    core.Api
}

func NewCreateUserApi() *CreateUserApi {
    return &CreateUserApi{}
}

// Validate handles request validation.
func (h *CreateUserApi) Validate(c *core.Ctx) error {
    return http.ProcessData[request.CreateUser](c)
}

// Handle processes the create user request.
// @Summary Create a new user
// @Description Creates a new user account
// @Tags Users
// @Accept json
// @Produce json
// @Param data body request.CreateUser true "User creation data"
// @Success 201 {object} response.User
// @Failure 400 {object} response.Error
// @Security ApiKeyAuth
// @Router /users [post]
func (h *CreateUserApi) Handle(c *core.Ctx) error {
    // 1. Get validated request from context
    requestData := c.GetData(constants.Request).(request.CreateUser)

    // 2. Convert to DTO
    createUserDto := requestData.ToDto()

    // 3. Call service
    user, err := services.CreateUser(createUserDto)
    if err != nil {
        return c.Error(response.Error{
            Code:    core.StatusBadRequest,
            Message: err.Error(),
        })
    }

    // 4. Transform to response
    userResponse := transformers.ToUserResponse(*user)

    // 5. Return JSON response
    return c.Status(core.StatusCreated).JSON(userResponse)
}
```

#### 7. Register Route

```go
// File: internal/http/routes/api_routes.go

usersGroup := apiRouter.Group("/users")
usersGroup.POST("", app.Apply(
    middleware.JWTAuth(),
    middleware.CheckRolesMiddleware(types.RoleAdmin),
)(user.NewCreateUserApi()))
```

#### Complete Flow:

```
POST /api/v1/users
Content-Type: application/json
Authorization: Bearer <token>

{
  "email": "john@example.com",
  "password": "SecurePass123",
  "fullname": "John Doe"
}

         ↓

Controller.Validate(c)
  → ProcessData[request.CreateUser](c)
    → Parse JSON
    → Sanitize
    → Validate (email format, min length, etc.)
    → Store in context

         ↓

Controller.Handle(c)
  → Get request from context
  → Convert to DTO: requestData.ToDto()
  → Call service: services.CreateUser(dto)

         ↓

Service Layer
  → Validate business rules (email uniqueness)
  → Create domain model
  → Persist to database
  → Return models.User

         ↓

Controller
  → Transform model: transformers.ToUserResponse(user)
  → Return JSON response

         ↓

HTTP 201 Created
{
  "id": 123,
  "email": "john@example.com",
  "fullname": "John Doe",
  "created_at": "2025-01-15T10:30:00Z"
}
```

---

### Example 2: Update User (PUT Request)

#### Request Flow with Path ID:

```
PUT /api/v1/users/123
Content-Type: application/json
Authorization: Bearer <token>

{
  "fullname": "Jane Doe",
  "password": "NewPassword456"
}

         ↓

Controller.Validate(c)
  → ProcessUpdateData[request.UpdateUser](c)
    → Extract ID from URL path (123)
    → Parse JSON body
    → Sanitize
    → Call requestData.SetID(123)  // Sets DTO.ID = 123
    → Validate
    → Store in context

         ↓

Controller.Handle(c)
  → Get request from context
  → Convert to DTO: requestData.ToDto()  // DTO has ID=123
  → Call service: services.UpdateUser(dto)

         ↓

Service Layer
  → Find user by dto.ID (123)
  → Apply updates from DTO
  → Persist changes
  → Return updated models.User

         ↓

Controller
  → Transform: transformers.ToUserResponse(user)
  → Return JSON

         ↓

HTTP 200 OK
{
  "id": 123,
  "email": "john@example.com",
  "fullname": "Jane Doe",
  "updated_at": "2025-01-15T11:00:00Z"
}
```

**Important:** `SetID()` must use a **pointer receiver**:

```go
// CORRECT - pointer receiver
func (r *UpdateUser) SetID(id int) {
    r.ID = id
}

// INCORRECT - value receiver (won't modify the struct)
func (r UpdateUser) SetID(id int) {
    r.ID = id  // This creates a copy, original struct unchanged!
}
```

---

### Example 3: List Users (GET Request)

#### Query Parameter Flow:

```
GET /api/v1/users?keyword=john&page=2&per_page=20&order_by=created_at
Authorization: Bearer <token>

         ↓

Controller.Validate(c)
  → ProcessFilter(c)
    → Parse query params: keyword, page, per_page, order_by
    → Create dto.Filter
    → Sanitize
    → Validate (page >= 1, per_page <= 100, etc.)
    → Store in context

         ↓

Controller.Handle(c)
  → Get filter from context
  → Call service: services.FindUsers(filterDto)
  → Service returns: ([]models.User, totalCount, error)

         ↓

Controller
  → Create metadata: response.Meta{Page: 2, PerPage: 20, Total: 154}
  → Transform list: transformers.ToListResponse(users, transformers.ToUserResponse)
  → Return response.ListUser{Meta, Data}

         ↓

HTTP 200 OK
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

---

## Best Practices

### 1. DTO Design

✅ **DO:**
- Use DTOs as the single source of truth for validation
- Include comprehensive `doc` tags for documentation
- Use specific validation rules (`oneof`, `min`, `max`)
- Group related fields logically
- Use semantic naming (CreateX, UpdateX, FilterX)

❌ **DON'T:**
- Put business logic in DTOs (keep them pure data)
- Reuse input DTOs as output (use separate Response types)
- Skip validation tags ("be explicit, not implicit")
- Use generic names (`Data`, `Payload`)

### 2. Request Objects

✅ **DO:**
- **ALWAYS implement `ToDto()` method** - this is REQUIRED by the interface
- Always embed the corresponding DTO(s)
- For single DTO: `ToDto()` returns that DTO
- For multiple DTOs: `ToDto()` returns the primary/default DTO + implement `ToDtoXXX()` methods for alternatives
- Use pointer receiver for `SetID()` in update requests
- Keep requests thin (minimal logic)
- Use descriptive names for multiple conversion methods (`ToDtoAdmin()`, `ToDtoUser()`, `ToDtoInfo()`)
- Make `ToDto()` delegate to a specific `ToDtoXXX()` method when using multiple DTOs

❌ **DON'T:**
- Skip implementing `ToDto()` - it's required by the interface!
- Add business logic to requests
- Duplicate DTO fields (use embedding)
- Think `ToDtoXXX()` methods replace `ToDto()` - they supplement it
- Mix concerns in a single DTO (create separate DTOs for different contexts)

### 3. Validation

✅ **DO:**
- Validate at the HTTP boundary (in controller's `Validate()`)
- Use type-safe generics (`ProcessData[T]`, `ProcessUpdateData[T]`)
- Sanitize input before validation
- Return detailed field-level errors

❌ **DON'T:**
- Skip validation ("trust user input")
- Validate in service layer (that's for business rules)
- Swallow validation errors

### 4. Service Layer

✅ **DO:**
- Accept DTOs as parameters
- Return domain models (not DTOs or responses)
- Validate business rules (not field formats)
- Use repository pattern for data access
- Return errors using `github.com/gflydev/core/errors`

❌ **DON'T:**
- Accept HTTP context or request objects
- Return Response objects (that's the transformer's job)
- Access database directly (use repositories)
- Use `fmt.Errorf` (use gFly's error package)

### 5. Transformers

✅ **DO:**
- Handle null conversion (`sql.NullX` → `*Type`)
- Format URLs, paths, dates consistently
- Enrich with related data when needed
- Name transformers clearly (`ToUserResponse`, `ToListResponse`)
- Use generic transformers for lists

❌ **DON'T:**
- Put business logic in transformers
- Query database for every field (batch load related data)
- Return different response structures for same entity

### 6. Response Objects

✅ **DO:**
- Use pointers for nullable fields
- Include comprehensive `doc` tags
- Match JSON key naming convention (snake_case or camelCase consistently)
- Version responses when making breaking changes

❌ **DON'T:**
- Expose internal IDs or sensitive data unnecessarily
- Return raw domain models as JSON
- Mix required and optional fields without clear indication

### 7. Controller Pattern

✅ **DO:**
- Implement both `Validate()` and `Handle()` methods
- Use appropriate `ProcessData`, `ProcessUpdateData`, or `ProcessFilter`
- Include complete Swagger annotations
- Return structured errors via `c.Error(response.Error{...})`

❌ **DON'T:**
- Skip validation step
- Put business logic in controllers
- Return inconsistent error formats

### 8. Error Handling

✅ **DO:**
- Use `response.Error` for validation/input errors
- Use `response.Unauthorized` for auth errors
- Use `response.NotFound` for missing resources
- Include helpful error messages

❌ **DON'T:**
- Expose stack traces or internal errors to clients
- Return generic "error occurred" messages
- Mix error response formats

### 9. Testing

✅ **DO:**
- Test DTOs with valid and invalid data
- Test services with DTOs (no HTTP context needed)
- Test transformers with various null/non-null combinations
- Test validation rules comprehensively

❌ **DON'T:**
- Skip testing edge cases (null, empty, max length)
- Test only happy paths

### 10. Documentation

✅ **DO:**
- Document all public functions with GoDoc
- Include Parameters and Returns sections
- Use `doc` tags on all struct fields
- Complete Swagger annotations on controllers
- Keep CLAUDE.md updated with patterns

❌ **DON'T:**
- Skip documentation ("code is self-documenting")
- Use vague descriptions
- Forget to run `make doc` after API changes

---

## Summary

This guide provides a comprehensive reference for working with Request, DTO, Response, and Transformer objects in the gFly framework. The key takeaways:

1. **DTOs** are the single source of truth for data structure and validation
2. **Requests** are thin HTTP wrappers that convert to DTOs
   - ⚠️ **CRITICAL**: Every Request **MUST** implement `ToDto()` method (required by interface)
   - For multiple DTOs: `ToDto()` returns primary DTO + add `ToDtoXXX()` methods for alternatives
3. **Services** operate on DTOs and return domain models
4. **Transformers** convert domain models to response objects
5. **Responses** define the JSON output structure
6. **Validation** happens at the HTTP boundary using DTO rules
7. **Separation of concerns** is maintained throughout the stack

By following these patterns and best practices, you'll create maintainable, testable, and well-documented APIs that align with Clean Architecture principles and gFly framework conventions.

---

## Quick Reference

### File Locations

| Component | Central Location | Module Location |
|-----------|-----------------|-----------------|
| DTO | `/internal/dto/` | `/pkg/modules/{module}/dto/` |
| Request | `/internal/http/request/` | `/pkg/modules/{module}/request/` |
| Response | `/internal/http/response/` | `/pkg/modules/{module}/response/` |
| Transformer | `/internal/http/transformers/` | `/pkg/modules/{module}/transformers/` |
| Controller | `/internal/http/controllers/api/` | `/pkg/modules/{module}/api/` |
| Service | `/internal/services/` | `/pkg/modules/{module}/services/` |

### Common Commands

```bash
# Generate Swagger docs after API changes
make doc

# Run tests
make test

# Run linter
make lint

# Check code quality
make check
```

### Validation Tags Quick Reference

```go
validate:"required"              // Field must be present
validate:"omitempty"             // Field is optional
validate:"email"                 // Valid email format
validate:"min=8,max=255"         // Length constraints
validate:"gte=1,lte=100"         // Numeric constraints
validate:"oneof=active pending"  // Enum validation
validate:"url"                   // Valid URL
validate:"uuid"                  // Valid UUID
```

---

**Last Updated:** 2025-01-22
**gFly Framework Version:** Compatible with github.com/gflydev/core v1.x
