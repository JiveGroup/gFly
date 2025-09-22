# Application Directory

This directory contains the core application code organized into modular components.

## Purpose

The app directory is the heart of the application, containing:
- Core business logic and domain models
- HTTP request handling and routing
- Command-line interface and scheduled tasks
- Services and utilities
- Error handling and notifications

## Structure

- **console/**: Command-line interface, scheduled tasks, and queue workers
  - **commands/**: CLI commands
  - **queues/**: Background job processing
  - **schedules/**: Scheduled tasks
- **constant/**: Application-wide constants and configuration values
- **domain/**: Core business logic and domain models
  - **models/**: Domain entities
  - **repository/**: Data access layer
- **dto/**: Data Transfer Objects for API requests and responses
- **errors/**: Custom error types and error handling
- **http/**: HTTP request handling and routing
  - **controllers/**: Request handlers
  - **middleware/**: HTTP middleware
  - **request/**: Request validation
  - **response/**: Response formatting
  - **routes/**: Route definitions
  - **transformers/**: Data transformers
- **notifications/**: Notification templates and delivery
- **services/**: Business logic services
- **utils/**: Utility functions and helpers

## Best Practices

- Follow the Single Responsibility Principle for all components
- Keep controllers thin, delegating business logic to services
- Use domain models to encapsulate business rules
- Implement proper error handling and validation
- Write tests for all components
- Document public APIs and interfaces
- Maintain separation of concerns between layers
- Use dependency injection for better testability
