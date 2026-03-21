---
name: init-project
description: Initialize a new Go Fiber project with complete setup including structure, database, redis, middleware, docker, and swagger
disable-model-invocation: true
argument-hint: [project-name]
---

# Init Project Skill

This skill initializes a complete Go Fiber project from scratch with all necessary components and configuration.

## What This Skill Does

Creates a production-ready Go Fiber project with:

✅ **Project Structure** - Complete directory layout
✅ **Database Setup** - GORM with PostgreSQL/MySQL support
✅ **Redis Setup** - Redis client for caching
✅ **Middleware** - CORS, Logger, Recovery, JWT Auth, Rate Limiting
✅ **Docker** - Dockerfile and docker-compose.yml
✅ **Development Tools** - Makefile, .gitignore, Swagger
✅ **Example Feature** - Sample user module/resource to follow
✅ **Configuration** - config.yaml and environment setup

## Usage

```
/init-project <project-name>
```

Example:
```
/init-project my-awesome-api
```

## Implementation Steps

### Step 1: Ask User for Architecture Choice

Before creating the project, ask the user:

**Question:** "Which architecture pattern would you like to use?"

**Options:**
1. **Module-based (Package by Feature)** - Best for medium-large projects, DDD
   - Creates `modules/` structure
   - Each feature in its own folder
   - Example included in `modules/user/`

2. **Layer-based (Package by Layer)** - Best for small-medium projects
   - Creates `api/controllers/`, `data/services/`, `data/repositories/`
   - Features separated by technical layer
   - Example included across layers

### Step 2: Create Project Directory

```bash
mkdir -p <project-name>
cd <project-name>
```

### Step 3: Initialize Go Module

```bash
go mod init <project-name>
```

### Step 4: Create Directory Structure

**If Module-based:**
```
<project-name>/
├── api/
│   ├── middleware/
│   └── routes/
├── bootstrap/
├── core/
│   ├── utilities/
│   └── logs/
├── domain/
│   └── constants/
├── modules/
│   └── user/           # Example module
├── config.yaml
├── main.go
├── Dockerfile
├── docker-compose.yml
├── Makefile
└── .gitignore
```

**If Layer-based:**
```
<project-name>/
├── api/
│   ├── controllers/    # HTTP handlers
│   ├── middleware/
│   └── routes/
├── data/
│   ├── services/       # Business logic
│   └── repositories/   # Data access
├── domain/
│   ├── models/         # Entities
│   ├── dto/           # DTOs
│   └── constants/
├── bootstrap/
├── core/
│   ├── utilities/
│   └── logs/
├── config.yaml
├── main.go
├── Dockerfile
├── docker-compose.yml
├── Makefile
└── .gitignore
```

### Step 5: Generate Core Files

#### 5.1 main.go
```go
package main

import (
	"fmt"
	"log"
	"<project-name>/api/routes"
	"<project-name>/bootstrap"
)

func main() {
	app := bootstrap.App()
	globalEnv := app.Env
	fiber := app.Fiber
	db := app.Postgres
	rd := app.Redis

	routes.Setup(fiber, db, rd)

	log.Fatal(fiber.Listen(fmt.Sprintf(":%d", globalEnv.App.Port)))
}
```

#### 5.2 bootstrap/app.go
```go
package bootstrap

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Application struct {
	Env      *Config
	Fiber    *fiber.App
	Postgres *gorm.DB
	Redis    *redis.Client
}

func App() *Application {
	app := &Application{}
	app.Env = NewConfig()
	app.Postgres = NewPostgresDatabase(app.Env)
	app.Redis = NewRedisClient(app.Env)
	app.Fiber = NewFiberApp()

	return app
}

func NewFiberApp() *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:      "Go Fiber API",
		ServerHeader: "Fiber",
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))

	return app
}
```

#### 5.3 bootstrap/config.go
```go
package bootstrap

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	App      AppConfig      ` + "`yaml:\"App\"`" + `
	Database DatabaseConfig ` + "`yaml:\"database\"`" + `
	Redis    RedisConfig    ` + "`yaml:\"redis\"`" + `
	JWT      JWTConfig      ` + "`yaml:\"jwt\"`" + `
}

type AppConfig struct {
	Port int ` + "`yaml:\"port\"`" + `
}

type DatabaseConfig struct {
	Host     string ` + "`yaml:\"master_host\"`" + `
	Port     int    ` + "`yaml:\"master_port\"`" + `
	Username string ` + "`yaml:\"master_username\"`" + `
	Password string ` + "`yaml:\"master_password\"`" + `
	DBName   string ` + "`yaml:\"master_dbname\"`" + `
}

type RedisConfig struct {
	Host     string ` + "`yaml:\"host\"`" + `
	Port     int    ` + "`yaml:\"port\"`" + `
	Password string ` + "`yaml:\"password\"`" + `
	DB       int    ` + "`yaml:\"db\"`" + `
}

type JWTConfig struct {
	AccessSecret  string ` + "`yaml:\"access_token\"`" + `
	RefreshSecret string ` + "`yaml:\"refresh_token\"`" + `
}

func NewConfig() *Config {
	data, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		log.Fatalf("Failed to parse config: %v", err)
	}

	return &config
}
```

#### 5.4 bootstrap/database.go
```go
package bootstrap

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewPostgresDatabase(config *Config) *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Database.Host,
		config.Database.Port,
		config.Database.Username,
		config.Database.Password,
		config.Database.DBName,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Database connected successfully")
	return db
}
```

#### 5.5 bootstrap/redis.go
```go
package bootstrap

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

func NewRedisClient(config *Config) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Redis.Host, config.Redis.Port),
		Password: config.Redis.Password,
		DB:       config.Redis.DB,
	})

	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	log.Println("Redis connected successfully")
	return client
}
```

