# Module Skill - Quick Reference

## What is This?

A Claude Code skill that helps you create and manage module-based (package by feature) project structure for Go Fiber applications.

## Quick Commands

| Command | Description |
|---------|-------------|
| `/module create <name>` | Create a new module with full CRUD |
| `/module add <module> <resource>` | Add resource to existing module |
| `/module list` | List all modules in the project |
| `/module show <name>` | Show structure of a module |

## File Structure Created

```
modules/<module-name>/
├── handler.go       # HTTP handlers (controllers)
├── service.go       # Business logic
├── repository.go    # Database operations
├── model.go         # Domain entity
├── dto.go          # Request/Response structures
└── routes.go       # Route registration
```

## Dependency Flow

```
Handler → Service → Repository → Database
```

Each layer only depends on the layer below it.

## Example: Create User Module

```
User: "I need user management with CRUD operations"
You: /module create user
```

This creates all necessary files with boilerplate code for:
- Create user
- Get user by ID
- List users
- Update user
- Delete user

## What You Need to Do After Creation

1. **Customize the Model** (`model.go`)
   - Add specific fields for your entity
   - Add GORM tags and constraints

2. **Update DTOs** (`dto.go`)
   - Match request/response fields to your model
   - Add validation tags

3. **Implement Business Logic** (`service.go`)
   - Add custom validation
   - Add business rules
   - Add custom methods

4. **Register Routes** (`api/routes/setup.go`)
   ```go
   import "go-fiber/modules/user"

   func Setup(app *fiber.App, db *gorm.DB, rd *redis.Client) {
       v1 := app.Group("/api/v1")
       user.SetupRoutes(v1, db)
   }
   ```

5. **Run Migrations**
   - Create database migrations for your models
   - Run migrations

6. **Generate Swagger Docs**
   ```bash
   make swag
   ```

7. **Test**
   - Use Swagger UI at `http://localhost:3000/swagger`
   - Or test with curl/Postman

## Why Module-Based?

**Module structure:**
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
```

**Benefits:**
- All related code in one place
- Strong feature boundaries
- Easy to navigate
- Independent modules
- Great for teams
- Scalable architecture

## API Endpoints Generated

For a module named `user`, you get:

```
POST   /api/v1/user/      - Create
GET    /api/v1/user/      - List all
GET    /api/v1/user/:id   - Get by ID
PUT    /api/v1/user/:id   - Update
DELETE /api/v1/user/:id   - Delete
```

## Code Templates

The skill generates production-ready code with:
- ✅ Proper error handling
- ✅ Dependency injection
- ✅ Swagger documentation
- ✅ GORM integration
- ✅ Fiber best practices
- ✅ Request validation
- ✅ Response formatting

## Supporting Files

- `architecture.md` - Detailed explanation of module-based architecture
- `example-usage.md` - Real-world examples and customization
- `SKILL.md` - Full skill implementation with templates

## Learn More

1. Read `architecture.md` for architecture concepts
2. Check `example-usage.md` for practical examples
3. Use `/module create` to start building

## Tips

- Keep modules independent
- Use DTOs for API boundaries
- Put business logic in services
- Keep handlers thin
- Use dependency injection
- Return errors, don't panic
- Add Swagger docs as you go

---

**Ready to start?** Try:
```
/module create user
```
