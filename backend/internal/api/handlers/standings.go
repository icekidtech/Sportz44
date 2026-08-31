package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/icekidtech/Sportz44/backend/internal/services"
)

// StandingsHandler serves league table and top-scorer endpoints.
type StandingsHandler struct {
	standings *services.StandingsService
}

// NewStandingsHandler creates a new StandingsHandler.
func NewStandingsHandler(standings *services.StandingsService) *StandingsHandler {
	return &StandingsHandler{standings: standings}
}

// GetStandings handles GET /api/standings?league=<competitionID>&season=<season>
func (h *StandingsHandler) GetStandings(c *gin.Context) {
	compID, err := strconv.ParseUint(c.Query("league"), 10, 64)
	if err != nil || compID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "league query param is required"})
		return
	}
	season := c.Query("season")
	rows, err := h.standings.GetStandings(c.Request.Context(), uint(compID), season)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to compute standings"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"standings": rows})
}

// GetTopScorers handles GET /api/standings/:league/top-scorers?season=<season>
func (h *StandingsHandler) GetTopScorers(c *gin.Context) {
	compID, err := strconv.ParseUint(c.Param("league"), 10, 64)
	if err != nil || compID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid league id"})
		return
	}
	season := c.Query("season")
	limit := atoiDefault(c.Query("limit"), 20)
	rows, err := h.standings.GetTopScorers(c.Request.Context(), uint(compID), season, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to compute top scorers"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"top_scorers": rows})
}
