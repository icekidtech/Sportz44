package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/icekidtech/Sportz44/backend/internal/api/cookies"
	"github.com/icekidtech/Sportz44/backend/pkg/jwt"
)

// Context keys set by the auth middleware.
const (
	ContextUserID = "userID"
	ContextRole   = "role"
)

// Auth validates the access token from the HttpOnly cookie and injects the
// user ID + role into the request context.
func Auth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr, err := c.Cookie(cookies.AccessCookieName)
		if err != nil || tokenStr == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		claims, err := jwt.Parse(tokenStr, secret)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}

		c.Set(ContextUserID, claims.UserID)
		c.Set(ContextRole, claims.Role)
		c.Next()
	}
}

// RequireRole restricts a route to users without one of the given roles.
func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, ok := c.Get(ContextRole)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		for _, r := range roles {
			if role == r {
				c.Next()
				return
			}
		}
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
	}
}
