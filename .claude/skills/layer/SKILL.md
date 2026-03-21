---
name: layer
description: Create and manage layer-based project structure with separate controllers/, services/, and repositories/ folders organized by technical layer
disable-model-invocation: false
argument-hint: [resource-name] [options]
---

# Layer Skill

This skill helps you create and manage a **layer-based (package by layer)** project structure for Go Fiber applications, where code is organized by technical layer (controllers, services, repositories) rather than by business domain.

## Architecture Philosophy

**Layer-based structure** organizes code by technical responsibility:

```
api/
  controllers/
    user_controller.go
    product_controller.go
    order_controller.go
data/
  services/
    user_service.go
    product_service.go
    order_service.go
  repositories/
    user_repository.go
    product_repository.go
    order_repository.go
domain/
  models/
    user.go
    product.go
    order.go
  dto/
    user_dto.go
    product_dto.go
    order_dto.go
```

## When to Use Layer-Based

✅ **Good for:**
- Small to medium projects (3-10 features)
- Teams familiar with MVC/layered architecture
- Projects with consistent patterns across all features
- When all features have similar complexity
- Educational/learning projects
- When you want centralized view of all controllers/services

❌ **Consider module-based instead if:**
- Large projects (10+ features)
- Features have very different complexities
- Want stronger feature boundaries
- Teams work on separate domains
- Following Domain-Driven Design (DDD)

## Usage

1. **Create a new resource**: `/layer create <resource-name>`
2. **List all resources**: `/layer list`
3. **Show resource files**: `/layer show <resource-name>`

## Arguments

- `$0` - Action: `create`, `list`, `show`
- `$1` - Resource name (e.g., `user`, `product`, `order`)

## Implementation Steps

### For `create` action:

1. **Analyze existing structure**:
   - Check if `api/controllers/`, `data/services/`, `data/repositories/` exist
   - Check if `domain/models/`, `domain/dto/` exist
   - Create directories if needed

2. **Generate controller** in `api/controllers/<resource>_controller.go`:
   - Create controller struct that depends on service
   - Implement HTTP handler methods (Create, Get, Update, Delete, List)
   - Use Fiber context for request/response
   - Include proper error handling and validation
   - Add Swagger annotations
   - Follow naming: `UserController`, `ProductController`

3. **Generate service** in `data/services/<resource>_service.go`:
   - Create service struct that depends on repository
   - Implement business logic methods
   - Handle business rules and validations
   - Use dependency injection pattern
   - Follow naming: `UserService`, `ProductService`

4. **Generate repository** in `data/repositories/<resource>_repository.go`:
   - Create repository struct that depends on database
   - Implement data access methods using GORM
   - Include Create, FindByID, Update, Delete, List operations
   - Use proper error handling
   - Follow naming: `UserRepository`, `ProductRepository`

5. **Generate model** in `domain/models/<resource>.go`:
   - Create the domain entity struct with GORM tags
   - Include common fields (ID, CreatedAt, UpdatedAt, DeletedAt)
   - Add JSON tags for serialization
   - Add validation tags
   - Follow naming: `User`, `Product`

6. **Generate DTOs** in `domain/dto/<resource>_dto.go`:
   - Create request DTOs (Create<Resource>Request, Update<Resource>Request)
   - Create response DTOs (<Resource>Response, <Resource>ListResponse)
   - Add validation tags
   - Add JSON tags

7. **Generate constants** in `domain/constants/<resource>_constants.go`:
   - Create status enums (e.g., StatusActive, StatusInactive)
   - Create error message constants (e.g., ErrNotFound, ErrInvalidInput)
   - Create validation constants (min/max values, patterns)
   - Export all constants for use across the application

8. **Update routes** in `api/routes/setup.go`:
   - Import new controller
   - Register routes under resource prefix
   - Wire up dependencies (repository -> service -> controller)

## Code Templates