#### 5.6 config.yaml
```yaml
App:
  port: 3000

database:
  master_host: localhost
  master_port: 5432
  master_username: postgres
  master_password: postgres
  master_dbname: myapp

redis:
  host: localhost
  port: 6379
  password: ""
  db: 0

jwt:
  access_token: your-secret-key-here
  refresh_token: your-refresh-secret-here
```

### Step 6: Create Middleware

#### 6.1 api/middleware/jwt.go
```go
package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuth(secret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing authorization header",
			})
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token",
			})
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.Locals("userID", claims["user_id"])
		}

		return c.Next()
	}
}
```

### Step 7: Create Routes Setup

Create `api/routes/setup.go` based on architecture choice.

**If Module-based:**
```go
package routes

import (
	"<project-name>/modules/user"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func Setup(app *fiber.App, db *gorm.DB, rd *redis.Client) {
	api := app.Group("/api/v1")

	// Register module routes
	user.SetupRoutes(api, db)

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
		})
	})
}
```

**If Layer-based:**
```go
package routes

import (
	"<project-name>/api/controllers"
	"<project-name>/data/repositories"
	"<project-name>/data/services"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func Setup(app *fiber.App, db *gorm.DB, rd *redis.Client) {
	api := app.Group("/api/v1")

	// User routes
	userRepo := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepo)
	userController := controllers.NewUserController(userService)

	users := api.Group("/users")
	users.Post("/", userController.Create)
	users.Get("/", userController.List)
	users.Get("/:id", userController.GetByID)
	users.Put("/:id", userController.Update)
	users.Delete("/:id", userController.Delete)

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
		})
	})
}
```

### Step 8: Create Example User Feature

**If Module-based:** Use `/module create user`
**If Layer-based:** Use `/layer create user`

### Step 9: Create Docker Files

#### 9.1 Dockerfile
```dockerfile
FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o main .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/config.yaml .

EXPOSE 3000

CMD ["./main"]
```

#### 9.2 docker-compose.yml
```yaml
version: '3.8'

services:
  app:
    build: .
    ports:
      - "3000:3000"
    depends_on:
      - postgres
      - redis
    environment:
      - DATABASE_HOST=postgres
      - REDIS_HOST=redis
    volumes:
      - .:/app

  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
      POSTGRES_DB: myapp
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data

volumes:
  postgres_data:
  redis_data:
```

### Step 10: Create Development Tools

#### 10.1 Makefile
```makefile
.PHONY: run swag build clean docker-up docker-down

run:
	go run main.go

swag:
	swag init -g ./bootstrap/swagger.go

build:
	go build -o bin/app main.go

clean:
	rm -rf bin/

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-build:
	docker-compose up --build -d

test:
	go test -v ./...

install:
	go mod download
	go install github.com/swaggo/swag/cmd/swag@latest
```

#### 10.2 .gitignore
```
# Binaries
*.exe
*.exe~
*.dll
*.so
*.dylib
bin/
dist/

# Test binary
*.test

# Output
*.out

# Go workspace file
go.work

# Dependencies
vendor/

# IDE
.vscode/
.idea/
*.swp
*.swo
*~

# OS
.DS_Store
Thumbs.db

# Config (if contains secrets)
# config.yaml

# Logs
*.log

# Swagger
docs/
```

### Step 11: Install Dependencies

```bash
go get github.com/gofiber/fiber/v2
go get gorm.io/gorm
go get gorm.io/driver/postgres
go get github.com/redis/go-redis/v9
go get github.com/golang-jwt/jwt/v5
go get gopkg.in/yaml.v3
go get github.com/swaggo/swag/cmd/swag
go get github.com/swaggo/fiber-swagger
```

### Step 12: Initialize Swagger

Create `bootstrap/swagger.go`:
```go
package bootstrap

// @title Go Fiber API
// @version 1.0
// @description API Server for Go Fiber Application
// @host localhost:3000
// @BasePath /api/v1
// @schemes http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
```

### Step 13: Display Success Message

After everything is created, show:

```
✅ Project "<project-name>" initialized successfully!

📁 Project structure created with <architecture-choice> architecture
🗄️  Database (PostgreSQL) configured
🔴 Redis configured
🔧 Middleware (CORS, Logger, Recovery, JWT) added
🐳 Docker files created
🛠️  Makefile created
📝 Example user feature created

Next steps:
1. cd <project-name>
2. Update config.yaml with your database credentials
3. Install Swagger: go install github.com/swaggo/swag/cmd/swag@latest
4. Generate Swagger docs: make swag
5. Run with Docker: make docker-up
   OR run locally: make run
6. Visit http://localhost:3000/health to check if it's running
7. Visit http://localhost:3000/swagger/index.html for API docs

To add more features:
- Module-based: /module create <feature-name>
- Layer-based: /layer create <feature-name>

Happy coding! 🚀
```

## Important Notes

1. **Architecture Choice** - Always ask the user first
2. **Database** - Default to PostgreSQL, but code should work with MySQL too
3. **Configuration** - Use config.yaml for easy configuration
4. **Example Feature** - Include user example to show the pattern
5. **Production Ready** - Include all essential middleware and error handling
6. **Docker Support** - Full docker-compose setup for development
7. **Swagger** - Pre-configured for API documentation

## Example Usage

```
User: "I want to start a new Go Fiber project for a blog API"
Assistant: /init-project blog-api

[Asks: Which architecture?]
User selects: Module-based

[Creates complete project]

User can now:
/module create post
/module create comment
/module create category
```

The project is immediately ready for development!
