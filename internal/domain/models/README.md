# Models Directory

This directory contains the data models used throughout the gFly application. Models represent the structure of data in the database and provide a way to interact with that data in a type-safe manner.

## Overview

Models in gFly follow a consistent structure:
- Each model corresponds to a database table
- Models use struct tags to define database mappings
- Constants are defined for table names and status values
- Models use the `mb.MetaData` field for table metadata

## Available Models

### User Model

The `User` model represents application users and is stored in the `users` table.

**Key Fields:**
- `ID`: Unique identifier (primary key)
- `Email`: User's email address
- `Password`: Hashed password
- `Fullname`: User's full name
- `Status`: User's status (active, pending, blocked)
- Various timestamp fields (created_at, updated_at, verified_at, etc.)

**Status Constants:**
- `UserStatusActive`: "active"
- `UserStatusPending`: "pending"
- `UserStatusBlocked`: "blocked"

### Role Model

The `Role` model represents user roles and is stored in the `roles` table.

**Key Fields:**
- `ID`: Unique identifier (primary key)
- `Name`: Role name
- `Slug`: URL-friendly role identifier
- `CreatedAt`: Timestamp when the role was created
- `UpdatedAt`: Timestamp when the role was last updated

**Role Types:**
- `RoleAdmin`: Administrator role
- `RoleModerator`: Moderator role
- `RoleMember`: Member role
- `RoleUser`: Regular user role
- `RoleGuest`: Guest role

The `RoleType` enum provides helper methods:
- `Name()`: Get the string name of a role
- `Ordinal()`: Get the numeric value of a role
- `Values()`: Get all role names
- `ByName()`: Get a role type by its string name

### UserRole Model

The `UserRole` model represents the many-to-many relationship between users and roles, stored in the `user_roles` table.

**Key Fields:**
- `ID`: Unique identifier (primary key)
- `RoleID`: Foreign key to the roles table
- `UserID`: Foreign key to the users table
- `CreatedAt`: Timestamp when the relationship was created

## Usage Example

```
// Creating a new user
user := models.User{
    Email:     "user@example.com",
    Password:  hashedPassword,
    Fullname:  "Example User",
    Phone:     "123-456-7890",
    Status:    models.UserStatusActive,
    CreatedAt: time.Now(),
    UpdatedAt: time.Now(),
}

// Using the model builder to save to database
err := mb.CreateModel(&user)
```

## Database Mapping

Models use struct tags to map to database columns:

```
type Example struct {
    MetaData mb.MetaData `db:"-" model:"table:examples"`
    ID       int         `db:"id" model:"name:id; type:serial,primary"`
    Name     string      `db:"name" model:"name:name"`
}
```

The `db` tag specifies the column name, while the `model` tag provides additional metadata for the model builder.
