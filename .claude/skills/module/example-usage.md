# Module Skill - Example Usage

## Quick Start Examples

### Example 1: Creating a User Module

**User Request:**
> "I need to create a user management module with CRUD operations"

**Your Action:**
```
/module create user
```

**What Gets Created:**
```
modules/user/
├── handler.go       # User HTTP handlers (Create, Get, Update, Delete, List)
├── service.go       # User business logic
├── repository.go    # User data access with GORM
├── model.go         # User entity with common fields
├── dto.go          # User DTOs (CreateRequest, UpdateRequest, Response)
└── routes.go       # User route registration
```

**Next Steps for User:**
1. Add specific fields to `model.go` (e.g., Name, Email, Password)
2. Update DTOs in `dto.go` to match model fields
3. Update service methods to handle business logic
4. Run migrations to create the users table
5. Generate Swagger docs: `make swag`
6. Test the endpoints

---

### Example 2: Creating a Product Module

**User Request:**
> "Create a product module"

**Your Action:**
```
/module create product
```

**Then Customize the Model:**

Edit `modules/product/model.go`:
```go
package product

import (
	"time"
	"gorm.io/gorm"
)

type Model struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	Name        string         `gorm:"size:255;not null" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	Price       float64        `gorm:"not null" json:"price"`
	Stock       int            `gorm:"default:0" json:"stock"`
	CategoryID  uint           `gorm:"not null" json:"category_id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Model) TableName() string {
	return "products"
}
```

Edit `modules/product/dto.go`:
```go
package product

import "time"

type CreateRequest struct {
	Name        string  `json:"name" validate:"required,min=1,max=255"`
	Description string  `json:"description"`
	Price       float64 `json:"price" validate:"required,gt=0"`
	Stock       int     `json:"stock" validate:"gte=0"`
	CategoryID  uint    `json:"category_id" validate:"required"`
}

type UpdateRequest struct {
	Name        *string  `json:"name,omitempty" validate:"omitempty,min=1,max=255"`
	Description *string  `json:"description,omitempty"`
	Price       *float64 `json:"price,omitempty" validate:"omitempty,gt=0"`
	Stock       *int     `json:"stock,omitempty" validate:"omitempty,gte=0"`
	CategoryID  *uint    `json:"category_id,omitempty"`
}

type Response struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Stock       int       `json:"stock"`
	CategoryID  uint      `json:"category_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ListResponse struct {
	Data  []*Response `json:"data"`
	Total int64       `json:"total"`
}
```

---

### Example 3: Module with Relations

**User Request:**
> "Create an order module that relates to users and products"

**Your Action:**
```
/module create order
```

**Customize with Relations:**

Edit `modules/order/model.go`:
```go
package order

import (
	"time"
	"gorm.io/gorm"
)

type Model struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	UserID    uint           `gorm:"not null;index" json:"user_id"`
	TotalAmount float64      `gorm:"not null" json:"total_amount"`
	Status    string         `gorm:"size:50;default:'pending'" json:"status"`
	Items     []OrderItem    `gorm:"foreignKey:OrderID" json:"items"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type OrderItem struct {
	ID        uint    `gorm:"primarykey" json:"id"`
	OrderID   uint    `gorm:"not null;index" json:"order_id"`
	ProductID uint    `gorm:"not null" json:"product_id"`
	Quantity  int     `gorm:"not null" json:"quantity"`
	Price     float64 `gorm:"not null" json:"price"`
}

func (Model) TableName() string {
	return "orders"
}

func (OrderItem) TableName() string {
	return "order_items"
}
```

---

### Example 4: Listing All Modules

**User Request:**
> "What modules do we have?"

**Your Action:**
```
/module list
```

**Expected Output:**
```
Modules in this project:
1. user - User management
2. product - Product catalog
3. order - Order processing
```

---

### Example 5: Showing Module Structure

**User Request:**
> "Show me the user module structure"

**Your Action:**
```
/module show user
```

**Expected Output:**
```
modules/user/
├── handler.go       - HTTP handlers with 5 endpoints
│   ├── Create(c *fiber.Ctx) error
│   ├── GetByID(c *fiber.Ctx) error
│   ├── List(c *fiber.Ctx) error
│   ├── Update(c *fiber.Ctx) error
│   └── Delete(c *fiber.Ctx) error
├── service.go       - Business logic layer
│   ├── Create(*CreateRequest) (*Response, error)
│   ├── GetByID(uint) (*Response, error)
│   ├── List() ([]*Response, error)
│   ├── Update(uint, *UpdateRequest) (*Response, error)
│   └── Delete(uint) error
├── repository.go    - Data access layer
│   ├── Create(*Model) error
│   ├── FindByID(uint) (*Model, error)
│   ├── FindAll() ([]*Model, error)
│   ├── Update(*Model) error
│   └── Delete(*Model) error
├── model.go        - Domain entity (User model)
├── dto.go          - Request/Response DTOs
└── routes.go       - Route registration
```

