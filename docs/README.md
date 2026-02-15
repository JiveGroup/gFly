# Documentation Hub

Welcome to the gFly project documentation. This directory contains comprehensive guides and references for developing with the gFly framework.

## Table of Contents

- [Quick Start Guides](#quick-start-guides)
- [Core Development Guides](#core-development-guides)
- [System Integration Guides](#system-integration-guides)
- [Framework Reference](#framework-reference)
- [Contributing](#contributing)

---

## Quick Start Guides

### [gFly Framework Best Practices](gFly.md)
Expert guide covering Go and gFly framework best practices, architecture patterns, and development guidelines.

**Topics:**
- Architecture patterns (Clean Architecture, DDD, microservices)
- Project structure guidelines
- Security and resilience patterns
- Testing and documentation standards

---

## Core Development Guides

These guides cover the essential patterns and workflows for building APIs with gFly.

### [CRUD API Development Guide](CRUD_API_GUIDE.md)
Complete guide for building Create, Read, Update, Delete (CRUD) operations.

**Topics:**
- HTTP helper functions (ProcessData, ProcessFilter, ProcessPathID, ProcessUpdateData)
- Create (POST), Read (GET), Update (PUT/PATCH), Delete (DELETE), List (GET Collection) operations
- Advanced filtering with custom parameters
- Response patterns and HTTP status codes
- Complete CRUD examples with best practices

**Use this guide when:** Building new API endpoints or refactoring existing CRUD operations.

---

### [Data Flow Architecture Guide](DATA_FLOW_GUIDE.md)
Comprehensive guide to the Request → DTO → Service → Model → Transformer → Response architecture.

**Topics:**
- Layered architecture overview
- DTO (Data Transfer Objects) patterns and validation
- Request objects with ToDto() conversion
- Response object structure
- Transformer patterns for Model → Response conversion
- Complete flow examples with multiple DTOs

**Use this guide when:** Understanding data flow through the application or implementing new features.

---

### [Model Builder (mb) Usage Guide](MODEL_BUILDER_GUIDE.md)
Complete reference for database operations using the gFly Model Builder.

**Topics:**
- Setup and initialization
- CRUD operations (CreateModel, GetModelByID, UpdateModel, DeleteModel)
- Query builder patterns (Where, Join, OrderBy, Limit)
- Advanced filtering (WhereGroup, When conditions)
- Pagination and ordering
- Joins and relationships
- Transactions with try/catch pattern
- Error handling patterns

**Use this guide when:** Writing database queries, implementing repositories, or optimizing data access.

---

### [Response Handling Guide](RESPONSE_GUIDE.md)
Reference for HTTP response patterns, types, and transformers.

**Topics:**
- Generic response types (List, Success, Error)
- Custom response type definitions
- Transformer patterns for domain models
- Helper functions (PathID, Parse, Validate, FilterData)
- HTTP status code conventions
- Complete examples for all response scenarios

**Use this guide when:** Returning data from controllers or creating new response types.

---

## System Integration Guides

### [Email Notification Guide](EMAIL_NOTIFICATION_GUIDE.md)
Guide for implementing email notifications in gFly applications.

**Topics:**
- Notification system architecture
- Creating notification structs with ToEmail() method
- Pongo2 email templates
- SMTP configuration and MailHog for development
- Sending notifications with error handling
- Best practices for email delivery

**Use this guide when:** Implementing user notifications, alerts, or transactional emails.

---

### [Storage System Guide](STORAGE_GUIDE.md)
Comprehensive guide to the unified storage abstraction layer.

**Topics:**
- Storage providers (local filesystem, S3-compatible cloud storage)
- Setup and configuration
- Basic storage operations (Put, Get, Delete, Exists)
- File upload patterns
- Presigned URLs for secure uploads/downloads
- File legitimization (temporary → permanent)
- Storage API endpoints
- Transformers and public URL generation

**Use this guide when:** Implementing file uploads, managing user content, or switching storage backends.

---

### [Storage Configuration](STORAGE_CONFIG.md)
Environment configuration reference for storage systems.

**Topics:**
- Common storage settings
- Local storage configuration
- CS3 (S3-compatible) storage configuration
- Environment variables reference
- Storage provider options

**Use this guide when:** Setting up storage for development or production environments.

---

## Framework Reference

### Architecture Flow

Understanding how data flows through the gFly application:

```
HTTP Request
    ↓
Request Layer (validation + ToDto())
    ↓
DTO Layer (validated data)
    ↓
Service Layer (business logic)
    ↓
Repository Layer (data access via Model Builder)
    ↓
Database
    ↓
Domain Model
    ↓
Transformer (Model → Response)
    ↓
Response Layer
    ↓
HTTP Response (JSON)
```

### Common Workflows

#### Building a New CRUD API Endpoint

1. **Define DTO** - Create data transfer object in `internal/dto/`
   - See: [DATA_FLOW_GUIDE.md](DATA_FLOW_GUIDE.md#dto-data-transfer-objects)

2. **Create Request** - Build request struct in `internal/http/request/`
   - See: [DATA_FLOW_GUIDE.md](DATA_FLOW_GUIDE.md#request-objects)

3. **Create Response** - Define response struct in `internal/http/response/`
   - See: [RESPONSE_GUIDE.md](RESPONSE_GUIDE.md#custom-response-types)

4. **Create Transformer** - Build transformer in `internal/http/transformers/`
   - See: [DATA_FLOW_GUIDE.md](DATA_FLOW_GUIDE.md#transformer-objects)

5. **Implement Service** - Add business logic in `internal/services/`
   - See: [MODEL_BUILDER_GUIDE.md](MODEL_BUILDER_GUIDE.md) for database operations

6. **Create Controller** - Build API controller in `internal/http/controllers/api/`
   - See: [CRUD_API_GUIDE.md](CRUD_API_GUIDE.md) for patterns

7. **Register Route** - Add route in `internal/http/routes/api_routes.go`
   - See: [CRUD_API_GUIDE.md](CRUD_API_GUIDE.md#best-practices)

#### Adding Email Notifications

1. **Create Notification Struct** - Define in `internal/notifications/`
2. **Implement ToEmail()** - Return email data structure
3. **Create Template** - Build Pongo2 template in `resources/views/mails/`
4. **Send Notification** - Use `notification.Send()` in service layer

See: [EMAIL_NOTIFICATION_GUIDE.md](EMAIL_NOTIFICATION_GUIDE.md)

#### Implementing File Uploads

1. **Configure Storage** - Set environment variables
2. **Create Upload Endpoint** - Use presigned URLs or direct upload
3. **Legitimize Files** - Move from temporary to permanent storage
4. **Transform URLs** - Include public URLs in responses

See: [STORAGE_GUIDE.md](STORAGE_GUIDE.md) and [STORAGE_CONFIG.md](STORAGE_CONFIG.md)

---

## Document Index

| Document | Purpose | Key Topics |
|----------|---------|------------|
| [gFly.md](gFly.md) | Framework best practices | Architecture, structure, security |
| [CRUD_API_GUIDE.md](CRUD_API_GUIDE.md) | API development patterns | CRUD operations, HTTP helpers |
| [DATA_FLOW_GUIDE.md](DATA_FLOW_GUIDE.md) | Data architecture | Request, DTO, Response, Transformer |
| [MODEL_BUILDER_GUIDE.md](MODEL_BUILDER_GUIDE.md) | Database operations | Queries, transactions, joins |
| [RESPONSE_GUIDE.md](RESPONSE_GUIDE.md) | Response handling | Response types, transformers |
| [EMAIL_NOTIFICATION_GUIDE.md](EMAIL_NOTIFICATION_GUIDE.md) | Email system | Notifications, templates |
| [STORAGE_GUIDE.md](STORAGE_GUIDE.md) | File storage | Uploads, presigned URLs |
| [STORAGE_CONFIG.md](STORAGE_CONFIG.md) | Storage configuration | Environment setup |

---

## Contributing

When adding new documentation:

1. Follow the existing structure with Table of Contents
2. Include practical examples for each concept
3. Add cross-references to related guides
4. Update this README.md with links to new documents
5. Use clear headings and code blocks for readability

---

## Need Help?

- Start with [gFly.md](gFly.md) for general framework understanding
- Use [CRUD_API_GUIDE.md](CRUD_API_GUIDE.md) for API development
- Refer to [DATA_FLOW_GUIDE.md](DATA_FLOW_GUIDE.md) for architecture questions
- Check [MODEL_BUILDER_GUIDE.md](MODEL_BUILDER_GUIDE.md) for database queries

For more information about the gFly framework, visit [https://www.gfly.dev](https://www.gfly.dev)
