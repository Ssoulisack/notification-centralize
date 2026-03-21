# Validation Guide

All DTOs **MUST** include validation tags using the `go-playground/validator` package.

## Important: Use Response Utility

**Always use the response utility (`core/response`) instead of `fiber.Map`** for consistent API responses.

✅ **Correct:**
```go
return response.ValidationError(c, err.Error())
return response.BadRequest(c, "Invalid request body", nil)
```

❌ **Don't do this:**
```go
return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
    "error": "Validation failed",
})
```

## Installation

```bash
go get github.com/go-playground/validator/v10
```

## Basic Usage

### In Handler/Controller

```go
import (
    "yourproject/core/response"
    "github.com/go-playground/validator/v10"
    "github.com/gofiber/fiber/v2"
)

var validate = validator.New()

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

    // Continue with business logic
    result, err := h.service.Create(&req)
    if err != nil {
        return response.InternalServerError(c, err.Error(), nil)
    }

    return response.Created(c, "Resource created successfully", result)
}
```

## Validation Tags Reference

### Required & Optional

```go
type Request struct {
    Name string `json:"name" validate:"required"`           // Mandatory
    Bio  string `json:"bio" validate:"omitempty"`          // Optional, but validate if provided
    Age  *int   `json:"age,omitempty" validate:"omitempty,gte=0"` // Optional for updates
}
```

### String Validation

```go
type Request struct {
    // Length validation
    Username string `json:"username" validate:"required,min=3,max=20"`
    Bio      string `json:"bio" validate:"omitempty,max=500"`
    Code     string `json:"code" validate:"required,len=6"`

    // Format validation
    Email    string `json:"email" validate:"required,email"`
    Website  string `json:"website" validate:"omitempty,url"`
    Phone    string `json:"phone" validate:"omitempty,e164"`

    // Character types
    Username string `json:"username" validate:"required,alphanum"`
    Code     string `json:"code" validate:"required,numeric"`
    Name     string `json:"name" validate:"required,alpha"`

    // UUID
    ID       string `json:"id" validate:"required,uuid"`
}
```

### Numeric Validation

```go
type Request struct {
    // Greater than / Less than
    Age      int     `json:"age" validate:"required,gte=0,lte=150"`
    Price    float64 `json:"price" validate:"required,gt=0"`
    Quantity int     `json:"quantity" validate:"required,min=1"`

    // Exact value
    Rating   int     `json:"rating" validate:"required,oneof=1 2 3 4 5"`
}
```

### Enum Validation

```go
type Request struct {
    // Must be one of the listed values
    Status string `json:"status" validate:"required,oneof=active inactive pending"`
    Role   string `json:"role" validate:"required,oneof=admin user guest"`
    Type   string `json:"type" validate:"required,oneof=personal business enterprise"`
}
```

### Date/Time Validation

```go
type Request struct {
    // RFC3339 format: 2006-01-02T15:04:05Z07:00
    CreatedAt time.Time `json:"created_at" validate:"required"`

    // Custom datetime format
    BirthDate string `json:"birth_date" validate:"required,datetime=2006-01-02"`
}
```

### Array/Slice Validation

```go
type Request struct {
    // Array length
    Tags []string `json:"tags" validate:"required,min=1,max=10"`
    IDs  []int    `json:"ids" validate:"required,dive,gte=1"` // Each element >= 1

    // Nested struct validation
    Items []Item `json:"items" validate:"required,dive"`
}

type Item struct {
    Name  string  `json:"name" validate:"required"`
    Price float64 `json:"price" validate:"required,gt=0"`
}
```

### Complex Validation

```go
type Request struct {
    // Multiple conditions
    Password string `json:"password" validate:"required,min=8,max=128,containsany=!@#$%"`

    // Conditional validation (eqfield, nefield, gtfield, ltfield)
    Password        string `json:"password" validate:"required,min=8"`
    PasswordConfirm string `json:"password_confirm" validate:"required,eqfield=Password"`

    // Required with/without
    Email string `json:"email" validate:"required_with=Phone"` // Required if Phone is provided
    Phone string `json:"phone"`
}
```

## Common Validation Patterns

### Create Request