---

## Integration Example

### Updating Routes Setup

After creating modules, update `api/routes/setup.go`:

```go
package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	// Import your modules
	"go-fiber/modules/user"
	"go-fiber/modules/product"
	"go-fiber/modules/order"
)

func Setup(app *fiber.App, db *gorm.DB, rd *redis.Client) {
	// API v1 group
	v1 := app.Group("/api/v1")

	// Register module routes
	user.SetupRoutes(v1, db)
	product.SetupRoutes(v1, db)
	order.SetupRoutes(v1, db)
}
```

### Generated Routes

With the above setup, you'll have:

**User Module:**
- `POST   /api/v1/user/`      - Create user
- `GET    /api/v1/user/`      - List users
- `GET    /api/v1/user/:id`   - Get user by ID
- `PUT    /api/v1/user/:id`   - Update user
- `DELETE /api/v1/user/:id`   - Delete user

**Product Module:**
- `POST   /api/v1/product/`      - Create product
- `GET    /api/v1/product/`      - List products
- `GET    /api/v1/product/:id`   - Get product by ID
- `PUT    /api/v1/product/:id`   - Update product
- `DELETE /api/v1/product/:id`   - Delete product

**Order Module:**
- `POST   /api/v1/order/`      - Create order
- `GET    /api/v1/order/`      - List orders
- `GET    /api/v1/order/:id`   - Get order by ID
- `PUT    /api/v1/order/:id`   - Update order
- `DELETE /api/v1/order/:id`   - Delete order

---

## Testing Your Module

### Manual Testing with curl

```bash
# Create a user
curl -X POST http://localhost:3000/api/v1/user \
  -H "Content-Type: application/json" \
  -d '{"name": "John Doe", "email": "john@example.com"}'

# Get all users
curl http://localhost:3000/api/v1/user

# Get user by ID
curl http://localhost:3000/api/v1/user/1

# Update user
curl -X PUT http://localhost:3000/api/v1/user/1 \
  -H "Content-Type: application/json" \
  -d '{"name": "Jane Doe"}'

# Delete user
curl -X DELETE http://localhost:3000/api/v1/user/1
```

### Using Swagger UI

1. Generate docs: `make swag`
2. Run the app: `make run`
3. Open: `http://localhost:3000/swagger/index.html`
4. Test endpoints interactively

---

## Adding Middleware to a Module

Edit `modules/user/routes.go`:

```go
package user

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"go-fiber/api/middleware" // Your middleware package
)

func SetupRoutes(router fiber.Router, db *gorm.DB) {
	// Initialize dependencies
	repo := NewRepository(db)
	service := NewService(repo)
	handler := NewHandler(service)

	// Create module route group
	group := router.Group("/user")

	// Public routes
	group.Post("/register", handler.Register)
	group.Post("/login", handler.Login)

	// Protected routes with JWT middleware
	protected := group.Group("", middleware.JWTAuth())
	protected.Get("/", handler.List)
	protected.Get("/:id", handler.GetByID)
	protected.Put("/:id", handler.Update)
	protected.Delete("/:id", handler.Delete)
}
```

---

## Advanced: Module with Custom Methods

You can extend the generated code with custom business logic:

**In `modules/user/service.go`:**

```go
// Add custom method beyond CRUD
func (s *Service) ChangePassword(userID uint, oldPassword, newPassword string) error {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return err
	}

	// Verify old password
	if !verifyPassword(user.Password, oldPassword) {
		return errors.New("incorrect password")
	}

	// Hash new password
	hashedPassword, err := hashPassword(newPassword)
	if err != nil {
		return err
	}

	user.Password = hashedPassword
	return s.repo.Update(user)
}

func (s *Service) GetUserByEmail(email string) (*Response, error) {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	return toResponse(user), nil
}
```

**In `modules/user/repository.go`:**

```go
// Add custom query method
func (r *Repository) FindByEmail(email string) (*Model, error) {
	var user Model
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
```

**In `modules/user/handler.go`:**

```go
// Add custom endpoint
func (h *Handler) ChangePassword(c *fiber.Ctx) error {
	userID := c.Locals("userID").(uint) // From JWT middleware

	var req struct {
		OldPassword string `json:"old_password" validate:"required"`
		NewPassword string `json:"new_password" validate:"required,min=8"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request",
		})
	}

	if err := h.service.ChangePassword(userID, req.OldPassword, req.NewPassword); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Password changed successfully",
	})
}
```

**In `modules/user/routes.go`:**

```go
// Add the custom route
protected.Post("/change-password", handler.ChangePassword)
```

---

## Summary

The module skill helps you:

1. **Create** new modules with complete CRUD operations
2. **Organize** code by business domain
3. **Scale** your application with clear boundaries
4. **Maintain** code more easily with related code together
5. **Collaborate** with teams owning different modules

Use `/module create <name>` to get started, then customize the generated code to fit your specific needs!
