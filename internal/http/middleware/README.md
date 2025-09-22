# Middleware

This directory contains HTTP middleware components that process requests and responses.

## Purpose

Middleware is responsible for:
- Processing HTTP requests before they reach controllers
- Processing HTTP responses before they are sent to clients
- Implementing cross-cutting concerns like authentication, logging, and CORS
- Modifying request or response data
- Terminating requests early when necessary (e.g., for unauthorized access)

## Usage

Middleware can be applied globally, to route groups, or to individual routes:

```
// Example pseudocode for middleware implementation and usage
// Note: This is not actual implementation code

// In a middleware file (auth_middleware.go):
func AuthMiddleware() MiddlewareFunc {
    return func(ctx *fasthttp.RequestCtx, next NextFunc) {
        // Get token from request
        token := ctx.Request.Header.Peek("Authorization")

        // Validate token
        if !isValidToken(token) {
            // Return unauthorized response
            ctx.Response.SetStatusCode(401)
            ctx.Response.SetBodyString("Unauthorized")
            return
        }

        // Token is valid, continue to next middleware or handler
        next(ctx)
    }
}

// In routes.go:
// Apply middleware globally
app.Use(middleware.Logger())

// Apply middleware to a route group
apiGroup := app.Group("/api")
apiGroup.Use(middleware.AuthMiddleware())

// Apply middleware to a specific route
app.Get("/admin", middleware.AdminOnly(), controllers.AdminDashboard)
```

## Common Middleware Types

- **Authentication**: Verifies user identity
- **Authorization**: Checks user permissions
- **CORS**: Handles Cross-Origin Resource Sharing
- **Logging**: Records request and response information
- **Rate Limiting**: Prevents abuse by limiting request frequency
- **Request ID**: Adds unique identifiers to requests for tracing
- **Recovery**: Catches panics and returns appropriate error responses
- **Compression**: Compresses response bodies
- **Caching**: Caches responses for improved performance

## Best Practices

- Keep middleware focused on a single responsibility
- Order middleware carefully based on dependencies
- Use middleware for cross-cutting concerns only
- Document middleware behavior and configuration options
- Handle errors appropriately within middleware
- Test middleware in isolation
- Consider performance implications, especially for global middleware
- Use dependency injection for middleware that requires external services
