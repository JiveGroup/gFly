# Repository Directory

This directory contains the repository implementations used throughout the gFly application. Repositories provide an abstraction layer between the domain models and the database, encapsulating the logic for retrieving and persisting data.

## Overview

Repositories in gFly follow a consistent structure:
- Each repository is defined by an interface
- Concrete implementations of these interfaces provide the actual data access logic
- A repository pool is used to collect and access all repositories

## Repository Pattern

The repository pattern is used to:
- Decouple business logic from data access logic
- Provide a consistent API for data access
- Make testing easier by allowing repository implementations to be mocked
- Centralize data access code to avoid duplication

## Available Repositories

### Role Repository

The `RoleRepository` implements the `IRoleRepository` interface and provides methods for working with user roles.

**Interface Methods:**
- `GetRolesByUserID(userID int) []models.Role`: Retrieves all roles assigned to a specific user
- `AddRoleForUserID(userID int, slug string) error`: Assigns a role to a user by role slug

**Implementation Details:**
- Uses the model builder (`mb`) for database operations
- Uses the query builder (`qb.) for constructing SQL queries
- Handles errors and logs them appropriately

## Repository Pool

The repository pool is a singleton that provides access to all repositories in the application:

```
// Repositories struct for collect all app repositories.
type Repositories struct {
    IRoleRepository
}

// Pool a repository pool to store all
var Pool = &Repositories{
    &RoleRepository{},
}
```

## Usage Example

```
// Get roles for a user
roles := repository.Pool.GetRolesByUserID(userId)

// Add a role to a user
err := repository.Pool.AddRoleForUserID(userId, "admin")
if err != nil {
    // Handle error
}
```

## Creating New Repositories

To create a new repository:

1. Define an interface for the repository
2. Create a struct that implements the interface
3. Add the interface to the `Repositories` struct
4. Initialize the implementation in the `Pool` variable

Example:

```
// Define the interface
type IUserRepository interface {
    GetUserByID(id int) (*models.User, error)
    // Other methods...
}

// Create the implementation
type UserRepository struct {
}

// Implement the methods
func (q *UserRepository) GetUserByID(id int) (*models.User, error) {
    // Implementation...
}

// Update the Repositories struct and Pool
type Repositories struct {
    IRoleRepository
    IUserRepository
}

var Pool = &Repositories{
    &RoleRepository{},
    &UserRepository{},
}
```