```go
type CreateProductRequest struct {
    Name        string  `json:"name" validate:"required,min=3,max=255"`
    Description string  `json:"description" validate:"omitempty,max=1000"`
    Price       float64 `json:"price" validate:"required,gt=0"`
    Stock       int     `json:"stock" validate:"required,gte=0"`
    CategoryID  uint    `json:"category_id" validate:"required,gte=1"`
    Tags        []string `json:"tags" validate:"omitempty,min=1,max=10,dive,min=2,max=50"`
    Status      string  `json:"status" validate:"required,oneof=active inactive draft"`
}
```

### Update Request (Optional Fields)

```go
type UpdateProductRequest struct {
    Name        *string  `json:"name,omitempty" validate:"omitempty,min=3,max=255"`
    Description *string  `json:"description,omitempty" validate:"omitempty,max=1000"`
    Price       *float64 `json:"price,omitempty" validate:"omitempty,gt=0"`
    Stock       *int     `json:"stock,omitempty" validate:"omitempty,gte=0"`
    CategoryID  *uint    `json:"category_id,omitempty" validate:"omitempty,gte=1"`
    Status      *string  `json:"status,omitempty" validate:"omitempty,oneof=active inactive draft"`
}
```

### User Registration

```go
type RegisterRequest struct {
    Email           string `json:"email" validate:"required,email"`
    Password        string `json:"password" validate:"required,min=8,max=128"`
    PasswordConfirm string `json:"password_confirm" validate:"required,eqfield=Password"`
    FirstName       string `json:"first_name" validate:"required,min=2,max=50,alpha"`
    LastName        string `json:"last_name" validate:"required,min=2,max=50,alpha"`
    Phone           string `json:"phone" validate:"omitempty,e164"`
    DateOfBirth     string `json:"date_of_birth" validate:"required,datetime=2006-01-02"`
    Terms           bool   `json:"terms" validate:"required,eq=true"`
}
```

## Custom Error Messages

For better user experience, format validation errors:

```go
import (
    "fmt"
    "yourproject/core/response"
    "github.com/go-playground/validator/v10"
    "github.com/gofiber/fiber/v2"
)

func (h *Handler) Create(c *fiber.Ctx) error {
    var req CreateRequest

    // Parse request body
    if err := c.BodyParser(&req); err != nil {
        return response.BadRequest(c, "Invalid request body", nil)
    }

    // Validate request
    if err := validate.Struct(&req); err != nil {
        // Format validation errors for better UX
        errors := make(map[string]string)
        for _, err := range err.(validator.ValidationErrors) {
            field := err.Field()
            switch err.Tag() {
            case "required":
                errors[field] = fmt.Sprintf("%s is required", field)
            case "email":
                errors[field] = fmt.Sprintf("%s must be a valid email", field)
            case "min":
                errors[field] = fmt.Sprintf("%s must be at least %s characters", field, err.Param())
            case "max":
                errors[field] = fmt.Sprintf("%s must be at most %s characters", field, err.Param())
            case "gt":
                errors[field] = fmt.Sprintf("%s must be greater than %s", field, err.Param())
            case "gte":
                errors[field] = fmt.Sprintf("%s must be at least %s", field, err.Param())
            default:
                errors[field] = fmt.Sprintf("%s is invalid", field)
            }
        }

        return response.ValidationError(c, errors)
    }

    // Continue with business logic
    result, err := h.service.Create(&req)
    if err != nil {
        return response.InternalServerError(c, err.Error(), nil)
    }

    return response.Created(c, "Resource created successfully", result)
}
```

## All Validation Tags

