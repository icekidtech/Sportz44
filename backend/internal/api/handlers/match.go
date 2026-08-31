package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/icekidtech/Sportz44/backend/internal/repository"
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

// ListMatches returns a paginated list of matches.
// Query params: club, competition, status, season, date_from, date_to, page, limit.
func (h *MatchHandler) ListMatches(c *gin.Context) {
	f := repository.MatchFilters{
		Status: c.Query("status"),
		Season: c.Query("season"),
		Page:   atoiDefault(c.Query("page"), 1),
		Limit:  atoiDefault(c.Query("limit"), 50),
	}
	if v := c.Query("club"); v != "" {
		f.ClubID = uint(atoiDefault(v, 0))
	}
	if v := c.Query("competition"); v != "" {
		f.CompetitionID = uint(atoiDefault(v, 0))
	}
	if v := c.Query("date_from"); v != "" {
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			f.DateFrom = &t
		}
	}
	if v := c.Query("date_to"); v != "" {
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			f.DateTo = &t
		}
	}

	matches, total, err := h.svc.ListMatches(c.Request.Context(), f)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"matches": matches,
		"total":   total,
		"page":    f.Page,
		"limit":   f.Limit,
	})
}

// GetMatch returns a single match by ID.
func (h *MatchHandler) GetMatch(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid match id"})
		return
	}
	m, err := h.svc.GetMatch(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "match not found"})
		return
	}
	c.JSON(http.StatusOK, m)
}

// GetMatchEvents returns events for a match.
func (h *MatchHandler) GetMatchEvents(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid match id"})
		return
	}
	events, err := h.svc.GetMatchEvents(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"events": events})
}

// GetMatchLineup returns the lineup for a match.
func (h *MatchHandler) GetMatchLineup(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid match id"})
		return
	}
	lineup, err := h.svc.GetMatchLineup(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"lineup": lineup})
}

// GetMatchStats returns the statistics for a match.
func (h *MatchHandler) GetMatchStats(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid match id"})
		return
	}
	stats, err := h.svc.GetMatchStats(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"stats": stats})
}

// SyncMatchEvents triggers a backfill of events for a single match.
func (h *MatchHandler) SyncMatchEvents(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid match id"})
		return
	}
	if err := h.svc.SyncMatchEvents(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "synced"})
}

// SyncMatchLineup triggers a backfill of the lineup for a single match.
func (h *MatchHandler) SyncMatchLineup(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid match id"})
		return
	}
	if err := h.svc.SyncMatchLineup(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "synced"})
}

// SyncMatchStats triggers a backfill of the statistics for a single match.
func (h *MatchHandler) SyncMatchStats(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid match id"})
		return
	}
	if err := h.svc.SyncMatchStats(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "synced"})
}

// SyncMatchDetails triggers a backfill of events, lineup, and stats for a
// single match in one call.
func (h *MatchHandler) SyncMatchDetails(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid match id"})
		return
	}
	if err := h.svc.SyncMatchDetails(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "synced"})
}

func atoiDefault(s string, def int) int {
	if s == "" {
		return def
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return n
}
