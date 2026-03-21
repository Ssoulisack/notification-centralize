package handler

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	apicontext "github.com/your-org/notification-center/internal/api/context"
	"github.com/your-org/notification-center/internal/model"
	"github.com/your-org/notification-center/internal/service"
)

type NotificationHandler struct {
	svc      *service.NotificationService
	validate *validator.Validate
	logger   *slog.Logger
}

func NewNotificationHandler(svc *service.NotificationService, logger *slog.Logger) *NotificationHandler {
	return &NotificationHandler{
		svc:      svc,
		validate: validator.New(),
		logger:   logger,
	}
}

func (h *NotificationHandler) Send(c *gin.Context) {
	projectID, err := apicontext.ProjectIDFromContext(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no project context"})
		return
	}

	var req model.SendRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON: " + err.Error()})
		return
	}

	if err := h.validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "validation: " + err.Error()})
		return
	}

	resp, err := h.svc.Send(c.Request.Context(), projectID, &req)
	if err != nil {
		h.logger.Error("send notification failed", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, resp)
}

func (h *NotificationHandler) Get(c *gin.Context) {
	projectID, err := apicontext.ProjectIDFromContext(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no project context"})
		return
	}

	id := c.Param("id")

	n, err := h.svc.GetStatus(c.Request.Context(), projectID, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "notification not found"})
		return
	}

	c.JSON(http.StatusOK, n)
}

func (h *NotificationHandler) List(c *gin.Context) {
	projectID, err := apicontext.ProjectIDFromContext(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no project context"})
		return
	}

	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	notifications, err := h.svc.ListByUser(c.Request.Context(), projectID, userID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"notifications": notifications, "count": len(notifications)})
}
