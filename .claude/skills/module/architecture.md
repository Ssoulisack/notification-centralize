# Module-Based Architecture Guide

## Overview

Module-based architecture organizes code by **business domain/feature**. Each module is self-contained with all its layers (handler, service, repository, model) living together.

## Benefits of Module-Based Architecture

### 1. Better Cohesion
All code related to a single feature lives together, making it easier to understand and modify.

### 2. Reduced Coupling
Modules are more independent, reducing unintended dependencies between features.

### 3. Easier Navigation
When working on a feature, all related code is in one place instead of scattered across multiple layer directories.

### 4. Better Scalability
As the application grows, modules remain manageable and don't create deeply nested layer structures.

### 5. Team Collaboration
Different teams can own different modules with minimal conflicts.

### 6. Clearer Boundaries
Module boundaries naturally align with business domains, making it easier to enforce separation of concerns.

## Module Structure
```
modules/
  user/
    handler.go       # User HTTP handlers
    service.go       # User business logic
    repository.go    # User data access
    model.go         # User domain model
    dto.go          # User request/response DTOs
    routes.go       # User routes
  product/
    handler.go
    service.go
    repository.go
    model.go
    dto.go
    routes.go
  order/
    handler.go
    service.go
    repository.go
    model.go
    dto.go
    routes.go
```

**Structure:**
- Everything for a feature is in one place
- Clear module boundaries
- Easy to understand and navigate
- Changes to one module don't affect others
- Flat, shallow directory structure
- Self-contained and independent

## File Responsibilities

### handler.go
- HTTP request handling
- Request validation
- Response formatting
- Error handling
- Swagger documentation
- Calls service layer

### service.go
- Business logic
- Business rule validation
- Orchestration between repositories
- Transaction management
- Error handling
- Calls repository layer

### repository.go
- Database operations (CRUD)
- Query building
- Data mapping
- Uses GORM for ORM operations

### model.go
- Domain entity definition
- Database schema (GORM tags)
- Table name specification
- Entity relationships

### dto.go
- Request DTOs (CreateRequest, UpdateRequest)
- Response DTOs (Response, ListResponse)
- Validation tags
- JSON serialization tags

### routes.go
- Route registration
- Dependency injection setup
- Route grouping
- Middleware application

## Dependency Flow

```
HTTP Request
    ↓
Handler (handler.go)
    ↓
Service (service.go)
    ↓
Repository (repository.go)
    ↓
Database
```

Each layer only knows about the layer directly below it:
- Handler depends on Service
- Service depends on Repository
- Repository depends on Database

## Module Creation Checklist

When creating a new module:

- [ ] Create module directory in `modules/`
- [ ] Create handler.go with HTTP handlers
- [ ] Create service.go with business logic
- [ ] Create repository.go with data access
- [ ] Create model.go with domain entity
- [ ] Create dto.go with request/response structures
- [ ] Create routes.go with route registration
- [ ] Register module routes in `api/routes/setup.go`
- [ ] Run database migrations if needed
- [ ] Generate Swagger documentation
- [ ] Write tests

## Best Practices

### 1. Keep Modules Independent
Avoid importing types or functions from other modules. If you need to share, create a shared package in `core/` or `domain/`.

### 2. Use Dependency Injection
Always inject dependencies through constructors:
```go
func NewService(repo *Repository) *Service {
    return &Service{repo: repo}
}
```

### 3. Return Errors, Don't Panic
```go
// Good
func (s *Service) Create(req *CreateRequest) (*Response, error) {
    if err := s.validate(req); err != nil {
        return nil, err
    }
    // ...
}

// Bad
func (s *Service) Create(req *CreateRequest) *Response {
    if err := s.validate(req); err != nil {
        panic(err) // Don't do this
    }
    // ...
}
```

### 4. Use DTOs for API Boundaries
Never expose your domain models directly:
```go
// Good - Using DTO
func (h *Handler) Create(c *fiber.Ctx) error {
    var req CreateRequest // DTO
    // ...
}

// Bad - Exposing domain model
func (h *Handler) Create(c *fiber.Ctx) error {
    var model Model // Domain model
    // ...
}
```

### 5. Keep Business Logic in Services
```go
// Good
func (s *Service) CreateUser(req *CreateUserRequest) error {
    if !s.isEmailValid(req.Email) {
        return errors.New("invalid email")
    }
    // Business logic here
}

// Bad - Business logic in handler
func (h *Handler) CreateUser(c *fiber.Ctx) error {
    var req CreateUserRequest
    c.BodyParser(&req)
    if !strings.Contains(req.Email, "@") { // Don't do this
        return c.Status(400).JSON(...)
    }
    // ...
}
```

### 6. Use Proper HTTP Status Codes
- 200 OK - Successful GET, PUT
- 201 Created - Successful POST
- 204 No Content - Successful DELETE
- 400 Bad Request - Invalid input
- 404 Not Found - Resource not found
- 500 Internal Server Error - Server error

### 7. Add Swagger Documentation
Always document your API endpoints with Swagger comments.

## Example: Creating a Product Module

```bash
/module create product
```

This creates:
```
modules/product/
├── handler.go      # Product handlers
├── service.go      # Product business logic
├── repository.go   # Product data access
├── model.go        # Product entity
├── dto.go         # Product DTOs
└── routes.go      # Product routes
```

Then you can add fields to the Product model and DTOs:

```go
// model.go
type Model struct {
    ID          uint           `gorm:"primarykey" json:"id"`
    Name        string         `gorm:"size:255;not null" json:"name"`
    Price       float64        `gorm:"not null" json:"price"`
    Description string         `gorm:"type:text" json:"description"`
    CreatedAt   time.Time      `json:"created_at"`
    UpdatedAt   time.Time      `json:"updated_at"`
    DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// dto.go
type CreateRequest struct {
    Name        string  `json:"name" validate:"required,min=1,max=255"`
    Price       float64 `json:"price" validate:"required,gt=0"`
    Description string  `json:"description"`
}
```

## Starting with Module-Based

To build your project with module-based architecture:

1. Create `modules/` directory at project root
2. For each feature/domain:
   ```bash
   /module create user
   /module create product
   /module create order
   ```
3. Customize each module's model and DTOs
4. Implement business logic in services
5. Register module routes in `api/routes/setup.go`
6. Test each module independently

## Shared Code

Some code needs to be shared across modules:

### Core Utilities (`core/`)
- Encryption, hashing
- File operations
- Date/time formatting
- Custom validators

### Domain Constants (`domain/constants/`)
- Shared constants
- Enums
- Error codes

### Middleware (`api/middleware/`)
- Authentication
- Authorization
- Logging
- Rate limiting
- CORS

These remain in their current locations and can be imported by any module.

## Testing Modules

Each module should have its own tests:

```
modules/user/
├── handler.go
├── handler_test.go
├── service.go
├── service_test.go
├── repository.go
└── repository_test.go
```

This keeps tests close to the code they test.

## Summary

Module-based architecture provides:
- ✅ Code organized by feature/domain
- ✅ Related code stays together
- ✅ Low coupling between features
- ✅ Easy codebase navigation
- ✅ Excellent scalability
- ✅ Clear business domain boundaries
- ✅ Independent, testable modules
- ✅ Team ownership of domains
