package handlers

import (
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/your-org/notification-center/internal/middleware"
	"github.com/your-org/notification-center/internal/services"
)

// InboxHandler handles user inbox endpoints.
type InboxHandler struct {
	notifService *services.NotificationService
	userService  *services.UserSyncService
	logger       *slog.Logger
}

// NewInboxHandler creates a new inbox handler.
func NewInboxHandler(
	notifService *services.NotificationService,
	userService *services.UserSyncService,
	logger *slog.Logger,
) *InboxHandler {
	return &InboxHandler{
		notifService: notifService,
		userService:  userService,
		logger:       logger,
	}
}

// GetInbox returns the user's notification inbox.
// GET /inbox
func (h *InboxHandler) GetInbox(c *gin.Context) {
	user, err := middleware.GetUserFromContext(c)
	if err != nil {
		Unauthorized(c, "not authenticated")
		return
	}

	dbUser, err := h.userService.GetUserByKeycloakID(c.Request.Context(), user.KeycloakID)
	if err != nil {
		NotFound(c, "user not found")
		return
	}

	// Get project ID from query or use default
	projectIDStr := c.Query("project_id")
	if projectIDStr == "" {
		BadRequest(c, "project_id is required")
		return
	}

	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		BadRequest(c, "invalid project_id")
		return
	}

	limit, offset := GetPaginationParams(c)

	notifications, total, err := h.notifService.GetUserInbox(c.Request.Context(), projectID, dbUser.ID, limit, offset)
	if err != nil {
		h.logger.Error("failed to get inbox", "error", err)
		InternalError(c, "failed to get inbox")
		return
	}

	Paginated(c, notifications, total, limit, offset)
}

// GetUnreadCount returns the count of unread notifications.
// GET /inbox/unread/count
func (h *InboxHandler) GetUnreadCount(c *gin.Context) {
	user, err := middleware.GetUserFromContext(c)
	if err != nil {
		Unauthorized(c, "not authenticated")
		return
	}

	dbUser, err := h.userService.GetUserByKeycloakID(c.Request.Context(), user.KeycloakID)
	if err != nil {
		NotFound(c, "user not found")
		return
	}

	projectIDStr := c.Query("project_id")
	if projectIDStr == "" {
		BadRequest(c, "project_id is required")
		return
	}

	projectID, err := uuid.Parse(projectIDStr)
	if err != nil {
		BadRequest(c, "invalid project_id")
		return
	}

	count, err := h.notifService.GetUnreadCount(c.Request.Context(), projectID, dbUser.ID)
	if err != nil {
		h.logger.Error("failed to get unread count", "error", err)
		InternalError(c, "failed to get unread count")
		return
	}

	Success(c, gin.H{"unread_count": count})
}

// MarkAsRead marks a notification as read.
// POST /inbox/:id/read
func (h *InboxHandler) MarkAsRead(c *gin.Context) {
	user, err := middleware.GetUserFromContext(c)
	if err != nil {
		Unauthorized(c, "not authenticated")
		return
	}

	dbUser, err := h.userService.GetUserByKeycloakID(c.Request.Context(), user.KeycloakID)
	if err != nil {
		NotFound(c, "user not found")
		return
	}

	notificationID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		BadRequest(c, "invalid notification ID")
		return
	}

	if err := h.notifService.MarkAsRead(c.Request.Context(), notificationID, dbUser.ID); err != nil {
		h.logger.Error("failed to mark as read", "error", err)
		InternalError(c, "failed to mark as read")
		return
	}

	Success(c, gin.H{"message": "marked as read"})
}