### Controller Template (api/controllers/<resource>_controller.go):
```go
package controllers

import (
	"go-fiber/core/response"
	"go-fiber/data/services"
	"go-fiber/domain/dto"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate = validator.New()

type <Resource>Controller struct {
	service *services.<Resource>Service
}

func New<Resource>Controller(service *services.<Resource>Service) *<Resource>Controller {
	return &<Resource>Controller{service: service}
}

// Create<Resource> godoc
// @Summary Create <Resource>
// @Description Create a new <resource>
// @Tags <resource>
// @Accept json
// @Produce json
// @Param request body dto.Create<Resource>Request true "Create Request"
// @Success 201 {object} dto.<Resource>Response
// @Failure 400 {object} map[string]interface{}
// @Router /<resource> [post]
func (ctrl *<Resource>Controller) Create(c *fiber.Ctx) error {
	var req dto.Create<Resource>Request

	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body", nil)
	}

	// Validate request
	if err := validate.Struct(&req); err != nil {
		return response.ValidationError(c, err.Error())
	}

	result, err := ctrl.service.Create(&req)
	if err != nil {
		return response.InternalServerError(c, err.Error(), nil)
	}

	return response.Created(c, "<Resource> created successfully", result)
}

// Get<Resource>ByID godoc
// @Summary Get <Resource>
// @Description Get <resource> by ID
// @Tags <resource>
// @Produce json
// @Param id path int true "<Resource> ID"
// @Success 200 {object} dto.<Resource>Response
// @Failure 404 {object} map[string]interface{}
// @Router /<resource>/{id} [get]
func (ctrl *<Resource>Controller) GetByID(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return response.BadRequest(c, "Invalid ID", nil)
	}

	result, err := ctrl.service.GetByID(uint(id))
	if err != nil {
		return response.NotFound(c, err.Error())
	}

	return response.OK(c, result)
}

// List<Resource>s godoc
// @Summary List <Resource>s
// @Description Get all <resource>s
// @Tags <resource>
// @Produce json
// @Success 200 {object} dto.<Resource>ListResponse
// @Router /<resource> [get]
func (ctrl *<Resource>Controller) List(c *fiber.Ctx) error {
	results, err := ctrl.service.List()
	if err != nil {
		return response.InternalServerError(c, err.Error(), nil)
	}

	return response.OK(c, dto.<Resource>ListResponse{Data: results})
}

// Update<Resource> godoc
// @Summary Update <Resource>
// @Description Update <resource> by ID
// @Tags <resource>
// @Accept json
// @Produce json
// @Param id path int true "<Resource> ID"
// @Param request body dto.Update<Resource>Request true "Update Request"
// @Success 200 {object} dto.<Resource>Response
// @Failure 400 {object} map[string]interface{}
// @Router /<resource>/{id} [put]
func (ctrl *<Resource>Controller) Update(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return response.BadRequest(c, "Invalid ID", nil)
	}

	var req dto.Update<Resource>Request
	if err := c.BodyParser(&req); err != nil {
		return response.BadRequest(c, "Invalid request body", nil)
	}

	// Validate request
	if err := validate.Struct(&req); err != nil {
		return response.ValidationError(c, err.Error())
	}

	result, err := ctrl.service.Update(uint(id), &req)
	if err != nil {
		return response.InternalServerError(c, err.Error(), nil)
	}

	return response.OK(c, result)
}

// Delete<Resource> godoc
// @Summary Delete <Resource>
// @Description Delete <resource> by ID
// @Tags <resource>
// @Param id path int true "<Resource> ID"
// @Success 204
// @Failure 400 {object} map[string]interface{}
// @Router /<resource>/{id} [delete]
func (ctrl *<Resource>Controller) Delete(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return response.BadRequest(c, "Invalid ID", nil)
	}

	if err := ctrl.service.Delete(uint(id)); err != nil {
		return response.InternalServerError(c, err.Error(), nil)
	}

	return response.NoContent(c)
}
```

