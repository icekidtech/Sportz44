package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/icekidtech/Sportz44/backend/internal/services"
)

// MatchHandler handles match-related HTTP requests.
type MatchHandler struct {
	svc *services.MatchService
}

// NewMatchHandler creates a new MatchHandler.
func NewMatchHandler(svc *services.MatchService) *MatchHandler {
	return &MatchHandler{svc: svc}
}

// SyncCompetition triggers a sync of fixtures for a competition/season.
// Query params: competition (provider ID/code), season (e.g. "2026").
func (h *MatchHandler) SyncCompetition(c *gin.Context) {
	competitionID := c.Query("competition")
	season := c.Query("season")
	if competitionID == "" || season == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "competition and season are required"})
		return
	}
	if err := h.svc.SyncCompetition(c.Request.Context(), competitionID, season); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "synced"})
}
