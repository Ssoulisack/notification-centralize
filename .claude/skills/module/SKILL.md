---
name: module
description: Create module-based project structure organizing code by feature/domain with each module containing its own handler, service, repository, model, and routes
disable-model-invocation: false
argument-hint: [module-name] [options]
---

# Module Skill

This skill creates and manages a **module-based (package by feature)** project structure for Go Fiber applications. Code is organized by business domain/feature, with each module being self-contained.

## Module-Based Architecture

Organize code by **business domain** where each module contains all its layers:

```
modules/
  user/
    handler.go       # User HTTP handlers
    service.go       # User business logic
    repository.go    # User data access
    model.go         # User entity
    dto.go          # User DTOs
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

## Why Module-Based?

✅ **Benefits:**
- All code for a feature lives together
- Strong feature boundaries
- Easy to navigate - everything in one folder
- Independent modules with minimal coupling
- Scales well as project grows
- Teams can own specific modules
- Clear business domain boundaries

## Usage

When the user requests to create a new module or feature, use this skill to:

1. **Create a new module**: `/module create <module-name>`
2. **Add a resource to existing module**: `/module add <module-name> <resource-name>`
3. **List all modules**: `/module list`
4. **Show module structure**: `/module show <module-name>`

## Arguments

- `$0` - Action: `create`, `add`, `list`, `show`
- `$1` - Module name (e.g., `user`, `product`, `order`)
- `$2` - Resource name (optional, for `add` action)

## Implementation Steps

### For `create` action:

1. **Create module directory structure**:
   ```
   modules/<module-name>/
   ├── handler.go       # HTTP handlers (controllers)
   ├── service.go       # Business logic
   ├── repository.go    # Database operations
   ├── model.go         # Domain models/entities
   ├── dto.go          # Data Transfer Objects
   └── routes.go       # Module-specific routes
   ```

2. **Generate handler.go**:
   - Create a handler struct that depends on the service
   - Implement HTTP handler methods (Create, Get, Update, Delete, List)
   - Use Fiber context for request/response
   - Include proper error handling and validation
   - Add Swagger annotations

3. **Generate service.go**:
   - Create a service struct that depends on the repository
   - Implement business logic methods
   - Handle business rules and validations
   - Use dependency injection pattern

4. **Generate repository.go**:
   - Create a repository struct that depends on database connection
   - Implement data access methods using GORM
   - Include Create, FindByID, Update, Delete, List operations
   - Use proper error handling

5. **Generate model.go**:
   - Create the domain entity struct with GORM tags
   - Include common fields (ID, CreatedAt, UpdatedAt, DeletedAt)
   - Add JSON tags for serialization
   - Add validation tags

6. **Generate dto.go**:
   - Create request DTOs (CreateRequest, UpdateRequest)
   - Create response DTOs (Response, ListResponse)
   - Add validation tags
   - Add JSON tags

7. **Generate constants file** in `domain/constants/<module>_constants.go`:
   - Create status enums (e.g., StatusActive, StatusInactive)
   - Create error message constants (e.g., ErrNotFound, ErrInvalidInput)
   - Create validation constants (min/max values, patterns)
   - Export all constants for use across the application

8. **Generate routes.go**:
   - Create a function to register module routes
   - Group routes under module prefix (e.g., `/api/v1/<module>`)
   - Apply middleware as needed
   - Return the route group

9. **Update main routes file**:
   - Import the new module
   - Register module routes in `api/routes/setup.go`

10. **Create a module initialization function**:
   - Wire up dependencies (repository -> service -> handler)
   - Return configured handler ready for route registration

### For `add` action:

1. Check if module exists in `modules/<module-name>/`
2. Add new handler methods for the resource
3. Add new service methods
4. Add new repository methods if needed
5. Update routes.go with new routes

### For `list` action:

1. List all directories in `modules/` folder
2. Show module structure for each

### For `show` action:

1. Display the structure of the specified module
2. List all files and their purposes
3. Show exported functions/methods

## Code Templates

### Handler Template (handler.go):
```go
package <module>

import (
	"<project>/core/response"

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

// Create godoc
// @Summary Create <Module>
// @Description Create a new <module>
// @Tags <module>
// @Accept json
// @Produce json
// @Param request body CreateRequest true "Create Request"
// @Success 201 {object} Response
// @Failure 400 {object} map[string]interface{}
// @Router /<module> [post]
func (h *Handler) Create(c *fiber.Ctx) error {
	var req CreateRequest

	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body", nil)
	}

	// Validate request
	if err := validate.Struct(&req); err != nil {
		return response.ValidationError(c, err.Error())
	}

	result, err := h.service.Create(&req)
	if err != nil {
		return response.InternalServerError(c, err.Error(), nil)
	}

	return response.Created(c, "<Module> created successfully", result)
}

// GetByID godoc
// @Summary Get <Module>
// @Description Get <module> by ID
// @Tags <module>
// @Produce json
// @Param id path int true "<Module> ID"
// @Success 200 {object} Response
// @Failure 404 {object} map[string]interface{}
// @Router /<module>/{id} [get]
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