| Tag | Description | Example |
|-----|-------------|---------|
| `required` | Field is mandatory | `validate:"required"` |
| `omitempty` | Skip validation if empty | `validate:"omitempty,email"` |
| `len=N` | Exact length | `validate:"len=10"` |
| `min=N` | Minimum value/length | `validate:"min=3"` |
| `max=N` | Maximum value/length | `validate:"max=100"` |
| `eq=N` | Equal to value | `validate:"eq=5"` |
| `ne=N` | Not equal to value | `validate:"ne=0"` |
| `gt=N` | Greater than | `validate:"gt=0"` |
| `gte=N` | Greater than or equal | `validate:"gte=18"` |
| `lt=N` | Less than | `validate:"lt=100"` |
| `lte=N` | Less than or equal | `validate:"lte=150"` |
| `oneof=a b c` | One of listed values | `validate:"oneof=admin user guest"` |
| `email` | Valid email | `validate:"email"` |
| `url` | Valid URL | `validate:"url"` |
| `uri` | Valid URI | `validate:"uri"` |
| `e164` | Valid E.164 phone | `validate:"e164"` |
| `uuid` | Valid UUID | `validate:"uuid"` |
| `alpha` | Alphabetic only | `validate:"alpha"` |
| `alphanum` | Alphanumeric only | `validate:"alphanum"` |
| `numeric` | Numeric string | `validate:"numeric"` |
| `number` | Numeric value | `validate:"number"` |
| `hexadecimal` | Hexadecimal string | `validate:"hexadecimal"` |
| `ip` | Valid IP address | `validate:"ip"` |
| `ipv4` | Valid IPv4 address | `validate:"ipv4"` |
| `ipv6` | Valid IPv6 address | `validate:"ipv6"` |
| `datetime=layout` | Valid datetime | `validate:"datetime=2006-01-02"` |
| `eqfield=Field` | Equal to another field | `validate:"eqfield=Password"` |
| `nefield=Field` | Not equal to another field | `validate:"nefield=OldPassword"` |
| `gtfield=Field` | Greater than another field | `validate:"gtfield=StartDate"` |
| `ltfield=Field` | Less than another field | `validate:"ltfield=EndDate"` |
| `required_with=Field` | Required if Field is present | `validate:"required_with=Phone"` |
| `required_without=Field` | Required if Field is not present | `validate:"required_without=Email"` |
| `dive` | Validate array/slice elements | `validate:"dive,min=1"` |

## Best Practices

1. **Always validate** - Every DTO should have validation tags
2. **Use response utility** - Always use `core/response` instead of `fiber.Map`
3. **Use constants** - Reference constants from `domain/constants/` for enums
4. **Update requests** - Use pointers with `omitempty` for optional fields
5. **Clear errors** - Format validation errors for better UX
6. **Business logic** - Complex validation should be in the service layer
7. **Consistent messages** - Use constants for error messages
8. **Test validation** - Write tests for validation rules

## Example: Complete Product DTO

```go
package product

import "time"

type CreateRequest struct {
    Name        string   `json:"name" validate:"required,min=3,max=255"`
    Description string   `json:"description" validate:"omitempty,max=1000"`
    Price       float64  `json:"price" validate:"required,gt=0"`
    Stock       int      `json:"stock" validate:"required,gte=0"`
    CategoryID  uint     `json:"category_id" validate:"required,gte=1"`
    Tags        []string `json:"tags" validate:"omitempty,min=1,max=10,dive,min=2,max=50"`
    Status      string   `json:"status" validate:"required,oneof=active inactive draft"`
    SKU         string   `json:"sku" validate:"required,alphanum,len=8"`
}

type UpdateRequest struct {
    Name        *string  `json:"name,omitempty" validate:"omitempty,min=3,max=255"`
    Description *string  `json:"description,omitempty" validate:"omitempty,max=1000"`
    Price       *float64 `json:"price,omitempty" validate:"omitempty,gt=0"`
    Stock       *int     `json:"stock,omitempty" validate:"omitempty,gte=0"`
    CategoryID  *uint    `json:"category_id,omitempty" validate:"omitempty,gte=1"`
    Status      *string  `json:"status,omitempty" validate:"omitempty,oneof=active inactive draft"`
}

type Response struct {
    ID          uint      `json:"id"`
    Name        string    `json:"name"`
    Description string    `json:"description"`
    Price       float64   `json:"price"`
    Stock       int       `json:"stock"`
    CategoryID  uint      `json:"category_id"`
    Status      string    `json:"status"`
    SKU         string    `json:"sku"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
```

---

## Summary

**Remember:**
1. Validation in DTOs is your first line of defense - Always validate input!
2. Use `core/response` utility for all API responses - No `fiber.Map`!
3. Format validation errors for better user experience
4. Use constants for enums and error messages

**Response Utility Quick Reference:**
```go
// Success
return response.Created(c, "Resource created", result)
return response.OK(c, result)
return response.NoContent(c)

// Errors
return response.BadRequest(c, "Invalid input", nil)
return response.ValidationError(c, err.Error())
return response.NotFound(c, "Resource not found")
return response.InternalServerError(c, err.Error(), nil)
```
