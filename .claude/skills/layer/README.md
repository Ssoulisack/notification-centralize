# Layer Skill - Quick Reference

## What is This?

A Claude Code skill for creating **layer-based (package by layer)** project structure for Go Fiber applications.

## Quick Commands

| Command | Description |
|---------|-------------|
| `/layer create <name>` | Create a new resource with full CRUD |
| `/layer list` | List all resources in the project |
| `/layer show <name>` | Show files for a resource |

## File Structure Created

```
api/controllers/<resource>_controller.go    # HTTP handlers
data/services/<resource>_service.go          # Business logic
data/repositories/<resource>_repository.go   # Database operations
domain/models/<resource>.go                  # Domain entity
domain/dto/<resource>_dto.go                # DTOs
```

## Dependency Flow

```
Controller → Service → Repository → Database
```

## Example: Create User Resource

```
User: "I need user management"
You: /layer create user
```

This creates:
- `api/controllers/user_controller.go`
- `data/services/user_service.go`
- `data/repositories/user_repository.go`
- `domain/models/user.go`
- `domain/dto/user_dto.go`

With CRUD operations:
- Create user
- Get user by ID
- List users
- Update user
- Delete user

## Layer-Based Structure

All code organized by **technical layer**:

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

✅ **Use `/layer` for:**
- Small to medium projects (3-10 features)
- Traditional MVC/layered architecture
- Teams familiar with this pattern
- When features have similar complexity
- Learning/educational projects

❌ **Use `/module` instead for:**
- Large projects (10+ features)
- Domain-Driven Design
- Features with varying complexity
- Stronger feature boundaries

## API Endpoints Generated

For a resource named `user`:

```
POST   /api/v1/user/      - Create
GET    /api/v1/user/      - List all
GET    /api/v1/user/:id   - Get by ID
PUT    /api/v1/user/:id   - Update
DELETE /api/v1/user/:id   - Delete
```

## After Creation

1. **Customize the Model** (`domain/models/<resource>.go`)
2. **Update DTOs** (`domain/dto/<resource>_dto.go`)
3. **Add Business Logic** (`data/services/<resource>_service.go`)
4. **Routes auto-registered** in `api/routes/setup.go`
5. **Run migrations**
6. **Generate Swagger**: `make swag`
7. **Test endpoints**

## Code Features

The skill generates:
- ✅ Full CRUD operations
- ✅ Dependency injection
- ✅ Swagger documentation
- ✅ Error handling
- ✅ Request validation
- ✅ GORM integration
- ✅ Fiber best practices

## Comparison

| Aspect | Layer-Based (`/layer`) | Module-Based (`/module`) |
|--------|----------------------|------------------------|
| Organization | By technical layer | By feature/domain |
| Structure | `controllers/`, `services/` | `modules/user/`, `modules/product/` |
| Best for | Small-medium projects | Large projects |
| Navigation | Jump between folders | Everything in one folder |
| Coupling | Can be higher | Lower |
| Boundaries | Technical | Business domain |

## Tips

- Keep controllers thin
- Put business logic in services
- Use DTOs for API boundaries
- Return errors, don't panic
- Add Swagger docs
- Use dependency injection

---

**Ready to start?** Try:
```
/layer create user
```
