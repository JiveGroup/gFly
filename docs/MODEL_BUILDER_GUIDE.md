# gFly Framework: Model Builder (mb) Usage Guide

This comprehensive guide documents all Model Builder patterns, rules, and best practices used in the service layer of the ThietNgon e-commerce platform.

## Table of Contents

1. [Setup & Initialization](#1-setup--initialization)
2. [Basic CRUD Operations](#2-basic-crud-operations)
3. [Query Builder Patterns](#3-query-builder-patterns)
4. [Advanced Filtering](#4-advanced-filtering)
5. [Pagination & Ordering](#5-pagination--ordering)
6. [Joins & Relationships](#6-joins--relationships)
7. [Transactions](#7-transactions)
8. [Raw SQL Queries](#8-raw-sql-queries)
9. [Error Handling Patterns](#9-error-handling-patterns)
10. [Best Practices & Rules](#10-best-practices--rules)

---

## 1. Setup & Initialization

### Import Pattern

```go
import (
    mb "github.com/gflydev/db"
    dbNull "github.com/gflydev/db/null"
)
```

### Database Driver Registration

```go
// In cmd/web/main.go
mb.Register(dbPSQL.New())  // Register PostgreSQL driver
mb.Load()                   // Load the model builder
```

### Test Environment Setup

```go
// In test/setup_test.go
mb.Register(dbPSQL.New())
mb.Load()
```

---

## 2. Basic CRUD Operations

### 2.1 CREATE Operations

#### Generic CreateModel
Used for creating new records in the database.

**Pattern:**
```go
model := &models.Model{
    Field1: value1,
    Field2: value2,
    CreatedAt: time.Now(),
}

err := mb.CreateModel(model)
```

**Real Examples:**

**Example 1: Create Product**
```go
product := &models.Product{
    BrandID:       dbNull.Int32(int32(brand.ID)),
    VarietyID:     dbNull.Int32(int32(variety.ID)),
    Name:          createDTO.Name,
    Slug:          fmt.Sprintf("%s-%d", strings.ToLower(strings.ReplaceAll(createDTO.Name, " ", "-")), time.Now().Unix()),
    Description:   createDTO.Description,
    Content:       createDTO.Content,
    Price:         createDTO.Price,
    Currency:      currency,
    CreatedAt:     time.Now(),
}

err = mb.CreateModel(product)
```

**Example 2: Create Coupon**
```go
coupon := &models.Coupon{
    Code:        createDto.Code,
    Description: dbNull.String(createDto.Description),
    Type:        createDto.Type,
    Quantity:    createDto.Quantity,
    CreatedAt:   time.Now(),
}

err := mb.CreateModel(coupon)
```

**Example 3: Create Cart**
```go
cart := &models.Cart{
    UserID:    userID,
    Name:      name,
    SessionID: str.Random(10),
    CreatedAt: time.Now(),
}

err = mb.CreateModel(cart)
```

#### Transaction-Based Create
Used within transactions for atomic operations.

**Example: Create Order with Items**
```go
db := mb.Instance()
db.Begin()

// Create address
if err := db.Create(&address); err != nil {
    log.Errorf("Failed to create address: %v", err)
    try.Throw(errors.New("Failed to create address"))
}

// Create order
if err := db.Create(&order); err != nil {
    log.Errorf("Failed to create order: %v", err)
    try.Throw(errors.New("Failed to create order"))
}

// Create order items in a loop
for _, item := range cartItems {
    orderItem := models.OrderItem{
        OrderID:   order.ID,
        ProductID: item.ProductID,
        // ... other fields
    }
    if err := db.Create(&orderItem); err != nil {
        log.Errorf("Failed to create order item: %v", err)
        try.Throw(errors.New("Failed to create order item"))
    }
}
```

**Rules:**
- Always set `CreatedAt` to `time.Now()` for new records
- Use `dbNull` package for nullable fields
- Return the created model pointer for further operations
- Log errors before returning

---

### 2.2 READ Operations

#### Generic GetModelByID
Retrieve a single record by its primary key (ID).

**Pattern:**
```go
model, err := mb.GetModelByID[models.ModelName](id)
```

**Real Examples:**

**Example 1: Get Product by ID**
```go
product, err := mb.GetModelByID[models.Product](updateProductDto.ID)
if err != nil {
    log.Error(err)
    return nil, errors.ItemNotFound
}
```

**Example 2: Get Brand by ID**
```go
brand, err := mb.GetModelByID[models.Brand](createDTO.BrandID)
if err != nil {
    log.Error(err)
    return nil, errors.ItemNotFound
}
```

**Example 3: Get User by ID**
```go
user, err := mb.GetModelByID[models.User](updateUserDto.ID)
if err != nil {
    return nil, errors.New("User not found")
}
```

#### Generic GetModelBy
Retrieve a single record by any field.

**Pattern:**
```go
model, err := mb.GetModelBy[models.ModelName]("field_name", value)
```

**Real Examples:**

**Example 1: Get Coupon by Code**
```go
coupon, err := mb.GetModelBy[models.Coupon]("code", code)
if err != nil {
    return nil, errors.ItemNotFound
}
```

**Example 2: Get User by Email**
```go
existingUser, err := mb.GetModelBy[models.User]("email", updateProfileDto.Email)
if err == nil && existingUser != nil && existingUser.ID != user.ID {
    return nil, errors.New("Email already exists for another user")
}
```

**Example 3: Repository Helper Pattern**
```go
func (r *userRepository) getBy(field string, value any) *models.User {
    user, err := mb.GetModelBy[models.User](field, value)
    if err != nil {
        log.Error(err)
        return nil
    }
    return user
}

func (r *userRepository) FindByEmail(email string) *models.User {
    return r.getBy("email", email)
}
```

#### Query Builder - First()
Get a single record with custom conditions.

**Pattern:**
```go
var model models.ModelName
err := mb.Instance().
    Where("field", mb.Eq, value).
    First(&model)
```

**Real Examples:**

**Example 1: Get Active Cart**
```go
var cart models.Cart
err = mb.Instance().Where("user_id", mb.Eq, userID).
    OrderBy("created_at", mb.Desc).
    First(&cart)
```

**Example 2: Check Product in Cart**
```go
var cartItem models.CartItem
err = mb.Instance().
    Where("cart_id", mb.Eq, cart.ID).
    Where("product_id", mb.Eq, dtoData.ProductID).
    First(&cartItem)
```

**Example 3: Verify Address Ownership**
```go
var address models.Address
err := db.
    Where("id", mb.Eq, dtoData.AddressID).
    Where("user_id", mb.Eq, user.ID).
    First(&address)
```

#### Query Builder - Find()
Get multiple records with custom conditions.

**Pattern:**
```go
var models []models.ModelName
total, err := mb.Instance().
    Where("field", mb.Eq, value).
    Find(&models)
```

**Real Examples:**

**Example 1: Get Cart Items**
```go
var cartItems []models.CartItem
total, err := mb.Instance().Where("cart_id", mb.Eq, cartID).
    OrderBy("created_at", mb.Desc).
    Find(&cartItems)
```

**Example 2: Get Files by Content ID**
```go
var list []models.ContentFile
_, err := mb.Instance().
    Select(fmt.Sprintf("%s.*", models.TableContentFile)).
    Where("content_id", mb.Eq, contentID).
    Find(&list)
```

**Rules:**
- Use `GetModelByID` for simple ID lookups
- Use `GetModelBy` for lookups by any single field
- Use `Instance().Where().First()` for complex single-record queries
- Use `Instance().Where().Find()` for multiple records
- Always check error and handle `ItemNotFound` case
- `Find()` returns total count as first return value

---

### 2.3 UPDATE Operations

#### Generic UpdateModel
Update an existing record in the database.

**Pattern:**
```go
// 1. Fetch the model
model, err := mb.GetModelByID[models.ModelName](id)
if err != nil {
    return nil, errors.ItemNotFound
}

// 2. Update fields
model.Field1 = newValue1
model.Field2 = newValue2
model.UpdatedAt = dbNull.TimeNow()

// 3. Persist changes
err = mb.UpdateModel(model)
```

**Real Examples:**

**Example 1: Update Product**
```go
product, err := mb.GetModelByID[models.Product](updateProductDto.ID)
if err != nil {
    log.Error(err)
    return nil, errors.ItemNotFound
}

// Update fields if provided and different
if updateProductDto.Name != "" && updateProductDto.Name != product.Name {
    product.Name = updateProductDto.Name
}

if updateProductDto.Description != "" && updateProductDto.Description != product.Description {
    product.Description = updateProductDto.Description
}

product.UpdatedAt = dbNull.TimeNow()

err = mb.UpdateModel(product)
```

**Example 2: Update Coupon**
```go
coupon, err := mb.GetModelByID[models.Coupon](updateDto.ID)
if err != nil {
    return nil, errors.ItemNotFound
}

// Partially update
if updateDto.Quantity != 0 {
    coupon.Quantity = updateDto.Quantity
}

if updateDto.Description != "" {
    coupon.Description = dbNull.String(updateDto.Description)
}

coupon.UpdatedAt = dbNull.TimeNow()

err = mb.UpdateModel(coupon)
```

**Example 3: Update Cart Item Quantity**
```go
cartItem, err := mb.GetModelByID[models.CartItem](dtoData.CartItemID)
if err != nil {
    return err
}

cartItem.Quantity = dtoData.Quantity

return mb.UpdateModel(cartItem)
```

**Example 4: Update User Profile**
```go
// Update user fields
if updateProfileDto.FirstName != "" {
    user.FirstName = dbNull.String(updateProfileDto.FirstName)
}
if updateProfileDto.LastName != "" {
    user.LastName = dbNull.String(updateProfileDto.LastName)
}

user.UpdatedAt = dbNull.TimeNow()

err := mb.UpdateModel(user)
```

#### Transaction-Based Update

**Example: Update Order Total**
```go
db := mb.Instance()
db.Begin()

order.TotalPrice = totalPrice
if err := db.Update(order); err != nil {
    log.Errorf("Failed to update order total price: %v", err)
    try.Throw(errors.New("Failed to update order total price"))
}
```

**Rules:**
- Always fetch the model first before updating
- Update only changed fields (compare old vs new values)
- Always update `UpdatedAt` timestamp with `dbNull.TimeNow()`
- Use `dbNull` package for nullable fields
- Return the updated model for further operations
- Log errors before returning

---

### 2.4 DELETE Operations

#### Generic DeleteModel
Delete a record using the model instance.

**Pattern:**
```go
// 1. Fetch the model
model, err := mb.GetModelByID[models.ModelName](id)
if err != nil {
    return errors.ItemNotFound
}

// 2. Delete it
err = mb.DeleteModel(model)
```

**Real Examples:**

**Example 1: Delete Product**
```go
product, err := mb.GetModelByID[models.Product](productID)
if err != nil {
    log.Error(err)
    return errors.ItemNotFound
}

return mb.DeleteModel(product)
```

**Example 2: Delete Coupon**
```go
coupon, err := mb.GetModelByID[models.Coupon](couponID)
if err != nil {
    return errors.ItemNotFound
}

return mb.DeleteModel(coupon)
```

**Example 3: Delete User**
```go
user, err := mb.GetModelByID[models.User](userID)
if err != nil {
    return errors.New("User not found")
}

// Delete roles that sync with user
if err := repository.Pool.SyncRolesWithUser(userID, ""); err != nil {
    log.Errorf("Error while deleting user roles: %v", err)
    return errors.New("error occurs while deleting user roles")
}

// Delete user
if err := mb.DeleteModel(user); err != nil {
    log.Errorf("Error while deleting user: %v", err)
    return errors.New("error occurs while deleting user")
}
```

#### Direct Delete with Where Clause
Delete without fetching the model first.

**Pattern:**
```go
err := mb.Instance().Where("field", mb.Eq, value).Delete(models.ModelName{})
```

**Real Examples:**

**Example 1: Delete Cart**
```go
if err := db.Where("id", mb.Eq, cart.ID).Delete(models.Cart{}); err != nil {
    log.Errorf("Failed to delete cart: %v", err)
    try.Throw(errors.New("Failed to delete cart"))
}
```

**Example 2: Delete Cart Item**
```go
err = mb.Instance().Where("id", mb.Eq, cartItemID).
    Delete(&models.CartItem{})
```

**Rules:**
- Use `DeleteModel()` when you already have the model instance
- Use `Where().Delete()` for direct deletion without fetching
- Always validate existence before deletion (unless using Where clause)
- Delete related records first (foreign key constraints)
- Return appropriate error if record not found
- Consider soft deletes for audit trails

---

## 3. Query Builder Patterns

### 3.1 Instance() Method

The `mb.Instance()` method returns a query builder for constructing complex queries.

**Pattern:**
```go
dbInstance := mb.Instance()
```

**Usage Context:**
- Complex filtering
- Joins
- Custom ordering
- Pagination
- Aggregation

### 3.2 Select() Method

Specify which columns to retrieve.

**Pattern:**
```go
// Select all columns from a table
mb.Instance().Select(fmt.Sprintf("%s.*", models.TableName))

// Select specific columns
mb.Instance().Select("id", "name", "email")

// Select with aliases and joins
mb.Instance().Select(
    fmt.Sprintf("%s.*", models.TableProduct),
    fmt.Sprintf("%s.name as category_name", models.TableCategory),
)
```

**Real Examples:**

**Example 1: Select All Columns**
```go
builder := dbInstance.Select(fmt.Sprintf("%s.*", models.TableProduct))
```

**Example 2: Select with Join**
```go
builder := dbInstance.Select("DISTINCT users.id", "users.*").
    Join(mb.InnerJoin, "user_roles", mb.Condition{
        Field: "user_roles.user_id",
        Opt:   mb.Eq,
        Value: mb.ValueField("users.id"),
    })
```

**Example 3: Select Multiple Tables**
```go
builder := mb.Instance().
    Select(
        fmt.Sprintf("%s.*", models.TableRating),
        fmt.Sprintf("%s.email", models.TableUser),
        fmt.Sprintf("%s.fullname", models.TableUser),
        fmt.Sprintf("%s.avatar", models.TableUser),
    )
```

### 3.3 Where() Clauses

#### Simple Where

**Pattern:**
```go
mb.Instance().Where("field", operator, value)
```

**Available Operators:**
- `mb.Eq` - Equals (=)
- `mb.NotEq` - Not equals (!=)
- `mb.Like` - LIKE pattern matching
- `mb.In` - IN clause
- `mb.Null` - IS NULL check

**Real Examples:**

**Example 1: Simple Equality**
```go
err = mb.Instance().Where("user_id", mb.Eq, userID).First(&cart)
```

**Example 2: LIKE Operator**
```go
query.Where("name", mb.Like, "%"+filterDto.Keyword+"%")
```

**Example 3: IN Operator**
```go
builder.Where(models.TableProductCategory+".category_id", mb.In, categoryIDs)
```

**Example 4: NOT EQUAL**
```go
builder.Where(models.TableProduct+".id", mb.NotEq, productID)
```

**Example 5: NULL Check**
```go
Where("users.deleted_at", mb.Null, nil)
```

#### Multiple Where Clauses (AND Logic)

**Pattern:**
```go
mb.Instance().
    Where("field1", mb.Eq, value1).
    Where("field2", mb.Eq, value2).
    Where("field3", mb.Eq, value3)
```

**Real Examples:**

**Example 1: Verify Cart Item Ownership**
```go
var cartItem models.CartItem
err = mb.Instance().
    Where("cart_id", mb.Eq, cart.ID).
    Where("product_id", mb.Eq, dtoData.ProductID).
    First(&cartItem)
```

**Example 2: Check Duplicate Rating Item**
```go
err := mb.Instance().
    Where("item_type", mb.Eq, createDTO.ItemType).
    Where("item_ref", mb.Eq, createDTO.ItemRef).
    First(&existingItem)
```

**Example 3: Get Rating Stats**
```go
err := mb.Instance().
    Where("item_id", mb.Eq, itemID).
    Where("criterion_id", mb.Eq, criterionID).
    First(&stat)
```

### 3.4 Comparison Operators Reference

```go
// Equality
mb.Eq       // field = value
mb.NotEq    // field != value

// Pattern Matching
mb.Like     // field LIKE value

// Membership
mb.In       // field IN (values)

// Nullability
mb.Null     // field IS NULL

// Ordering
mb.Asc      // ORDER BY field ASC
mb.Desc     // ORDER BY field DESC

// Join Types
mb.InnerJoin  // INNER JOIN
```

### 3.5 ValueField() for Field References

Used in joins to reference another column instead of a literal value.

**Pattern:**
```go
mb.ValueField("table.column")
```

**Real Examples:**

**Example 1: Join Condition**
```go
Join(mb.InnerJoin, models.TableProductCategory,
    mb.Condition{
        Field: models.TableProductCategory + ".product_id",
        Opt:   mb.Eq,
        Value: mb.ValueField(models.TableProduct + ".id"),
    })
```

**Example 2: User Roles Join**
```go
Join(mb.InnerJoin, "user_roles", mb.Condition{
    Field: "user_roles.user_id",
    Opt:   mb.Eq,
    Value: mb.ValueField("users.id"),
}).
Join(mb.InnerJoin, "roles", mb.Condition{
    Field: "roles.id",
    Opt:   mb.Eq,
    Value: mb.ValueField("user_roles.role_id"),
})
```

---

## 4. Advanced Filtering

### 4.1 WhereGroup() for OR Logic

Use `WhereGroup()` to group multiple conditions with OR logic.

**Pattern:**
```go
builder.WhereGroup(func(queryGroup mb.WhereBuilder) *mb.WhereBuilder {
    queryGroup.
        Where("field1", mb.Like, "%value%").
        WhereOr("field2", mb.Like, "%value%").
        WhereOr("field3", mb.Like, "%value%")
    return &queryGroup
})
```

**Real Examples:**

**Example 1: Product Keyword Search**
```go
builder := dbInstance.Select(fmt.Sprintf("%s.*", models.TableProduct)).
    When(filterDto.Keyword != "", func(query mb.WhereBuilder) *mb.WhereBuilder {
        query.WhereGroup(func(queryGroup mb.WhereBuilder) *mb.WhereBuilder {
            queryGroup.
                Where("name", mb.Like, "%"+filterDto.Keyword+"%").
                WhereOr("description", mb.Like, "%"+filterDto.Keyword+"%").
                WhereOr("content", mb.Like, "%"+filterDto.Keyword+"%")

            return &queryGroup
        })
        return &query
    })
```

**Example 2: User Search**
```go
When(filterDto.Keyword != "", func(query mb.WhereBuilder) *mb.WhereBuilder {
    query.WhereGroup(func(queryGroup mb.WhereBuilder) *mb.WhereBuilder {
        queryGroup.Where("roles.name", mb.Like, "%"+filterDto.Keyword+"%").
            WhereOr("roles.slug", mb.Like, "%"+filterDto.Keyword+"%").
            WhereOr("users.email", mb.Like, "%"+filterDto.Keyword+"%").
            WhereOr("users.first_name", mb.Like, "%"+filterDto.Keyword+"%").
            WhereOr("users.last_name", mb.Like, "%"+filterDto.Keyword+"%").
            WhereOr("users.phone", mb.Like, "%"+filterDto.Keyword+"%")

        if slices.Contains(types.UserStatusList, types.UserStatus(filterDto.Keyword)) {
            queryGroup.WhereOr("users.status", mb.Eq, filterDto.Keyword)
        }

        return &queryGroup
    })

    return &query
})
```

**Example 3: Contact Submission Search**
```go
When(filterDto.Keyword != "", func(query mb.WhereBuilder) *mb.WhereBuilder {
    query.WhereGroup(func(queryGroup mb.WhereBuilder) *mb.WhereBuilder {
        queryGroup.
            Where("name", mb.Like, "%"+filterDto.Keyword+"%").
            WhereOr("email", mb.Like, "%"+filterDto.Keyword+"%").
            WhereOr("company", mb.Like, "%"+filterDto.Keyword+"%").
            WhereOr("subject", mb.Like, "%"+filterDto.Keyword+"%").
            WhereOr("message", mb.Like, "%"+filterDto.Keyword+"%")

        return &queryGroup
    })
    return &query
})
```

**Rules:**
- Use `WhereGroup()` for OR conditions within an AND context
- First condition uses `Where()`, subsequent use `WhereOr()`
- Always return `&queryGroup` at the end
- Commonly used for keyword searches across multiple fields

### 4.2 When() for Conditional Filtering

Use `When()` to conditionally add filters based on input parameters.

**Pattern:**
```go
builder.When(condition, func(query mb.WhereBuilder) *mb.WhereBuilder {
    query.Where("field", operator, value)
    return &query
})
```

**Real Examples:**

**Example 1: Product Filters**
```go
builder := dbInstance.Select(fmt.Sprintf("%s.*", models.TableProduct)).
    When(filterDto.IsFeatured != nil, func(query mb.WhereBuilder) *mb.WhereBuilder {
        query.Where("is_featured", mb.Eq, *filterDto.IsFeatured)
        return &query
    }).
    When(filterDto.Price != nil, func(query mb.WhereBuilder) *mb.WhereBuilder {
        query.Where("price", mb.Eq, *filterDto.Price)
        return &query
    }).
    When(filterDto.Currency != nil, func(query mb.WhereBuilder) *mb.WhereBuilder {
        query.Where("currency", mb.Eq, *filterDto.Currency)
        return &query
    }).
    When(filterDto.VarietyID != nil, func(query mb.WhereBuilder) *mb.WhereBuilder {
        query.Where("variety_id", mb.Eq, *filterDto.VarietyID)
        return &query
    }).
    When(filterDto.BrandID != nil, func(query mb.WhereBuilder) *mb.WhereBuilder {
        query.Where("brand_id", mb.Eq, *filterDto.BrandID)
        return &query
    }).
    When(filterDto.IsActive != nil, func(query mb.WhereBuilder) *mb.WhereBuilder {
        query.Where("is_active", mb.Eq, *filterDto.IsActive)
        return &query
    })
```

**Example 2: Order Filters**
```go
builder := mb.Instance().Select(fmt.Sprintf("%s.*", models.TableOrder)).
    When(filterDto.UserID > 0, func(query mb.WhereBuilder) *mb.WhereBuilder {
        query.Where("user_id", mb.Eq, filterDto.UserID)
        return &query
    }).
    When(filterDto.Keyword != "", func(query mb.WhereBuilder) *mb.WhereBuilder {
        query.WhereGroup(func(queryGroup mb.WhereBuilder) *mb.WhereBuilder {
            queryGroup.Where("note", mb.Like, "%"+filterDto.Keyword+"%")
            return &queryGroup
        })
        return &query
    })
```

**Example 3: Content Filters**
```go
builder := dbInstance.Select(fmt.Sprintf("%s.*", models.TableContent)).
    When(queryDto.TypeID > 0, func(query mb.WhereBuilder) *mb.WhereBuilder {
        query.Where("type_id", mb.Eq, queryDto.TypeID)
        return &query
    }).
    When(queryDto.TaxonomyID > 0, func(query mb.WhereBuilder) *mb.WhereBuilder {
        query.Where("taxonomy_id", mb.Eq, queryDto.TaxonomyID)
        return &query
    }).
    When(filterDto.Keyword != "", func(query mb.WhereBuilder) *mb.WhereBuilder {
        query.WhereGroup(func(queryGroup mb.WhereBuilder) *mb.WhereBuilder {
            queryGroup.
                Where("title", mb.Like, "%"+filterDto.Keyword+"%").
                WhereOr("description", mb.Like, "%"+filterDto.Keyword+"%").
                WhereOr("content", mb.Like, "%"+filterDto.Keyword+"%")
            return &queryGroup
        })
        return &query
    })
```

**Rules:**
- Use `When()` for optional filters that depend on user input
- Check for non-zero/non-empty values before applying filters
- Chain multiple `When()` calls for different optional filters
- Always return `&query` at the end of the callback
- Use pointer fields in DTOs for optional filters (nil = not provided)

### 4.3 Complex Filtering Combinations

Combining multiple filtering techniques.

**Example: Rating Criteria with Multiple Filters**
```go
builder := dbInstance.Select(fmt.Sprintf("%s.*", models.TableRatingCriteria)).
    // Fixed filter
    When(filterDto.ScaleID != 0, func(query mb.WhereBuilder) *mb.WhereBuilder {
        query.Where("scale_id", mb.Eq, filterDto.ScaleID)
        return &query
    }).
    // Keyword search with OR logic
    When(filterDto.Keyword != "", func(query mb.WhereBuilder) *mb.WhereBuilder {
        query.WhereGroup(func(queryGroup mb.WhereBuilder) *mb.WhereBuilder {
            queryGroup.
                Where("name", mb.Like, "%"+filterDto.Keyword+"%").
                WhereOr("description", mb.Like, "%"+filterDto.Keyword+"%")

            return &queryGroup
        })
        return &query
    }).
    Limit(filterDto.PerPage, offset)
```

---

## 5. Pagination & Ordering

### 5.1 Pagination Pattern

Standard pagination implementation used across all list endpoints.

**Pattern:**
```go
var offset = 0

// Calculate offset
if filterDto.Page > 0 {
    offset = (filterDto.Page - 1) * filterDto.PerPage
}

builder.Limit(filterDto.PerPage, offset)
```

**Real Examples:**

**Example 1: Product Pagination**
```go
var offset = 0

// Calculate offset
if filterDto.Page > 0 {
    offset = (filterDto.Page - 1) * filterDto.PerPage
}

// ... query building ...

Limit(filterDto.PerPage, offset)
```

**Example 2: User Pagination**
```go
var offset = 0

if filterDto.Page > 0 {
    offset = (filterDto.Page - 1) * filterDto.PerPage
}

// ... query building ...

Limit(filterDto.PerPage, offset)
```

**Rules:**
- Default offset is 0 (first page)
- Formula: `offset = (Page - 1) * PerPage`
- Apply `Limit()` at the end of query building
- `Limit()` takes two parameters: limit and offset

### 5.2 Ordering Patterns

#### Simple OrderBy

**Pattern:**
```go
builder.OrderBy("field_name", mb.Asc)  // Ascending
builder.OrderBy("field_name", mb.Desc) // Descending
```

**Real Examples:**

**Example 1: Order by Created Date**
```go
builder.OrderBy("created_at", mb.Desc)
```

**Example 2: Order by ID**
```go
builder.OrderBy("id", mb.Asc)
```

**Example 3: Order with Table Prefix**
```go
builder.OrderBy(fmt.Sprintf("%s.created_at", models.TableOrder), mb.Desc)
```

#### Dynamic Ordering from DTO

Standard pattern for user-controlled sorting.

**Pattern:**
```go
if filterDto.OrderBy != "" {
    // Default order direction
    direction := mb.Asc
    orderKey := filterDto.OrderBy

    // Parse prefix for descending order
    if strings.HasPrefix(orderKey, "-") {
        direction = mb.Desc
        orderKey = orderKey[1:]
    }

    // Define allowed order fields
    orderFields := core.Data{
        "id":         fmt.Sprintf("%s.id", models.TableName),
        "name":       fmt.Sprintf("%s.name", models.TableName),
        "created_at": fmt.Sprintf("%s.created_at", models.TableName),
    }

    if field, ok := orderFields[orderKey]; ok {
        builder.OrderBy(field.(string), direction)
    }
}
```

**Real Examples:**

**Example 1: Product Dynamic Ordering**
```go
if filterDto.OrderBy != "" {
    // Set default order options
    direction := mb.Asc
    orderKey := filterDto.OrderBy

    // Parse parameter prefix
    if strings.HasPrefix(orderKey, "-") {
        direction = mb.Desc
        orderKey = orderKey[1:]
    }

    // OrderBy:id, brand_id, name, created_at, updated_at, deleted_at
    orderFields := core.Data{
        "id":         fmt.Sprintf("%s.id", models.TableProduct),
        "brand_id":   fmt.Sprintf("%s.brand_id", models.TableProduct),
        "name":       fmt.Sprintf("%s.name", models.TableProduct),
        "created_at": fmt.Sprintf("%s.created_at", models.TableProduct),
        "updated_at": fmt.Sprintf("%s.updated_at", models.TableProduct),
        "deleted_at": fmt.Sprintf("%s.deleted_at", models.TableProduct),
    }

    if field, ok := orderFields[orderKey]; ok {
        builder.OrderBy(field.(string), direction)
    }
}
```

**Example 2: User Dynamic Ordering**
```go
if filterDto.OrderBy != "" {
    // Default order by
    direction := mb.Asc
    orderKey := filterDto.OrderBy

    if strings.HasPrefix(filterDto.OrderBy, "-") {
        orderKey = filterDto.OrderBy[1:]
        direction = mb.Desc
    }

    var orderByFields = core.Data{
        "id":          fmt.Sprintf("%s.id", models.TableUser),
        "email":       fmt.Sprintf("%s.email", models.TableUser),
        "first_name":  fmt.Sprintf("%s.first_name", models.TableUser),
        "last_name":   fmt.Sprintf("%s.last_name", models.TableUser),
        "phone":       fmt.Sprintf("%s.phone", models.TableUser),
        "status":      fmt.Sprintf("%s.status", models.TableUser),
        "last_access": fmt.Sprintf("%s.last_access_at", models.TableUser),
    }

    if field, ok := orderByFields[orderKey]; ok {
        builder.OrderBy(field.(string), direction)
    }
}
```

**Example 3: Multiple OrderBy**
```go
OrderBy(fmt.Sprintf("%s.created_at", models.TableRating), mb.Desc).
OrderBy(fmt.Sprintf("%s.rating_value", models.TableRating), mb.Desc)
```

#### Default Ordering

Always provide a default order when no user preference specified.

**Example 1: Default Descending**
```go
// OrderBy by default
builder.OrderBy("created_at", mb.Desc)
```

**Example 2: Conditional Default**
```go
if filterDto.OrderBy != "" {
    // ... custom ordering ...
} else {
    // Default order by created_at desc
    builder.OrderBy(fmt.Sprintf("%s.created_at", models.TableOrder), mb.Desc)
}
```

**Rules:**
- Prefix `-` indicates descending order (e.g., "-created_at")
- Always whitelist allowed order fields for security
- Provide sensible default ordering
- Use table-prefixed column names for joins
- Multiple OrderBy() calls are supported (applied in order)

### 5.3 GroupBy

Used to group results by a column.

**Pattern:**
```go
builder.GroupBy("table.column")
```

**Real Example:**

**Example: Deduplicate Related Products**
```go
// Group by product ID to avoid duplicates
builder.GroupBy(models.TableProduct + ".id")
```

---

## 6. Joins & Relationships

### 6.1 InnerJoin Syntax

**Pattern:**
```go
Join(mb.InnerJoin, "target_table", mb.Condition{
    Field: "target_table.foreign_key",
    Opt:   mb.Eq,
    Value: mb.ValueField("source_table.primary_key"),
})
```

**Components:**
- **Join Type**: `mb.InnerJoin` (only INNER JOIN is shown in examples)
- **Target Table**: Table to join with
- **Condition**: Join condition struct with:
  - **Field**: Column from target table
  - **Opt**: Comparison operator (usually `mb.Eq`)
  - **Value**: Reference to source column using `mb.ValueField()`

### 6.2 Simple Join Examples

**Example 1: Product Categories Join**
```go
builder := dbInstance.Select(fmt.Sprintf("%s.*", models.TableProduct)).
    Join(mb.InnerJoin, models.TableProductCategory,
        mb.Condition{
            Field: models.TableProductCategory + ".product_id",
            Opt:   mb.Eq,
            Value: mb.ValueField(models.TableProduct + ".id"),
        })
```

**Example 2: Cart Items with Cart**
```go
_, err = mb.Instance().Select(models.TableCartItem+".*").
    Join(mb.InnerJoin, models.TableCart,
        mb.Condition{
            Field: models.TableCart + ".id",
            Opt:   mb.Eq,
            Value: mb.ValueField(models.TableCartItem + ".cart_id"),
        }).
    Where(models.TableCart+".user_id", mb.Eq, userID).
    OrderBy(models.TableCartItem+".id", mb.Desc).
    Find(&items)
```

**Example 3: Verify Cart Item Ownership**
```go
err = mb.Instance().Select(models.TableCartItem+".*").
    Join(mb.InnerJoin, models.TableCart,
        mb.Condition{
            Field: models.TableCart + ".id",
            Opt:   mb.Eq,
            Value: mb.ValueField(models.TableCartItem + ".cart_id"),
        }).
    Where(models.TableCartItem+".id", mb.Eq, cartItemID).
    Where(models.TableCart+".user_id", mb.Eq, userID).
    First(&cartItem)
```

### 6.3 Multiple Joins

**Example 1: User with Roles**
```go
builder := dbInstance.Select("DISTINCT users.id", "users.*").
    Join(mb.InnerJoin, "user_roles", mb.Condition{
        Field: "user_roles.user_id",
        Opt:   mb.Eq,
        Value: mb.ValueField("users.id"),
    }).
    Join(mb.InnerJoin, "roles", mb.Condition{
        Field: "roles.id",
        Opt:   mb.Eq,
        Value: mb.ValueField("user_roles.role_id"),
    })
```

**Example 2: Ratings with Users and Items**
```go
builder := mb.Instance().
    Select(
        fmt.Sprintf("%s.*", models.TableRating),
        fmt.Sprintf("%s.email", models.TableUser),
        fmt.Sprintf("%s.first_name", models.TableUser),
        fmt.Sprintf("%s.last_name", models.TableUser),
        fmt.Sprintf("%s.avatar", models.TableUser),
    ).
    Join(mb.InnerJoin, models.TableUser, mb.Condition{
        Field: fmt.Sprintf("%s.id", models.TableUser),
        Opt:   mb.Eq,
        Value: mb.ValueField(fmt.Sprintf("%s.user_id", models.TableRating)),
    }).
    Join(mb.InnerJoin, models.TableRatingItem, mb.Condition{
        Field: fmt.Sprintf("%s.id", models.TableRatingItem),
        Opt:   mb.Eq,
        Value: mb.ValueField(fmt.Sprintf("%s.item_id", models.TableRating)),
    })
```

**Example 3: Rating Stats by Pair**
```go
err := mb.Instance().
    Select(fmt.Sprintf("%s.*", models.TableRatingStats)).
    Join(mb.InnerJoin, models.TableRatingItem, mb.Condition{
        Field: fmt.Sprintf("%s.item_id", models.TableRatingStats),
        Opt:   mb.Eq,
        Value: mb.ValueField(fmt.Sprintf("%s.id", models.TableRatingItem)),
    }).
    Where(fmt.Sprintf("%s.item_type", models.TableRatingItem), mb.Eq, itemType).
    Where(fmt.Sprintf("%s.item_ref", models.TableRatingItem), mb.Eq, itemRef).
    First(&ratingStats)
```

### 6.4 Join with Filtering

**Example: Related Products by Categories**
```go
builder := dbInstance.Select(fmt.Sprintf("%s.*", models.TableProduct)).
    Join(mb.InnerJoin, models.TableProductCategory,
        mb.Condition{
            Field: models.TableProductCategory + ".product_id",
            Opt:   mb.Eq,
            Value: mb.ValueField(models.TableProduct + ".id"),
        })

// Add where conditions for category IDs
builder.Where(models.TableProductCategory+".category_id", mb.In, categoryIDs)

// Exclude the original product
builder.Where(models.TableProduct+".id", mb.NotEq, productID)

// Group by product ID to avoid duplicates
builder.GroupBy(models.TableProduct + ".id")

// Order by product ID in descending order
builder.OrderBy(models.TableProduct+".id", mb.Desc)

// Limit the number of results
builder.Limit(limit, 0)

// Execute the query
_, err = builder.Find(&relatedProducts)
```

**Rules:**
- Always use table-prefixed column names in joins
- Use `mb.ValueField()` to reference columns (not literal values)
- Apply WHERE clauses after joins
- Use DISTINCT or GroupBy to avoid duplicates
- Select specific columns to optimize performance

---

## 7. Transactions

### 7.1 Transaction Pattern with Try/Catch

Standard transaction pattern used for atomic operations.

**Pattern:**
```go
db := mb.Instance()
var err error

try.Perform(func() {
    // Begin transaction
    db.Begin()

    // Perform operations
    if err := db.Create(&model1); err != nil {
        try.Throw(err)
    }

    if err := db.Update(&model2); err != nil {
        try.Throw(err)
    }

    if err := db.Delete(&model3); err != nil {
        try.Throw(err)
    }

    // Commit transaction
    err = db.Commit()
    if err != nil {
        try.Throw(err)
    }
}).Catch(func(e try.E) {
    err = e.(error)
    log.Error(err)
    _ = db.Rollback() // Rollback on error
})

return err
```

**Real Example: Generate Order**
```go
db := mb.Instance()
var order models.Order
var err error

try.Perform(func() {
    // Begin transaction
    db.Begin()

    // Step 1: Verify and create address if needed
    var address models.Address
    err := db.
        Where("id", mb.Eq, dtoData.AddressID).
        Where("user_id", mb.Eq, user.ID).
        First(&address)

    if err != nil || address.ID == 0 {
        // Create new address
        address = models.Address{
            UserID:       user.ID,
            Type:         types.AddressTypeShipping,
            // ... other fields
        }

        if err := db.Create(&address); err != nil {
            log.Errorf("Failed to create address: %v", err)
            try.Throw(errors.New("Failed to create address"))
        }
    }

    // Step 2: Get cart
    cart, err := mb.GetModelByID[models.Cart](dtoData.CartID)
    if err != nil {
        log.Errorf("Failed to get cart by ID: %v", err)
        try.Throw(errors.New("Failed to get cart"))
    }

    cartItems := GetItemsByCartID(cart.ID)
    if len(cartItems) == 0 {
        try.Throw(errors.New("Cart is empty"))
    }

    // Step 3: Create order
    order = models.Order{
        UserID:     user.ID,
        TotalPrice: 1,
        CreatedAt:  time.Now(),
    }

    if err := db.Create(&order); err != nil {
        log.Errorf("Failed to create order: %v", err)
        try.Throw(errors.New("Failed to create order"))
    }

    // Step 4: Create order items
    var totalPrice float64 = 0
    for _, item := range cartItems {
        product, err := mb.GetModelByID[models.Product](item.ProductID)
        if err != nil {
            log.Errorf("Failed to get product by ID: %v", err)
            try.Throw(errors.New("Failed to get product"))
        }

        productPrice := product.Price
        totalPrice += float64(item.Quantity) * productPrice

        orderItem := models.OrderItem{
            OrderID:   order.ID,
            ProductID: item.ProductID,
            Price:     productPrice,
            Quantity:  int16(item.Quantity),
            CreatedAt: time.Now(),
        }

        if err := db.Create(&orderItem); err != nil {
            log.Errorf("Failed to create order item: %v", err)
            try.Throw(errors.New("Failed to create order item"))
        }
    }

    // Step 5: Update order total
    order.TotalPrice = totalPrice
    if err := db.Update(order); err != nil {
        log.Errorf("Failed to update order total price: %v", err)
        try.Throw(errors.New("Failed to update order total price"))
    }

    // Step 6: Create delivery records
    deliveryAddress := models.DeliveryAddress{
        Email:     address.Email,
        Phone:     address.Phone,
        AddressID: sql.NullInt32{Int32: int32(address.ID), Valid: address.ID > 0},
        // ... other fields
        CreatedAt: time.Now(),
    }

    if err := db.Create(&deliveryAddress); err != nil {
        log.Errorf("Failed to create delivery address: %v", err)
        try.Throw(errors.New("Failed to create delivery address"))
    }

    delivery := models.Delivery{
        OrderID:           order.ID,
        DeliveryAddressID: deliveryAddress.ID,
        Type:              types.DeliveryTypeBase,
        Status:            types.DeliveryStatusPrepare,
        CreatedAt:         time.Now(),
    }

    if err := db.Create(&delivery); err != nil {
        log.Errorf("Failed to create delivery: %v", err)
        try.Throw(errors.New("Failed to create delivery"))
    }

    // Step 7: Create delivery items
    for _, item := range cartItems {
        deliveryItem := models.DeliveryItem{
            OrderID:    order.ID,
            ProductID:  item.ProductID,
            DeliveryID: delivery.ID,
            Quantity:   uint8(item.Quantity),
            CreatedAt:  time.Now(),
        }

        if err := db.Create(&deliveryItem); err != nil {
            log.Errorf("Failed to create delivery item: %v", err)
            try.Throw(errors.New("Failed to create delivery item"))
        }
    }

    // Step 8: Delete the cart
    if err := db.Where("id", mb.Eq, cart.ID).Delete(models.Cart{}); err != nil {
        log.Errorf("Failed to delete cart: %v", err)
        try.Throw(errors.New("Failed to delete cart"))
    }

    // Commit the transaction
    err = db.Commit()
    if err != nil {
        try.Throw(err)
    }

    // Step 9: Send notification (outside transaction)
    err = notification.Send(notifications.NewOrder{
        Order:           order,
        User:            user,
        DeliveryAddress: deliveryAddress,
    })

    if err != nil {
        log.Errorf("Failed to send notification: %v", err)
        // Don't throw - order created successfully
    }
}).Catch(func(e try.E) {
    err = e.(error)
    log.Errorf("Error generating order: %v", err)
    _ = db.Rollback() // Roll back the transaction
})

if err != nil {
    return nil, err
}

return &order, nil
```

### 7.2 Transaction Operations

#### Within Transaction: Create
```go
if err := db.Create(&model); err != nil {
    log.Errorf("Failed to create: %v", err)
    try.Throw(errors.New("Failed to create"))
}
```

#### Within Transaction: Update
```go
if err := db.Update(model); err != nil {
    log.Errorf("Failed to update: %v", err)
    try.Throw(errors.New("Failed to update"))
}
```

#### Within Transaction: Delete
```go
if err := db.Where("id", mb.Eq, id).Delete(models.Model{}); err != nil {
    log.Errorf("Failed to delete: %v", err)
    try.Throw(errors.New("Failed to delete"))
}
```

#### Within Transaction: Query
```go
var model models.Model
err := db.Where("id", mb.Eq, id).First(&model)
if err != nil {
    try.Throw(errors.New("Record not found"))
}
```

**Rules:**
- Always use `db.Begin()` to start transaction
- Use `try.Throw()` to trigger rollback on errors
- Call `db.Commit()` explicitly at the end
- Rollback is automatic in the Catch block
- Log all errors with context
- Non-critical operations (like notifications) should be outside transactions
- Use descriptive error messages
- Handle partial failures gracefully

---

## 8. Raw SQL Queries

For complex queries that are difficult to express with the query builder, use raw SQL.

### 8.1 Raw() Method

**Pattern:**
```go
query := `SELECT ... FROM ... WHERE ...`
_, err := mb.Instance().Raw(query).Find(&results)
```

**Real Example: Recursive CTE for Content Tree**
```go
var contentList []models.Content
var dbInstance = mb.Instance()

query := `
WITH RECURSIVE content_tree AS (
    SELECT DISTINCT id, created_by, content_id, taxonomy_id, type_id, title, description, content, link, file,
    weight, created_at, updated_at FROM contents
    UNION ALL
    SELECT DISTINCT c.id, c.created_by, c.content_id, c.taxonomy_id, c.type_id, c.title, c.description, c.content, c.link,
    c.file, c.weight, c.created_at, c.updated_at FROM contents c
    INNER JOIN content_tree ct ON c.content_id = ct.id
)
SELECT DISTINCT id, created_by, content_id, taxonomy_id, type_id, title, description, content, link, file,
weight, created_at, updated_at
FROM content_tree
ORDER BY id ASC
`

_, err := dbInstance.Raw(query).Find(&contentList)
if err != nil {
    log.Error(err)
    return nil, 0, errors.New("Failed to fetch contents")
}
```

**When to Use Raw SQL:**
- Recursive CTEs (Common Table Expressions)
- Complex aggregations
- Window functions
- Database-specific features
- Performance-critical queries that need optimization

**Rules:**
- Use only when query builder is insufficient
- Document the purpose of complex queries
- Test raw queries thoroughly
- Be aware of SQL injection risks (use parameterized queries)
- Consider maintainability

---

## 9. Error Handling Patterns

### 9.1 Try/Catch Pattern

Standard error handling pattern using gflydev/core/try.

**Pattern:**
```go
try.Perform(func() {
    // Operations that may fail
    _, err = mb.Instance().Where("field", mb.Eq, value).Find(&items)
    if err != nil {
        try.Throw(err)
    }
}).Catch(func(e try.E) {
    log.Error(e)
    err = e.(error)
})
```

**Real Examples:**

**Example 1: Repository Pattern**
```go
try.Perform(func() {
    _, err = mb.Instance().Select(models.TableAddress+".*").
        Join(mb.InnerJoin, models.TableUser,
            mb.Condition{
                Field: models.TableAddress + ".user_id",
                Opt:   mb.Eq,
                Value: mb.ValueField(models.TableUser + ".id"),
            }).
        Where(models.TableUser+".id", mb.Eq, userID).
        Find(&items)
}).Catch(func(e try.E) {
    log.Error(e)
    err = e.(error)
})
```

**Example 2: Service with Finally**
```go
try.Perform(func() {
    // Create builder which query list orders
    builder := mb.Instance().Select(fmt.Sprintf("%s.*", models.TableOrder)).
        When(filterDto.UserID > 0, func(query mb.WhereBuilder) *mb.WhereBuilder {
            query.Where("user_id", mb.Eq, filterDto.UserID)
            return &query
        })

    // ... more query building ...

    // Query data
    total, err = builder.Find(&ordersList)
}).Finally(func() {
    if err != nil {
        log.Errorf("Error querying orders: %v", err)
    }
})
```

### 9.2 Direct Error Handling

Simple error checking without try/catch.

**Pattern:**
```go
model, err := mb.GetModelByID[models.Model](id)
if err != nil {
    log.Error(err)
    return nil, errors.ItemNotFound
}
```

**Real Examples:**

**Example 1: Get and Return Error**
```go
product, err := mb.GetModelByID[models.Product](updateProductDto.ID)
if err != nil {
    log.Error(err)
    return nil, errors.ItemNotFound
}
```

**Example 2: Create and Handle Error**
```go
try.Perform(func() {
    err = mb.CreateModel(cart)
}).Catch(func(e try.E) {
    err = e.(error)
    log.Error(err)
})

return cart, nil
```

**Example 3: Query with Error Check**
```go
total, err := mb.Instance().Where("cart_id", mb.Eq, cartID).
    OrderBy("created_at", mb.Desc).
    Find(&cartItems)

// Return an empty slice if no carts with items are found
if err != nil || total == 0 {
    cartItems = []models.CartItem{}
}
```

### 9.3 Error Logging

**Levels:**
- `log.Error(err)` - General errors
- `log.Errorf("Context: %v", err)` - Errors with context
- `log.Tracef("Success message")` - Success operations

**Real Examples:**

**Example 1: Error with Context**
```go
if err := db.Create(&address); err != nil {
    log.Errorf("Failed to create address: %v", err)
    try.Throw(errors.New("Failed to create address"))
}
```

**Example 2: Success Trace**
```go
log.Tracef("Successfully unlinked OAuth record %d for user %d", oauthID, userID)
```

**Example 3: Simple Error Log**
```go
if err != nil {
    log.Error(err)
    return nil, errors.ItemNotFound
}
```

### 9.4 Standard Error Returns

**Pattern 1: ItemNotFound**
```go
if err != nil {
    log.Error(err)
    return nil, errors.ItemNotFound
}
```

**Pattern 2: Custom Error**
```go
if err != nil {
    log.Error(err)
    return nil, errors.New("descriptive error message")
}
```

**Pattern 3: Internal Server Error**
```go
if err != nil {
    log.Error(err)
    return nil, errors.New("internal server error")
}
```

**Rules:**
- Always log errors before returning
- Use `errors.ItemNotFound` for missing records
- Provide descriptive error messages
- Use try/catch for complex operations
- Use direct error handling for simple operations
- Include context in error messages
- Don't expose internal details in user-facing errors

---

## 10. Best Practices & Rules

### 10.1 Service Layer Responsibilities

**DO:**
- Validate DTOs before database operations
- Fetch and verify related entities
- Use transactions for multiple related operations
- Log all errors with context
- Return domain errors (not database errors)
- Handle null values with dbNull package
- Set timestamps (CreatedAt, UpdatedAt)

**DON'T:**
- Expose database errors directly to controllers
- Skip validation of foreign keys
- Forget to rollback transactions on errors
- Use raw SQL unless necessary
- Skip error logging
- Return nil without checking errors

### 10.2 Generic Helper Patterns

#### Repository Helpers
```go
func (r *repository) getBy(field string, value any) *models.Model {
    model, err := mb.GetModelBy[models.Model](field, value)
    if err != nil {
        log.Error(err)
        return nil
    }
    return model
}

func (r *repository) GetByEmail(email string) *models.Model {
    return r.getBy("email", email)
}

func (r *repository) GetByID(id int) *models.Model {
    return r.getBy("id", id)
}
```

**Reference:** `internal/domain/repository/user_repository.go`

### 10.3 DTO Validation Pattern

Always validate DTOs before database operations.

**Example:**
```go
func CreateProduct(createDTO dto.CreateProduct) (*models.Product, error) {
    // Validate brand
    brand, err := mb.GetModelByID[models.Brand](createDTO.BrandID)
    if err != nil {
        log.Error(err)
        return nil, errors.ItemNotFound
    }

    // Validate currency
    currency := types.CurrencyUnit(createDTO.Currency)
    if !slices.Contains(types.CurrencyUnits, currency) {
        return nil, errors.InvalidParameter
    }

    // Validate variety
    variety, err := mb.GetModelByID[models.Variety](createDTO.VarietyID)
    if err != nil {
        log.Error(err)
        return nil, errors.ItemNotFound
    }

    // Validate attributes
    for _, val := range createDTO.Attributes {
        _, err := mb.GetModelByID[models.Attribute](val.ID)
        if err != nil {
            return nil, errors.InvalidParameter
        }
    }

    // Now create the product
    // ...
}
```

**Reference:** `internal/services/product_services.go`

### 10.4 Timestamp Management

**CreatedAt:**
```go
model := &models.Model{
    // ... fields ...
    CreatedAt: time.Now(),
}
```

**UpdatedAt:**
```go
model.UpdatedAt = dbNull.TimeNow()  // For nullable fields
// or
model.UpdatedAt = time.Now()  // For non-nullable fields
```

### 10.5 Null Value Handling

Use the `dbNull` package for nullable database fields.

**String:**
```go
dbNull.String(value)          // Create nullable string
dbNull.StringVal(nullString)  // Extract value
dbNull.StringNil(nullString)  // Convert to pointer/nil
```

**Int:**
```go
dbNull.Int32(value)          // Create nullable int32
```

**Time:**
```go
dbNull.TimeNow()             // Current time as nullable
dbNull.TimeNil(nullTime)     // Convert to pointer/nil
```

**Reference:** Throughout the codebase, especially:
- `internal/services/product_services.go`
- `internal/services/user_services.go`
- `internal/services/coupon_services.go`

### 10.6 Conditional Updates

Only update fields that have changed.

**Pattern:**
```go
if updateDto.Name != "" && updateDto.Name != model.Name {
    model.Name = updateDto.Name
}

if updateDto.Price > 0 && updateDto.Price != model.Price {
    model.Price = updateDto.Price
}

if updateDto.IsActive != nil && *updateDto.IsActive != model.IsActive {
    model.IsActive = *updateDto.IsActive
}
```

**Reference:** `internal/services/product_services.go`

### 10.7 Query Optimization

**Select Specific Columns:**
```go
// Good: Select only needed columns
Select(fmt.Sprintf("%s.*", models.TableProduct))

// Better: Select specific columns for large tables
Select("id", "name", "price", "created_at")
```

**Use Joins Wisely:**
```go
// Good: Join with filtering
Select(models.TableProduct + ".*").
    Join(mb.InnerJoin, models.TableCategory, condition).
    Where(models.TableCategory + ".is_active", mb.Eq, true)
```

**Pagination:**
```go
// Always paginate list endpoints
Limit(perPage, offset)
```

**Indexing:**
- Use indexed columns in WHERE clauses
- Use indexed columns for joins
- Use indexed columns for ORDER BY

### 10.8 Security Considerations

**Whitelist Order Fields:**
```go
orderFields := core.Data{
    "id":         fmt.Sprintf("%s.id", models.TableProduct),
    "name":       fmt.Sprintf("%s.name", models.TableProduct),
    "created_at": fmt.Sprintf("%s.created_at", models.TableProduct),
}

if field, ok := orderFields[orderKey]; ok {
    builder.OrderBy(field.(string), direction)
}
```

**Verify Ownership:**
```go
// Before deleting cart item, verify user owns it
err = mb.Instance().Select(models.TableCartItem+".*").
    Join(mb.InnerJoin, models.TableCart,
        mb.Condition{
            Field: models.TableCart + ".id",
            Opt:   mb.Eq,
            Value: mb.ValueField(models.TableCartItem + ".cart_id"),
        }).
    Where(models.TableCartItem+".id", mb.Eq, cartItemID).
    Where(models.TableCart+".user_id", mb.Eq, userID).
    First(&cartItem)
```

**Reference:** `internal/services/cart_service.go`

### 10.9 Testing Guidelines

**Test Setup:**
```go
func TestService(t *testing.T) {
    predictTesting()  // Initialize test environment
    // ... tests ...
}
```

**Reference:** Test files use `predictTesting()` from `test/setup_test.go`

### 10.10 Common Patterns Summary

| Pattern | Use Case |
|---------|----------|
| GetModelByID | Fetch by primary key |
| GetModelBy | Fetch by any field | 
| CreateModel | Insert new record | 
| UpdateModel | Update existing record | 
| DeleteModel | Delete by model | 
| Instance().Where().First() | Single record with conditions |
| Instance().Where().Find() | Multiple records | 
| When() | Conditional filtering |
| WhereGroup() | OR conditions | 
| Join() | Table relationships | 
| Transaction | Atomic operations |
| Try/Catch | Error handling | 

---

## Appendix: Quick Reference

### Import Statements
```go
import (
    mb "github.com/gflydev/db"
    dbNull "github.com/gflydev/db/null"
    "github.com/gflydev/core/log"
    "github.com/gflydev/core/errors"
    "github.com/gflydev/core/try"
)
```

### Operators
```go
mb.Eq       // =
mb.NotEq    // !=
mb.Like     // LIKE
mb.In       // IN
mb.Null     // IS NULL
mb.Asc      // ASC
mb.Desc     // DESC
mb.InnerJoin
```

### Common Functions
```go
mb.GetModelByID[T](id)
mb.GetModelBy[T](field, value)
mb.CreateModel(model)
mb.UpdateModel(model)
mb.DeleteModel(model)
mb.Instance()
mb.ValueField(column)
dbNull.String(value)
dbNull.Int32(value)
dbNull.TimeNow()
```

---

## Conclusion

This guide covers all Model Builder patterns found in the service layer of the ThietNgon e-commerce platform. Following these patterns ensures:

- **Consistency** across the codebase
- **Type Safety** with Go generics
- **Performance** through query optimization
- **Maintainability** with clear patterns
- **Security** through validation and ownership checks
- **Reliability** with proper error handling and transactions

For questions or clarifications, refer to the specific file references provided throughout this guide.
