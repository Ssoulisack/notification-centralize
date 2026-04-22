package handlers

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/your-org/notification-center/internal/data/services"
	"github.com/your-org/notification-center/internal/domain/dto"
	"github.com/your-org/notification-center/pkg/middleware"
	"github.com/your-org/notification-center/pkg/response"
)

type NotificationHandler interface {
	Send(c *gin.Context)
	SendBatch(c *gin.Context)
	List(c *gin.Context)
	Get(c *gin.Context)
}

// NotificationHandler handles notification endpoints.
type notificationHandler struct {
	notifService services.NotificationService
	userService  services.UserSyncService
	logger       *slog.Logger
}

// NewNotificationHandler creates a new notification handler.
func NewNotificationHandler(
	notifService services.NotificationService,
	userService services.UserSyncService,
	logger *slog.Logger,
) NotificationHandler {
	return &notificationHandler{
		notifService: notifService,
		userService:  userService,
		logger:       logger,
	}
}

// Send sends a notification.
// POST /projects/:id/notifications
func (h *notificationHandler) Send(c *gin.Context) {
	projectID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid project ID")
		return
	}

	user, err := middleware.GetUserFromContext(c)
	if err != nil {
		response.Unauthorized(c, "not authenticated")
		return
	}

	dbUser, err := h.userService.GetUserByKeycloakID(c.Request.Context(), user.KeycloakID)
	if err != nil {
		response.NotFound(c, "user not found")
		return
	}

	var req dto.SendRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	notification, err := h.notifService.Send(c.Request.Context(), projectID, dbUser.ID, &req)
	if err != nil {
		h.logger.Error("failed to send notification", "error", err)
		response.InternalError(c, "failed to send notification")
		return
	}

	response.Created(c, notification)
}

// SendBatch sends multiple notifications.
// POST /projects/:id/notifications/batch
func (h *notificationHandler) SendBatch(c *gin.Context) {
	projectID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid project ID")
		return
	}

	user, err := middleware.GetUserFromContext(c)
	if err != nil {
		response.Unauthorized(c, "not authenticated")
		return
	}

	dbUser, err := h.userService.GetUserByKeycloakID(c.Request.Context(), user.KeycloakID)
	if err != nil {
		response.NotFound(c, "user not found")
		return
	}

	var req struct {
		Notifications []dto.SendRequest `json:"notifications" binding:"required,min=1,max=100"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	var results []gin.H
	for _, notifReq := range req.Notifications {
		notification, err := h.notifService.Send(c.Request.Context(), projectID, dbUser.ID, &notifReq)
		if err != nil {
			results = append(results, gin.H{
				"success": false,
				"error":   err.Error(),
			})
		} else {
			results = append(results, gin.H{
				"success":      true,
				"notification": notification,
			})
		}
	}

	response.Success(c, gin.H{"results": results})
}

// List returns notifications for a project.
// GET /projects/:id/notifications
func (h *notificationHandler) List(c *gin.Context) {
	projectID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid project ID")
		return
	}

	limit, offset := response.GetPaginationParams(c)

	notifications, total, err := h.notifService.ListNotifications(c.Request.Context(), projectID, limit, offset)
	if err != nil {
		h.logger.Error("failed to list notifications", "error", err)
		response.InternalError(c, "failed to list notifications")
		return
	}

	response.Paginated(c, notifications, total, limit, offset)
}

// Get returns a single notification.
// GET /projects/:id/notifications/:notificationId
func (h *notificationHandler) Get(c *gin.Context) {
	projectID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "invalid project ID")
		return
	}

	notificationID, err := uuid.Parse(c.Param("notificationId"))
	if err != nil {
		response.BadRequest(c, "invalid notification ID")
		return
	}

	notification, err := h.notifService.GetNotification(c.Request.Context(), projectID, notificationID)
	if err != nil {
		response.NotFound(c, "notification not found")
		return
	}

	response.Success(c, notification)
}