// List godoc
// @Summary List <Module>s
// @Description Get all <module>s
// @Tags <module>
// @Produce json
// @Success 200 {object} ListResponse
// @Router /<module> [get]
func (h *Handler) List(c *fiber.Ctx) error {
	results, err := h.service.List()
	if err != nil {
		return response.InternalServerError(c, err.Error(), nil)
	}

	return response.OK(c, ListResponse{Data: results})
}

// Update godoc
// @Summary Update <Module>
// @Description Update <module> by ID
// @Tags <module>
// @Accept json
// @Produce json
// @Param id path int true "<Module> ID"
// @Param request body UpdateRequest true "Update Request"
// @Success 200 {object} Response
// @Failure 400 {object} map[string]interface{}
// @Router /<module>/{id} [put]
func (h *Handler) Update(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return response.BadRequest(c, "Invalid ID", nil)
	}

	var req UpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body", nil)
	}

	// Validate request
	if err := validate.Struct(&req); err != nil {
		return response.ValidationError(c, err.Error())
	}

	result, err := h.service.Update(uint(id), &req)
	if err != nil {
		return response.InternalServerError(c, err.Error(), nil)
	}

	return response.OK(c, result)
}

// Delete godoc
// @Summary Delete <Module>
// @Description Delete <module> by ID
// @Tags <module>
// @Param id path int true "<Module> ID"
// @Success 204
// @Failure 400 {object} map[string]interface{}
// @Router /<module>/{id} [delete]
func (h *Handler) Delete(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return response.BadRequest(c, "Invalid ID", nil)
	}

	if err := h.service.Delete(uint(id)); err != nil {
		return response.InternalServerError(c, err.Error(), nil)
	}

	return response.NoContent(c)
}
```

### Service Template (service.go):
```go
package <module>

import (
	"errors"
	"<project>/domain/constants"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(req *CreateRequest) (*Response, error) {
	model := &Model{
		// Map request fields to model
	}

	if err := s.repo.Create(model); err != nil {
		return nil, err
	}

	return toResponse(model), nil
}

func (s *Service) GetByID(id uint) (*Response, error) {
	model, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New(constants.Err<Module>NotFound)
	}

	return toResponse(model), nil
}

func (s *Service) List() ([]*Response, error) {
	models, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	responses := make([]*Response, len(models))
	for i, model := range models {
		responses[i] = toResponse(model)
	}

	return responses, nil
}

func (s *Service) Update(id uint, req *UpdateRequest) (*Response, error) {
	model, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Update model fields from request

	if err := s.repo.Update(model); err != nil {
		return nil, err
	}

	return toResponse(model), nil
}

func (s *Service) Delete(id uint) error {
	model, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	return s.repo.Delete(model)
}

func toResponse(model *Model) *Response {
	return &Response{
		ID:        model.ID,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}
}
```

### Repository Template (repository.go):
```go
package <module>

