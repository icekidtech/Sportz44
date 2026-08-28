package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/icekidtech/Sportz44/backend/internal/repository"
	"github.com/icekidtech/Sportz44/backend/internal/services"
)

// PlayerHandler serves player endpoints.
type PlayerHandler struct {
	players *services.PlayerService
}

// NewPlayerHandler creates a new PlayerHandler.
func NewPlayerHandler(players *services.PlayerService) *PlayerHandler {
	return &PlayerHandler{players: players}
}

// ListByClub handles GET /api/players?club=<id>
func (h *PlayerHandler) ListByClub(c *gin.Context) {
	clubID, err := strconv.ParseUint(c.Query("club"), 10, 64)
	if err != nil || clubID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "club query param is required"})
		return
	}
	players, err := h.players.ListByClub(c.Request.Context(), uint(clubID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list players"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"players": players})
}

// Get handles GET /api/players/:id
func (h *PlayerHandler) Get(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid player id"})
		return
	}
	p, err := h.players.Get(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "player not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get player"})
		return
	}
	c.JSON(http.StatusOK, p)
}

// GetSeasonStats handles GET /api/players/:id/stats
func (h *PlayerHandler) GetSeasonStats(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid player id"})
		return
	}
	stats, err := h.players.GetSeasonStats(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get player stats"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"stats": stats})
}

// GetForm handles GET /api/players/:id/form
func (h *PlayerHandler) GetForm(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid player id"})
		return
	}
	limit := atoiDefault(c.Query("limit"), 10)
	form, err := h.players.GetForm(c.Request.Context(), uint(id), limit)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "player not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get player form"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"form": form})
}

// GetInjuries handles GET /api/players/:id/injuries
func (h *PlayerHandler) GetInjuries(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid player id"})
		return
	}
	injuries, err := h.players.GetInjuries(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get player injuries"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"injuries": injuries})
}