### Service Template (data/services/<resource>_service.go):
```go
package services

import (
	"errors"
	"go-fiber/data/repositories"
	"go-fiber/domain/constants"
	"go-fiber/domain/dto"
	"go-fiber/domain/models"
)

type <Resource>Service struct {
	repo *repositories.<Resource>Repository
}

func New<Resource>Service(repo *repositories.<Resource>Repository) *<Resource>Service {
	return &<Resource>Service{repo: repo}
}

func (s *<Resource>Service) Create(req *dto.Create<Resource>Request) (*dto.<Resource>Response, error) {
	model := &models.<Resource>{
		// Map request fields to model
	}

	if err := s.repo.Create(model); err != nil {
		return nil, err
	}

	return to<Resource>Response(model), nil
}

func (s *<Resource>Service) GetByID(id uint) (*dto.<Resource>Response, error) {
	model, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New(constants.Err<Resource>NotFound)
	}

	return to<Resource>Response(model), nil
}

func (s *<Resource>Service) List() ([]*dto.<Resource>Response, error) {
	models, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	responses := make([]*dto.<Resource>Response, len(models))
	for i, model := range models {
		responses[i] = to<Resource>Response(model)
	}

	return responses, nil
}

func (s *<Resource>Service) Update(id uint, req *dto.Update<Resource>Request) (*dto.<Resource>Response, error) {
	model, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	// Update model fields from request

	if err := s.repo.Update(model); err != nil {
		return nil, err
	}

	return to<Resource>Response(model), nil
}

func (s *<Resource>Service) Delete(id uint) error {
	model, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	return s.repo.Delete(model)
}

func to<Resource>Response(model *models.<Resource>) *dto.<Resource>Response {
	return &dto.<Resource>Response{
		ID:        model.ID,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}
}
```

### Repository Template (data/repositories/<resource>_repository.go):
```go
package repositories

import (
	"gorm.io/gorm"
	"go-fiber/domain/models"
)

type <Resource>Repository struct {
	db *gorm.DB
}

func New<Resource>Repository(db *gorm.DB) *<Resource>Repository {
	return &<Resource>Repository{db: db}
}

func (r *<Resource>Repository) Create(model *models.<Resource>) error {
	return r.db.Create(model).Error
}

func (r *<Resource>Repository) FindByID(id uint) (*models.<Resource>, error) {
	var model models.<Resource>
	if err := r.db.First(&model, id).Error; err != nil {
		return nil, err
	}
	return &model, nil
}

func (r *<Resource>Repository) FindAll() ([]*models.<Resource>, error) {
	var models []*models.<Resource>
	if err := r.db.Find(&models).Error; err != nil {
		return nil, err
	}
	return models, nil
}

func (r *<Resource>Repository) Update(model *models.<Resource>) error {
	return r.db.Save(model).Error
}

func (r *<Resource>Repository) Delete(model *models.<Resource>) error {
	return r.db.Delete(model).Error
}
```

### Model Template (domain/models/<resource>.go):
```go
package models

import (
	"time"
	"gorm.io/gorm"
)

type <Resource> struct {
	ID        uint           ` + "`gorm:\"primarykey\" json:\"id\"`" + `
	Status    string         ` + "`gorm:\"size:50;default:'active'\" json:\"status\"`" + `
	CreatedAt time.Time      ` + "`json:\"created_at\"`" + `
	UpdatedAt time.Time      ` + "`json:\"updated_at\"`" + `
	DeletedAt gorm.DeletedAt ` + "`gorm:\"index\" json:\"-\"`" + `

	// Add your model fields here
	// Use constants from domain/constants/<resource>_constants.go for status values
}

func (<Resource>) TableName() string {
	return "<resource>s"
}
```

### DTO Template (domain/dto/<resource>_dto.go):
```go
package dto

import "time"

// Create<Resource>Request - ALWAYS include validation tags
type Create<Resource>Request struct {
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

// Update<Resource>Request - Use omitempty for optional updates, but ALWAYS validate when provided
type Update<Resource>Request struct {
	// Add your request fields here with validation tags
	// Use pointers for optional fields
	// Examples:
	// Name        *string  ` + "`json:\"name,omitempty\" validate:\"omitempty,min=3,max=255\"`" + `
	// Email       *string  ` + "`json:\"email,omitempty\" validate:\"omitempty,email\"`" + `
	// Price       *float64 ` + "`json:\"price,omitempty\" validate:\"omitempty,gt=0\"`" + `
	// Age         *int     ` + "`json:\"age,omitempty\" validate:\"omitempty,gte=0,lte=150\"`" + `
	// Status      *string  ` + "`json:\"status,omitempty\" validate:\"omitempty,oneof=active inactive pending\"`" + `
}

type <Resource>Response struct {
	ID        uint      ` + "`json:\"id\"`" + `
	Status    string    ` + "`json:\"status\"`" + `
	CreatedAt time.Time ` + "`json:\"created_at\"`" + `
	UpdatedAt time.Time ` + "`json:\"updated_at\"`" + `

	// Add your response fields here (no validation needed for responses)
}

type <Resource>ListResponse struct {
	Data  []*<Resource>Response ` + "`json:\"data\"`" + `
	Total int64                 ` + "`json:\"total\"`" + `
	Page  int                   ` + "`json:\"page,omitempty\"`" + `
	Limit int                   ` + "`json:\"limit,omitempty\"`" + `
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

### Constants Template (domain/constants/<resource>_constants.go):
```go
package constants

