# Architecture Guide: Module vs Layer

## Quick Start - New Project

**Starting from scratch?** Use:
```
/init-project my-api
```

This creates a complete Go Fiber project with database, redis, middleware, docker, and everything you need!

---

## Adding Features

This project has two skills for adding features to your Go Fiber application:

## `/module` - Module-Based (Package by Feature)

**Structure:**
```
modules/
  user/
    handler.go
    service.go
    repository.go
    model.go
    dto.go
    routes.go
  product/
    handler.go
    service.go
    repository.go
    model.go
    dto.go
    routes.go
```

**Organization:** By business domain/feature
**Philosophy:** Keep all code for a feature together

**Use when:**
- ✅ Building medium to large projects
- ✅ Following Domain-Driven Design (DDD)
- ✅ Features have varying complexity
- ✅ Want strong feature boundaries
- ✅ Teams own specific domains
- ✅ Long-term maintainability is key

**Command:** `/module create user`

---

## `/layer` - Layer-Based (Package by Layer)

**Structure:**
```
api/controllers/
  user_controller.go
  product_controller.go
data/services/
  user_service.go
  product_service.go
data/repositories/
  user_repository.go
  product_repository.go
domain/models/
  user.go
  product.go
```

**Organization:** By technical layer
**Philosophy:** Separate by technical responsibility

**Use when:**
- ✅ Building small to medium projects
- ✅ Team prefers traditional MVC/layered architecture
- ✅ All features have similar complexity
- ✅ Want centralized view of controllers/services
- ✅ Following common Go patterns

**Command:** `/layer create user`

---

## Quick Comparison

| Aspect | Module-Based | Layer-Based |
|--------|--------------|-------------|
| **Organization** | By feature | By technical layer |
| **Structure** | `modules/user/` | `controllers/`, `services/` |
| **Navigation** | Everything in one folder | Jump between folders |
| **Best for** | Large projects, DDD | Small-medium projects |
| **Coupling** | Low | Can be higher |
| **Boundaries** | Business domain | Technical layer |
| **Scalability** | Excellent | Good for smaller apps |
| **Team ownership** | Teams own modules | Teams own layers |

---

## Example Scenarios

### Scenario 1: Simple Blog API
**Project:** Basic blog with posts, comments, users
**Features:** 3-5 simple CRUD resources
**Recommendation:** `/layer` ✅
**Why:** Simple project, consistent patterns, layer-based is sufficient

### Scenario 2: E-commerce Platform
**Project:** Products, orders, payments, inventory, shipping
**Features:** 10+ complex features with different requirements
**Recommendation:** `/module` ✅
**Why:** Complex domains, varying complexity, better boundaries

### Scenario 3: Learning/Tutorial Project
**Project:** Following a tutorial or learning Go Fiber
**Recommendation:** `/layer` ✅
**Why:** Most tutorials use this pattern, easier to learn

### Scenario 4: Enterprise Application
**Project:** Large business application with multiple teams
**Features:** Many features, long-term maintenance
**Recommendation:** `/module` ✅
**Why:** Team ownership, clear boundaries, better scalability

---

## Migration Between Styles

You can switch between styles as your project evolves:

### From Layer to Module (Growing complexity)
```bash
# Create new module structure
/module create user

# Move code:
# controllers/user_controller.go → modules/user/handler.go
# services/user_service.go → modules/user/service.go
# repositories/user_repository.go → modules/user/repository.go
# models/user.go → modules/user/model.go
```

### From Module to Layer (Simplifying)
```bash
# Create layer structure
/layer create user

# Move code:
# modules/user/handler.go → controllers/user_controller.go
# modules/user/service.go → services/user_service.go
# modules/user/repository.go → repositories/user_repository.go
# modules/user/model.go → models/user.go
```

---

## Decision Tree

```
Start Here
    |
    ├─ Project size?
    |   ├─ Small (1-5 features) → Use /layer
    |   ├─ Medium (5-10 features) → Your choice, either works
    |   └─ Large (10+ features) → Use /module
    |
    ├─ Team size?
    |   ├─ Solo or small team → Use /layer
    |   └─ Multiple teams → Use /module
    |
    ├─ Feature complexity?
    |   ├─ All similar → Use /layer
    |   └─ Varying complexity → Use /module
    |
    └─ Long-term maintenance?
        ├─ Short-term project → Use /layer
        └─ Long-term/enterprise → Use /module
```

---

## Still Not Sure?

**Start with `/layer`** if:
- You're new to Go
- Building a prototype or MVP
- Project scope is unclear
- You prefer simplicity

**You can always migrate to `/module` later when:**
- Project grows complex
- Need stronger boundaries
- Multiple teams join
- Following DDD principles

---

## Both Skills Generate

✅ Full CRUD operations
✅ Dependency injection
✅ Swagger documentation
✅ Error handling
✅ Request validation
✅ GORM integration
✅ Fiber best practices

The difference is only in **how files are organized**, not what they contain.

---

## Get Started

**New Project:**
```bash
/init-project my-api
```

**Add Module-Based Feature:**
```bash
/module create user
```

**Add Layer-Based Feature:**
```bash
/layer create user
```

Choose what feels right for your project!

---

## All Available Skills

| Skill | Purpose |
|-------|---------|
| `/init-project` | Initialize complete new project |
| `/module` | Add module-based feature |
| `/layer` | Add layer-based feature |