import (
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(model *Model) error {
	return r.db.Create(model).Error
}

func (r *Repository) FindByID(id uint) (*Model, error) {
	var model Model
	if err := r.db.First(&model, id).Error; err != nil {
		return nil, err
	}
	return &model, nil
}

func (r *Repository) FindAll() ([]*Model, error) {
	var models []*Model
	if err := r.db.Find(&models).Error; err != nil {
		return nil, err
	}
	return models, nil
}

func (r *Repository) Update(model *Model) error {
	return r.db.Save(model).Error
}

func (r *Repository) Delete(model *Model) error {
	return r.db.Delete(model).Error
}
```

### Model Template (model.go):
```go
package <module>

import (
	"time"
	"gorm.io/gorm"
)

type Model struct {
	ID        uint           ` + "`gorm:\"primarykey\" json:\"id\"`" + `
	Status    string         ` + "`gorm:\"size:50;default:'active'\" json:\"status\"`" + `
	CreatedAt time.Time      ` + "`json:\"created_at\"`" + `
	UpdatedAt time.Time      ` + "`json:\"updated_at\"`" + `
	DeletedAt gorm.DeletedAt ` + "`gorm:\"index\" json:\"-\"`" + `

	// Add your model fields here
	// Use constants from domain/constants/<module>_constants.go for status values
}

func (Model) TableName() string {
	return "<module>s"
}
```

### DTO Template (dto.go):
```go
package <module>

import "time"

// CreateRequest - ALWAYS include validation tags
type CreateRequest struct {
	// Add your request fields here with validation tags
	// Examples:
	// Name        string  ` + "`json:\"name\" validate:\"required,min=3,max=255\"`" + `
	// Email       string  ` + "`json:\"email\" validate:\"required,email\"`" + `
	// Price       float64 ` + "`json:\"price\" validate:\"required,gt=0\"`" + `
	// Age         int     ` + "`json:\"age\" validate:\"required,gte=0,lte=150\"`" + `
	// Status      string  ` + "`json:\"status\" validate:\"required,oneof=active inactive pending\"`" + `
	// Website     string  ` + "`json:\"website\" validate:\"omitempty,url\"`" + `
	// Phone       string  ` + "`json:\"phone\" validate:\"omitempty,e164\"`" + `
}

// UpdateRequest - Use omitempty for optional updates, but ALWAYS validate when provided
type UpdateRequest struct {
	// Add your request fields here with validation tags
	// Use pointers for optional fields
	// Examples:
	// Name        *string  ` + "`json:\"name,omitempty\" validate:\"omitempty,min=3,max=255\"`" + `
	// Email       *string  ` + "`json:\"email,omitempty\" validate:\"omitempty,email\"`" + `
	// Price       *float64 ` + "`json:\"price,omitempty\" validate:\"omitempty,gt=0\"`" + `
	// Age         *int     ` + "`json:\"age,omitempty\" validate:\"omitempty,gte=0,lte=150\"`" + `
	// Status      *string  ` + "`json:\"status,omitempty\" validate:\"omitempty,oneof=active inactive pending\"`" + `
}

type Response struct {
	ID        uint      ` + "`json:\"id\"`" + `
	Status    string    ` + "`json:\"status\"`" + `
	CreatedAt time.Time ` + "`json:\"created_at\"`" + `
	UpdatedAt time.Time ` + "`json:\"updated_at\"`" + `

	// Add your response fields here (no validation needed for responses)
}

type ListResponse struct {
	Data  []*Response ` + "`json:\"data\"`" + `
	Total int64       ` + "`json:\"total\"`" + `
	Page  int         ` + "`json:\"page,omitempty\"`" + `
	Limit int         ` + "`json:\"limit,omitempty\"`" + `
}

// Common validation tags reference:
// required          - Field is mandatory
// omitempty         - Validate only if provided
// min=N, max=N      - String length or numeric value bounds
// gte=N, lte=N      - Greater/less than or equal (numeric)
// gt=N, lt=N        - Greater/less than (numeric)
// email             - Valid email format
// url               - Valid URL format
// e164              - Valid phone number (E.164 format)
// oneof=a b c       - Must be one of the listed values
// len=N             - Exact length
// alpha             - Alphabetic characters only
// alphanum          - Alphanumeric characters only
// numeric           - Numeric string
// uuid              - Valid UUID
// datetime=layout   - Valid datetime with specific layout
```

### Constants Template (domain/constants/<module>_constants.go):
```go
package constants

// <Module> status constants
const (
	<Module>StatusActive   = "active"
	<Module>StatusInactive = "inactive"
	<Module>StatusPending  = "pending"
)

// <Module> error messages
const (
	Err<Module>NotFound     = "<module> not found"
	Err<Module>InvalidInput = "invalid <module> input"
	Err<Module>AlreadyExists = "<module> already exists"
	Err<Module>Unauthorized = "unauthorized access to <module>"
)

// <Module> validation constants
const (
	<Module>NameMinLength = 3
	<Module>NameMaxLength = 255
	<Module>DescMaxLength = 1000
)

// Get all valid statuses
func Valid<Module>Statuses() []string {
	return []string{
		<Module>StatusActive,
		<Module>StatusInactive,
		<Module>StatusPending,
	}
}

// Check if status is valid
func IsValid<Module>Status(status string) bool {
	validStatuses := Valid<Module>Statuses()
	for _, s := range validStatuses {
		if s == status {
			return true
		}
	}
	return false
}
```

### Routes Template (routes.go):
```go
package <module>

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupRoutes(router fiber.Router, db *gorm.DB) {
	// Initialize dependencies
	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service)

	// Create module route group
	group := router.Group("/<module>")

	// Register routes
	group.Post("/", handler.Create)
	group.Get("/", handler.List)
	group.Get("/:id", handler.GetByID)
	group.Put("/:id", handler.Update)
	group.Delete("/:id", handler.Delete)
}
```

## Example Usage

User: "Create a user module"
→ `/module create user`

User: "Add authentication to the user module"
→ `/module add user auth`

User: "Show me all modules"
→ `/module list`

## Important Notes

1. **Always use dependency injection** - Repository -> Service -> Handler
2. **Each module is self-contained** - All related code stays together
3. **Follow Go naming conventions** - Exported types start with capital letters
4. **Use proper error handling** - Return errors, don't panic
5. **Add Swagger annotations** - Document your API as you build
6. **Keep modules independent** - Minimize cross-module dependencies
7. **Use interfaces when needed** - For testing and flexibility

## Migration Guide

If converting from layer-based to module-based:

1. Create `modules/` directory
2. For each feature/domain:
   - Create module directory
   - Move related handlers, services, repositories into the module
   - Create routes.go for module-specific routing
3. Update `api/routes/setup.go` to import and register module routes
4. Update imports throughout the codebase

## When to Use Module-Based

✅ **Use `/module` for:**
- Medium to large projects
- Domain-Driven Design (DDD)
- Features with varying complexity
- When you want strong feature boundaries
- Teams working on different domains
- Long-term maintainability

## When to Use This Skill

Use this skill when the user wants to:
- Create a new feature or domain module
- Add CRUD operations for a new entity
- Organize code by business domain
- Follow Domain-Driven Design principles
- Keep all related code together in one place

The skill will automatically detect the action based on the user's request and execute the appropriate steps.
