# Data Transfer Objects (DTOs)

This directory contains Data Transfer Objects (DTOs) used throughout the application. DTOs are simple objects that carry data between processes or layers of the application.

## Purpose

- Define structures for transferring data between different parts of the application
- Provide a clear contract for data exchange
- Separate data transfer concerns from domain models

## Usage

DTOs are typically used for:
- API request and response payloads
- Transferring data between application layers
- Mapping between domain models and external representations

## Best Practices

- Keep DTOs simple and focused on data transfer
- Avoid adding business logic to DTOs
- Use validation tags when appropriate for request/response validation
- Consider using separate DTOs for requests and responses
