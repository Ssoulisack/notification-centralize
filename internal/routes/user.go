package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/your-org/notification-center/internal/data/repository"
	"github.com/your-org/notification-center/internal/data/services"
	"github.com/your-org/notification-center/internal/handlers"
)

func SetupUserRoutes(router *gin.RouterGroup, deps *Dependencies) {
	userRepo := repository.NewUserRepository(deps.GormDB)
	projectRepo := repository.NewProjectRepository(deps.GormDB)
	userService := services.NewUserSyncService(userRepo, deps.Logger)
	handler := handlers.NewUserHandler(userService, projectRepo, deps.Logger)

	users := router.Group("/users")
	{
		users.GET("/me", handler.GetMe)
		users.PATCH("/me", handler.UpdateMe)
		users.GET("/me/projects", handler.GetMyProjects)
	}
}
