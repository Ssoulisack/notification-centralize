package handler

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	apicontext "github.com/your-org/notification-center/internal/api/context"
	"github.com/your-org/notification-center/internal/model"
	"github.com/your-org/notification-center/internal/store"
)

type DeviceHandler struct {
	store    store.DeviceStore
	validate *validator.Validate
	logger   *slog.Logger
}

func NewDeviceHandler(store store.DeviceStore, logger *slog.Logger) *DeviceHandler {
	return &DeviceHandler{
		store:    store,
		validate: validator.New(),
		logger:   logger,
	}
}

func (h *DeviceHandler) Register(c *gin.Context) {
	projectID, err := apicontext.ProjectIDFromContext(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no project context"})
		return
	}

	var req model.DeviceRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	device := &model.DeviceToken{
		ID:         uuid.New().String(),
		ProjectID:  projectID,
		UserID:     req.UserID,
		Token:      req.Token,
		Platform:   req.Platform,
		AppVersion: req.AppVersion,
		CreatedAt:  time.Now(),
		LastUsedAt: time.Now(),
	}

	if err := h.store.Register(c.Request.Context(), device); err != nil {
		h.logger.Error("register device failed", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": device.ID, "status": "registered"})
}

func (h *DeviceHandler) Remove(c *gin.Context) {
	projectID, err := apicontext.ProjectIDFromContext(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no project context"})
		return
	}

	id := c.Param("id")

	if err := h.store.Remove(c.Request.Context(), projectID, id); err != nil {
		h.logger.Error("remove device failed", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "removed"})
}
