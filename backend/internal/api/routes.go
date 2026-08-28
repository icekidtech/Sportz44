package api

import (
	"github.com/gin-gonic/gin"

	"github.com/icekidtech/Sportz44/backend/internal/api/handlers"
	"github.com/icekidtech/Sportz44/backend/internal/api/middleware"
)

// RegisterRoutes wires up all HTTP routes on the Gin engine.
func RegisterRoutes(r *gin.Engine, auth *handlers.AuthHandler, health *handlers.HealthHandler, match *handlers.MatchHandler, ws *handlers.WSHandler, jwtSecret string) {
	r.GET("/health", health.Check)

	// Public auth endpoints (cookies are set here).
	authGroup := r.Group("/api/auth")
	{
		authGroup.POST("/register", auth.Register)
		authGroup.POST("/login", auth.Login)
		authGroup.POST("/refresh", auth.Refresh)
		authGroup.POST("/logout", auth.Logout)
		authGroup.GET("/me", middleware.Auth(jwtSecret), auth.Me)
	}

	// Public WebSocket endpoint for live match updates.
	r.GET("/ws", ws.Handle)

	// Protected API group: requires a valid access cookie + CSRF token for
	// state-changing requests. Feature handlers will be added here.
	protected := r.Group("/api")
	protected.Use(middleware.Auth(jwtSecret), middleware.CSRFProtect())
	{
		// Match query endpoints (read-only).
		protected.GET("/matches", match.ListMatches)
		protected.GET("/matches/:id", match.GetMatch)
		protected.GET("/matches/:id/events", match.GetMatchEvents)
		protected.GET("/matches/:id/lineup", match.GetMatchLineup)

		// Match ingestion (admin-triggered sync).
		protected.GET("/matches/sync", match.SyncCompetition)
	}
}
