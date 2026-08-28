package api

import (
	"github.com/gin-gonic/gin"

	"github.com/icekidtech/Sportz44/backend/internal/api/handlers"
	"github.com/icekidtech/Sportz44/backend/internal/api/middleware"
)

// RegisterRoutes wires up all HTTP routes on the Gin engine.
func RegisterRoutes(r *gin.Engine, auth *handlers.AuthHandler, health *handlers.HealthHandler, match *handlers.MatchHandler, jwtSecret string) {
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

	// Protected API group: requires a valid access cookie + CSRF token for
	// state-changing requests. Feature handlers will be added here.
	protected := r.Group("/api")
	protected.Use(middleware.Auth(jwtSecret), middleware.CSRFProtect())
	{
		_ = protected // feature routes land here next (matches, players, community, ...)
		protected.GET("/matches/sync", match.SyncCompetition)
	}
}
