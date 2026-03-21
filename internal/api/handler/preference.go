package handler

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	apicontext "github.com/your-org/notification-center/internal/api/context"
	"github.com/your-org/notification-center/internal/model"
	"github.com/your-org/notification-center/internal/store"
)

type PreferenceHandler struct {
	store  store.PreferenceStore
	logger *slog.Logger
}

func NewPreferenceHandler(store store.PreferenceStore, logger *slog.Logger) *PreferenceHandler {
	return &PreferenceHandler{
		store:  store,
		logger: logger,
	}
}

func (h *PreferenceHandler) Get(c *gin.Context) {
	projectID, err := apicontext.ProjectIDFromContext(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no project context"})
		return
	}

	userID := c.Param("user_id")

	pref, err := h.store.Get(c.Request.Context(), projectID, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "preferences not found"})
		return
	}

	c.JSON(http.StatusOK, pref)
}

func (h *PreferenceHandler) Update(c *gin.Context) {
	projectID, err := apicontext.ProjectIDFromContext(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no project context"})
		return
	}

	userID := c.Param("user_id")

	var req model.PreferenceUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	pref := &model.UserPreference{
		ProjectID:       projectID,
		UserID:          userID,
		EnabledChannels: req.EnabledChannels,
		OptedOutEvents:  req.OptedOutEvents,
	}
	if req.QuietHoursStart != nil {
		pref.QuietHoursStart = *req.QuietHoursStart
	}
	if req.QuietHoursEnd != nil {
		pref.QuietHoursEnd = *req.QuietHoursEnd
	}

	if err := h.store.Upsert(c.Request.Context(), pref); err != nil {
		h.logger.Error("update preferences failed", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}