// <Resource> status constants
const (
	<Resource>StatusActive   = "active"
	<Resource>StatusInactive = "inactive"
	<Resource>StatusPending  = "pending"
)

// <Resource> error messages
const (
	Err<Resource>NotFound     = "<resource> not found"
	Err<Resource>InvalidInput = "invalid <resource> input"
	Err<Resource>AlreadyExists = "<resource> already exists"
	Err<Resource>Unauthorized = "unauthorized access to <resource>"
)

// <Resource> validation constants
const (
	<Resource>NameMinLength = 3
	<Resource>NameMaxLength = 255
	<Resource>DescMaxLength = 1000
)

// Get all valid statuses
func Valid<Resource>Statuses() []string {
	return []string{
		<Resource>StatusActive,
		<Resource>StatusInactive,
		<Resource>StatusPending,
	}
}

// Check if status is valid
func IsValid<Resource>Status(status string) bool {
	validStatuses := Valid<Resource>Statuses()
	for _, s := range validStatuses {
		if s == status {
			return true
		}
	}
	return false
}
```

### Routes Registration (api/routes/setup.go):
```go
package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"go-fiber/api/controllers"
	"go-fiber/data/repositories"
	"go-fiber/data/services"
)

func Setup(app *fiber.App, db *gorm.DB, rd *redis.Client) {
	v1 := app.Group("/api/v1")

	// User routes
	userRepo := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepo)
	userController := controllers.NewUserController(userService)

	userGroup := v1.Group("/user")
	userGroup.Post("/", userController.Create)
	userGroup.Get("/", userController.List)
	userGroup.Get("/:id", userController.GetByID)
	userGroup.Put("/:id", userController.Update)
	userGroup.Delete("/:id", userController.Delete)
}
```

## Example Usage

User: "Create a user resource"
→ `/layer create user`

Creates:
- `api/controllers/user_controller.go`
- `data/services/user_service.go`
- `data/repositories/user_repository.go`
- `domain/models/user.go`
- `domain/dto/user_dto.go`
- `domain/constants/user_constants.go` - Status enums, error messages, validation constants
- Updates `api/routes/setup.go`

User: "Show all resources"
→ `/layer list`

User: "Show user files"
→ `/layer show user`

## Important Notes

1. **Centralized view** - All controllers, services, and repositories are in their own folders
2. **Consistent patterns** - Same structure for every resource
3. **Easy to find** - Know exactly where each layer's code lives
4. **Dependency injection** - Repository -> Service -> Controller
5. **Follow Go conventions** - Exported types start with capital letters
6. **Proper error handling** - Return errors, don't panic
7. **Swagger documentation** - Document APIs as you build

## Migration Between Styles

### From Layer to Module:
If your project grows complex, you can migrate:
```bash
/module create user  # Creates modules/user/
# Move code from controllers/user_controller.go → modules/user/handler.go
# Move code from services/user_service.go → modules/user/service.go
# etc.
```

### From Module to Layer:
If you prefer centralized structure:
```bash
/layer create user  # Creates layer-based files
# Move code from modules/user/handler.go → controllers/user_controller.go
# etc.
```

## Benefits of Layer-Based

✅ **Pros:**
- Simple, familiar structure
- Easy to see all controllers/services at once
- Consistent patterns across resources
- Good for small-medium projects
- Easy to learn and onboard new developers
- Clear technical separation

❌ **Cons:**
- Related code scattered across directories
- Harder to navigate large codebases
- Feature boundaries less clear
- Can become unwieldy with many resources

## When to Use This Skill

Use `/layer` when:
- Building small to medium projects
- Your team prefers traditional MVC/layered architecture
- All features have similar complexity
- You want a centralized view of all controllers/services
- Following common Go patterns like those in many tutorials

Use `/module` instead when:
- Building large, complex applications
- Following Domain-Driven Design
- Features have varying complexity
- Want stronger feature boundaries
- Teams own different domains
