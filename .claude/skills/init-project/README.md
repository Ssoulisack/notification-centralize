# Init Project Skill - Quick Reference

## What is This?

A skill to initialize a **complete Go Fiber project** from scratch with everything you need to start coding immediately.

## Quick Start

```
/init-project my-api
```

You'll be asked which architecture to use:
- **Module-based** - For medium-large projects
- **Layer-based** - For small-medium projects

## What Gets Created

### ✅ Project Structure
Complete directory layout based on your architecture choice

### ✅ Database Setup
- PostgreSQL with GORM
- Connection pooling
- Migration-ready

### ✅ Redis Setup
- Redis client configured
- Ready for caching and sessions

### ✅ Middleware
- CORS - Cross-origin requests
- Logger - Request logging
- Recovery - Panic recovery
- JWT Auth - Token authentication
- Rate Limiting - API rate limits

### ✅ Docker
- Dockerfile for app
- docker-compose.yml with PostgreSQL and Redis
- Ready for containerized development

### ✅ Development Tools
- Makefile with common commands
- .gitignore configured
- Swagger/OpenAPI setup
- Health check endpoint

### ✅ Example Feature
- Complete user feature as reference
- Shows the pattern to follow
- Ready to customize

## Directory Structure

### Module-Based
```
my-api/
├── api/
│   ├── middleware/
│   └── routes/
├── bootstrap/
│   ├── app.go
│   ├── config.go
│   ├── database.go
│   └── redis.go
├── core/
│   ├── utilities/
│   └── logs/
├── domain/
│   └── constants/
├── modules/
│   └── user/
│       ├── handler.go
│       ├── service.go
│       ├── repository.go
│       ├── model.go
│       ├── dto.go
│       └── routes.go
├── config.yaml
├── main.go
├── Dockerfile
├── docker-compose.yml
├── Makefile
└── .gitignore
```

### Layer-Based
```
my-api/
├── api/
│   ├── controllers/
│   ├── middleware/
│   └── routes/
├── data/
│   ├── services/
│   └── repositories/
├── domain/
│   ├── models/
│   ├── dto/
│   └── constants/
├── bootstrap/
├── core/
├── config.yaml
├── main.go
├── Dockerfile
├── docker-compose.yml
├── Makefile
└── .gitignore
```

## After Initialization

### 1. Configure Database
Edit `config.yaml`:
```yaml
database:
  master_host: localhost
  master_port: 5432
  master_username: postgres
  master_password: your_password
  master_dbname: your_database
```

### 2. Install Swagger
```bash
go install github.com/swaggo/swag/cmd/swag@latest
export PATH=$PATH:$(go env GOPATH)/bin
```

### 3. Generate Swagger Docs
```bash
make swag
```

### 4. Run the Project

**With Docker:**
```bash
make docker-up
```

**Without Docker:**
```bash
# Make sure PostgreSQL and Redis are running
make run
```

### 5. Test It
- Health check: `http://localhost:3000/health`
- Swagger UI: `http://localhost:3000/swagger/index.html`
- API base: `http://localhost:3000/api/v1`

## Makefile Commands

| Command | Description |
|---------|-------------|
| `make run` | Run the application |
| `make swag` | Generate Swagger docs |
| `make build` | Build binary |
| `make docker-up` | Start with Docker |
| `make docker-down` | Stop Docker containers |
| `make test` | Run tests |
| `make install` | Install dependencies |

## Adding Features

**Module-based:**
```
/module create product
/module create order
```

**Layer-based:**
```
/layer create product
/layer create order
```

## What's Included

### Configuration (config.yaml)
- App settings (port)
- Database connection
- Redis connection
- JWT secrets

### Bootstrap Files
- `app.go` - Application initialization
- `config.go` - Configuration loading
- `database.go` - Database connection
- `redis.go` - Redis connection
- `swagger.go` - Swagger setup

### Middleware
- `jwt.go` - JWT authentication

### Routes
- Health check endpoint
- API v1 group setup
- Example user routes

### Docker
- Multi-stage Dockerfile
- docker-compose with PostgreSQL, Redis, and app
- Volumes for data persistence

## Example Usage

### Initialize Project
```
User: "I want to create a new e-commerce API"
You: /init-project ecommerce-api

[Prompts for architecture choice]
User: "Module-based"

[Creates complete project]

Success! Project created at ./ecommerce-api/
```

### Add Features
```
cd ecommerce-api
/module create product
/module create cart
/module create order
```

### Run
```
make docker-up
```

Visit `http://localhost:3000/swagger` to see your API!

## Dependencies Installed

- `github.com/gofiber/fiber/v2` - Web framework
- `gorm.io/gorm` - ORM
- `gorm.io/driver/postgres` - PostgreSQL driver
- `github.com/redis/go-redis/v9` - Redis client
- `github.com/golang-jwt/jwt/v5` - JWT
- `gopkg.in/yaml.v3` - YAML parser
- `github.com/swaggo/swag` - Swagger generator
- `github.com/swaggo/fiber-swagger` - Swagger UI

## Project Features

✅ Clean architecture
✅ Environment configuration
✅ Database migrations ready
✅ Redis caching ready
✅ JWT authentication
✅ CORS configured
✅ Request logging
✅ Error recovery
✅ API documentation
✅ Docker support
✅ Development tools
✅ Example code

## Notes

- Uses **PostgreSQL** by default (can switch to MySQL)
- **Port 3000** by default
- **JWT** secrets should be changed in production
- **Swagger** docs auto-generated from code comments
- **Docker Compose** includes all services

---

**Ready to start coding?**
```
/init-project my-awesome-api
```

Your complete Go Fiber project will be ready in seconds! 🚀
