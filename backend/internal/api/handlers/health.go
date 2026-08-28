package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/icekidtech/Sportz44/backend/pkg/cache"
)

// HealthHandler reports service health.
type HealthHandler struct {
	db  *gorm.DB
	rdb *cache.Redis
}

// NewHealthHandler creates a HealthHandler.
func NewHealthHandler(db *gorm.DB, rdb *cache.Redis) *HealthHandler {
	return &HealthHandler{db: db, rdb: rdb}
}

// Check returns 200 when the API, database, and Redis are reachable.
func (h *HealthHandler) Check(c *gin.Context) {
	status := http.StatusOK
	body := gin.H{"status": "ok"}

	sqlDB, err := h.db.DB()
	if err != nil || sqlDB.Ping() != nil {
		status = http.StatusServiceUnavailable
		body["database"] = "unreachable"
	} else {
		body["database"] = "ok"
	}

	if err := h.rdb.Ping(c.Request.Context()); err != nil {
		status = http.StatusServiceUnavailable
		body["redis"] = "unreachable"
	} else {
		body["redis"] = "ok"
	}

	c.JSON(status, body)
}
