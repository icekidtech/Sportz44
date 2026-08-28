package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/icekidtech/Sportz44/backend/internal/api/middleware"
	"github.com/icekidtech/Sportz44/backend/internal/models"
	"github.com/icekidtech/Sportz44/backend/internal/services"
)

// UserHandler serves user subscription and notification-preference endpoints.
type UserHandler struct {
	users *services.UserService
}

// NewUserHandler creates a new UserHandler.
func NewUserHandler(users *services.UserService) *UserHandler {
	return &UserHandler{users: users}
}

// ListSubscriptions handles GET /api/users/me/subscriptions
func (h *UserHandler) ListSubscriptions(c *gin.Context) {
	userID := c.GetUint(middleware.ContextUserID)
	subs, err := h.users.ListSubscriptions(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list subscriptions"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"subscriptions": subs})
}

// AddSubscription handles POST /api/users/me/subscriptions { "club_id": 1 }
func (h *UserHandler) AddSubscription(c *gin.Context) {
	var req struct {
		ClubID uint `json:"club_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "club_id is required"})
		return
	}
	userID := c.GetUint(middleware.ContextUserID)
	if err := h.users.AddSubscription(c.Request.Context(), userID, req.ClubID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add subscription"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"status": "subscribed"})
}

// RemoveSubscription handles DELETE /api/users/me/subscriptions/:clubID
func (h *UserHandler) RemoveSubscription(c *gin.Context) {
	clubID, err := strconv.ParseUint(c.Param("clubID"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid club id"})
		return
	}
	userID := c.GetUint(middleware.ContextUserID)
	if err := h.users.RemoveSubscription(c.Request.Context(), userID, uint(clubID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to remove subscription"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "unsubscribed"})
}

// GetNotificationPrefs handles GET /api/users/me/notification-preferences
func (h *UserHandler) GetNotificationPrefs(c *gin.Context) {
	userID := c.GetUint(middleware.ContextUserID)
	prefs, err := h.users.GetNotificationPrefs(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get notification preferences"})
		return
	}
	c.JSON(http.StatusOK, prefs)
}

// UpdateNotificationPrefs handles PUT /api/users/me/notification-preferences
func (h *UserHandler) UpdateNotificationPrefs(c *gin.Context) {
	userID := c.GetUint(middleware.ContextUserID)
	prefs, err := h.users.GetNotificationPrefs(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get notification preferences"})
		return
	}
	var req models.NotificationPreference
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	// Preserve the user ID and primary key; only update preference fields.
	prefs.KickoffReminder = req.KickoffReminder
	prefs.GoalAlerts = req.GoalAlerts
	prefs.RedCardAlerts = req.RedCardAlerts
	prefs.FullTimeResult = req.FullTimeResult
	prefs.NewsAlerts = req.NewsAlerts
	prefs.WhatsAppEnabled = req.WhatsAppEnabled
	prefs.PushEnabled = req.PushEnabled
	prefs.DoNotDisturbStart = req.DoNotDisturbStart
	prefs.DoNotDisturbEnd = req.DoNotDisturbEnd
	if err := h.users.UpdateNotificationPrefs(c.Request.Context(), prefs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update notification preferences"})
		return
	}
	c.JSON(http.StatusOK, prefs)
}
