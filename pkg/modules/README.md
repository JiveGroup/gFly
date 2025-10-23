# Modules

This directory contains reusable, self-contained modules that provide specific functionality and can be shared across different gFly applications.

## Purpose

The modules directory is used for:
- Providing standalone, pluggable functionality
- Implementing reusable business features that can be extracted from internal modules
- Creating framework-level modules that follow gFly conventions
- Offering package-level modules that can be imported by external applications

## Usage

Modules should be:
- Self-contained with their own routes, services, and domain logic
- Framework-agnostic where possible, with clear dependencies
- Well-documented with usage examples and API documentation
- Thoroughly tested with unit and integration tests
- Versioned and maintained independently when appropriate

## Module Structure

Each module should follow this structure:

```
module_name/
├── README.md              # Module documentation and usage guide
├── api/                   # HTTP controllers
│   └── controller.go
├── dto/                   # Data Transfer Objects
│   ├── request.go
│   └── response.go
├── domain/                # Domain models (optional if module-specific)
│   ├── models/
│   └── repository/
├── middleware/            # Module-specific middleware
│   └── middleware.go
├── routes/                # Route registration
│   └── routes.go         # RegisterApi() and RegisterWeb() functions
├── services/              # Business logic
│   └── service.go
└── module.go              # Module entry point and configuration
```

## Organization

Consider creating modules for:
- **Authentication & Authorization**: JWT, OAuth, session management, RBAC
- **Notifications**: Email, SMS, push notifications
- **File Management**: Upload, storage, processing
- **Payment Processing**: Payment gateway integrations
- **Audit Logging**: Activity tracking and logging
- **Search**: Full-text search, filtering
- **Analytics**: Usage tracking, metrics

## Integration Pattern

Modules should provide a simple registration interface:

```go
// module.go
package mymodule

import "github.com/gflydev/core"

type Module struct {
    config Config
}

func New(config Config) *Module {
    return &Module{config: config}
}

func (m *Module) Register(app core.IFly) {
    // Register routes, middleware, services
}
```

Usage in main application:

```go
import "github.com/yourproject/pkg/modules/mymodule"

func main() {
    app := core.NewApp()

    // Register module
    myMod := mymodule.New(mymodule.Config{...})
    myMod.Register(app)

    app.Run()
}
```

## Best Practices

- Keep modules loosely coupled from the main application
- Use dependency injection for external dependencies (database, cache, etc.)
- Provide clear configuration options via Config structs
- Document all public APIs with GoDoc comments and Swagger annotations
- Write comprehensive tests including integration tests
- Follow Clean Architecture principles within each module
- Use interfaces to define contracts and enable testing
- Version modules semantically if extracted as separate packages
- Provide migration scripts if the module requires database tables

## Difference from internal/modules

- `pkg/modules/`: Reusable, framework-level modules that can be shared across projects or extracted as separate packages
- `internal/modules/`: Application-specific modules that are tightly coupled to this project

## Migration Path

To move a module from `internal/modules/` to `pkg/modules/`:

1. Remove application-specific dependencies
2. Parameterize configuration via Config struct
3. Add comprehensive documentation and examples
4. Ensure all tests pass in isolation
5. Add module.go with registration interface
6. Update imports in the main application
