# Project Conventions & Standards

This document outlines all conventions and standards for the Go Fiber project.

## Table of Contents
- [Architecture Patterns](#architecture-patterns)
- [Response Format](#response-format)
- [Validation](#validation)
- [Constants](#constants)
- [Error Handling](#error-handling)
- [Code Organization](#code-organization)

---

## Architecture Patterns

### Choose Your Pattern

| Pattern | When to Use | Skills |
|---------|-------------|--------|
| **Module-Based** | Large projects, DDD, varying complexity | `/module` |
| **Layer-Based** | Small-medium projects, traditional MVC | `/layer` |

### Module-Based Structure
```
modules/
  user/
    handler.go    - HTTP handlers
    service.go    - Business logic
    repository.go - Data access
    model.go      - Domain entity
    dto.go        - Request/Response DTOs
    routes.go     - Route registration
```

### Layer-Based Structure
```
api/controllers/       - HTTP handlers
data/services/         - Business logic
data/repositories/     - Data access
domain/models/         - Domain entities
domain/dto/           - Request/Response DTOs
domain/constants/     - Constants & enums
```

---

## Response Format

### ✅ Always Use Response Utility

**Location:** `core/response/`

**Never use `fiber.Map`** - Always use the response utility for consistent API responses.

### Standard Response Structure

```json
{
  "success": true,
  "message": "Optional message",
  "data": { ... },
  "error": "Error message if failed",
  "details": "Additional error details"
}
```

### Usage Examples

```go
import "yourproject/core/response"

// Success Responses
return response.OK(c, result)                              // 200 OK
return response.Created(c, "User created", result)         // 201 Created
return response.NoContent(c)                               // 204 No Content

// Error Responses
return response.BadRequest(c, "Invalid input", nil)        // 400
return response.ValidationError(c, err.Error())            // 400
return response.Unauthorized(c, "Authentication required") // 401
return response.Forbidden(c, "Access denied")              // 403
return response.NotFound(c, "User not found")              // 404
return response.Conflict(c, "User already exists")         // 409
return response.InternalServerError(c, err.Error(), nil)   // 500
```

### ❌ Don't Do This

```go
// Never use fiber.Map
return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
    "error": "Invalid request",
})
```

---

## Validation

### ✅ Always Validate DTOs

Every DTO **MUST** have validation tags using `go-playground/validator`.

### Basic Pattern

```go
import (
    "yourproject/core/response"
    "github.com/go-playground/validator/v10"
    "github.com/gofiber/fiber/v2"
)

var validate = validator.New()

func (h *Handler) Create(c *fiber.Ctx) error {
    var req CreateRequest

    // Parse
    if err := c.BodyParser(&req); err != nil {
        return response.BadRequest(c, "Invalid request body", nil)
    }

    // Validate
    if err := validate.Struct(&req); err != nil {
        return response.ValidationError(c, err.Error())
    }

    // Business logic...
}
```

### Common Validation Tags

```go
type CreateRequest struct {
    Name   string  `json:"name" validate:"required,min=3,max=255"`
    Email  string  `json:"email" validate:"required,email"`
    Age    int     `json:"age" validate:"required,gte=0,lte=150"`
    Price  float64 `json:"price" validate:"required,gt=0"`
    Status string  `json:"status" validate:"required,oneof=active inactive"`
    Phone  string  `json:"phone" validate:"omitempty,e164"`
}

type UpdateRequest struct {
    Name   *string  `json:"name,omitempty" validate:"omitempty,min=3,max=255"`
    Email  *string  `json:"email,omitempty" validate:"omitempty,email"`
    Age    *int     `json:"age,omitempty" validate:"omitempty,gte=0,lte=150"`
    Price  *float64 `json:"price,omitempty" validate:"omitempty,gt=0"`
    Status *string  `json:"status,omitempty" validate:"omitempty,oneof=active inactive"`
}
```

**Rules:**
- Create requests: Required fields use `required`
- Update requests: Use pointers with `omitempty` for optional fields
- Always validate enums with `oneof`
- Use `email`, `url`, `e164` for format validation

---

## Constants

### ✅ Use Constants for Static Values

**Location:** `domain/constants/<module>_constants.go`

### Pattern

```go
package constants

// Status constants
const (
    UserStatusActive   = "active"
    UserStatusInactive = "inactive"
    UserStatusBanned   = "banned"
)

// Error messages
const (
    ErrUserNotFound     = "user not found"
    ErrUserInvalidInput = "invalid user input"
    ErrUserAlreadyExists = "user already exists"
)

// Validation constants
const (
    UserNameMinLength = 3
    UserNameMaxLength = 255
    UserAgeMin       = 18
    UserAgeMax       = 150
)

// Helper functions
func GetValidUserStatuses() []string {
    return []string{
        UserStatusActive,
        UserStatusInactive,
        UserStatusBanned,
    }
}

func IsValidUserStatus(status string) bool {
    for _, s := range GetValidUserStatuses() {
        if s == status {
            return true
        }
    }
    return false
}
```

### Usage in Code

```go
import "yourproject/domain/constants"

// In service
func (s *Service) Create(req *CreateRequest) (*Response, error) {
    model := &Model{
        Status: constants.UserStatusActive, // Use constant
    }
    // ...
}

// In validation
func (s *Service) GetByID(id uint) (*Response, error) {
    model, err := s.repo.FindByID(id)
    if err != nil {
        return nil, errors.New(constants.ErrUserNotFound) // Use constant
    }
    return toResponse(model), nil
}
```

**Rules:**
- Always use constants for status enums
- Always use constants for error messages
- Use constants for validation rules
- Provide helper functions for validation

---

## Error Handling

### Pattern

```go
// Handler/Controller Layer
func (h *Handler) Create(c *fiber.Ctx) error {
    result, err := h.service.Create(&req)
    if err != nil {
        return response.InternalServerError(c, err.Error(), nil)
    }
    return response.Created(c, "User created", result)
}

// Service Layer
func (s *Service) Create(req *CreateRequest) (*Response, error) {
    if err := s.validateBusinessRules(req); err != nil {
        return nil, err // Return error, don't handle HTTP
    }

    model := &Model{Status: constants.UserStatusActive}
    if err := s.repo.Create(model); err != nil {
        return nil, err
    }

    return toResponse(model), nil
}

// Repository Layer
func (r *Repository) Create(model *Model) error {
    return r.db.Create(model).Error // Return GORM error
}
```

**Rules:**
- Handler/Controller: Convert errors to HTTP responses using response utility
- Service: Return business errors, don't handle HTTP
- Repository: Return database errors as-is
- Use constants for error messages
- Never panic, always return errors

---

## Code Organization

### Project Structure

```
yourproject/
├── api/
│   ├── middleware/      - CORS, JWT, Logger, etc.
│   └── routes/         - Route registration
├── bootstrap/          - App initialization
├── core/
│   ├── response/       - Response utility ✅
│   ├── utilities/      - Helper functions
│   └── logs/          - Logging utilities
├── domain/
│   ├── constants/      - Constants & enums ✅
│   ├── models/        - Domain entities (layer-based)
│   └── dto/           - DTOs (layer-based)
├── modules/           - Module-based features ✅
│   └── user/
│       ├── handler.go
│       ├── service.go
│       ├── repository.go
│       ├── model.go
│       ├── dto.go
│       └── routes.go
├── data/              - Layer-based structure
│   ├── services/
│   └── repositories/
├── config.yaml
├── main.go
└── Makefile
```

### Dependency Injection

Always inject dependencies through constructors:

```go
// Repository
func NewRepository(db *gorm.DB) *Repository {
    return &Repository{db: db}
}

// Service
func NewService(repo *Repository) *Service {
    return &Service{repo: repo}
}

// Handler
func NewHandler(service *Service) *Handler {
    return &Handler{service: service}
}

// Wire up in routes
repo := NewRepository(db)
service := NewService(repo)
handler := NewHandler(service)
```

---

## Naming Conventions

### Files
- Lowercase with underscores: `user_controller.go`, `product_service.go`
- Module-based: `handler.go`, `service.go`, `repository.go`

### Types
- PascalCase: `UserController`, `ProductService`, `OrderRepository`
- Module-based: `Handler`, `Service`, `Repository`

### Functions
- PascalCase for exported: `Create`, `GetByID`, `Update`
- camelCase for private: `toResponse`, `validateInput`

### Constants
- PascalCase with prefix: `UserStatusActive`, `ErrUserNotFound`

---

## Checklist for New Features

When adding a new feature, ensure:

- [ ] Choose architecture pattern (module or layer-based)
- [ ] Generate with `/module create` or `/layer create`
- [ ] Add validation tags to all DTOs
- [ ] Use response utility (no `fiber.Map`)
- [ ] Create constants for statuses and errors
- [ ] Use constants in code
- [ ] Implement dependency injection
- [ ] Add Swagger documentation
- [ ] Register routes
- [ ] Test all endpoints

---

## Skills Available

| Skill | Purpose |
|-------|---------|
| `/init-project <name>` | Initialize new project from scratch |
| `/module create <name>` | Create module-based feature |
| `/layer create <name>` | Create layer-based feature |

---

## Quick Reference

### Handler/Controller Template

```go
package user

import (
    "yourproject/core/response"
    "github.com/go-playground/validator/v10"
    "github.com/gofiber/fiber/v2"
)

var validate = validator.New()

type Handler struct {
    service *Service
}

func NewHandler(service *Service) *Handler {
    return &Handler{service: service}
}

func (h *Handler) Create(c *fiber.Ctx) error {
    var req CreateRequest

    if err := c.BodyParser(&req); err != nil {
        return response.BadRequest(c, "Invalid request body", nil)
    }

    if err := validate.Struct(&req); err != nil {
        return response.ValidationError(c, err.Error())
    }

    result, err := h.service.Create(&req)
    if err != nil {
        return response.InternalServerError(c, err.Error(), nil)
    }

    return response.Created(c, "User created successfully", result)
}

func (h *Handler) GetByID(c *fiber.Ctx) error {
    id, err := c.ParamsInt("id")
    if err != nil {
        return response.BadRequest(c, "Invalid ID", nil)
    }

    result, err := h.service.GetByID(uint(id))
    if err != nil {
        return response.NotFound(c, err.Error())
    }

    return response.OK(c, result)
}
```

---

**Follow these conventions for consistent, maintainable code!**
