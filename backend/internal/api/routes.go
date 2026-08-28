package api

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"github.com/icekidtech/Sportz44/backend/internal/api/handlers"
	"github.com/icekidtech/Sportz44/backend/internal/api/middleware"
)

// RegisterRoutes wires up all HTTP routes on the Gin engine.
func RegisterRoutes(
	r *gin.Engine,
	auth *handlers.AuthHandler,
	health *handlers.HealthHandler,
	match *handlers.MatchHandler,
	player *handlers.PlayerHandler,
	standings *handlers.StandingsHandler,
	user *handlers.UserHandler,
	ws *handlers.WSHandler,
	rdb *redis.Client,
	rateLimitRequests int,
	rateLimitWindow time.Duration,
	jwtSecret string,
) {
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
	// state-changing requests, plus per-user rate limiting.
	protected := r.Group("/api")
	protected.Use(
		middleware.Auth(jwtSecret),
		middleware.CSRFProtect(),
		middleware.RateLimit(rdb, rateLimitRequests, rateLimitWindow),
	)
	{
		// Match query endpoints (read-only).
		protected.GET("/matches", match.ListMatches)
		protected.GET("/matches/:id", match.GetMatch)
		protected.GET("/matches/:id/events", match.GetMatchEvents)
		protected.GET("/matches/:id/lineup", match.GetMatchLineup)
		protected.GET("/matches/:id/stats", match.GetMatchStats)

		// Match ingestion (admin-triggered sync).
		protected.GET("/matches/sync", match.SyncCompetition)
		protected.GET("/matches/:id/sync-events", match.SyncMatchEvents)

		// Player endpoints.
		protected.GET("/players", player.ListByClub)
		protected.GET("/players/:id", player.Get)
		protected.GET("/players/:id/stats", player.GetSeasonStats)
		protected.GET("/players/:id/form", player.GetForm)
		protected.GET("/players/:id/injuries", player.GetInjuries)

		// Standings endpoints.
		protected.GET("/standings", standings.GetStandings)
		protected.GET("/standings/:league/top-scorers", standings.GetTopScorers)

		// User subscriptions & notification preferences.
		protected.GET("/users/me/subscriptions", user.ListSubscriptions)
		protected.POST("/users/me/subscriptions", user.AddSubscription)
		protected.DELETE("/users/me/subscriptions/:clubID", user.RemoveSubscription)
		protected.GET("/users/me/notification-preferences", user.GetNotificationPrefs)
		protected.PUT("/users/me/notification-preferences", user.UpdateNotificationPrefs)
	}
}
