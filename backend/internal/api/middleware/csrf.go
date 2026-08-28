package middleware

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CSRFCookieName is the (non-HttpOnly) double-submit cookie.
const CSRFCookieName = "sportz44_csrf"

// CSRFProtect implements the double-submit cookie pattern.
//
// Because auth uses HttpOnly cookies (sent automatically by the browser),
// state-changing requests need CSRF protection. A random token is stored in
// a readable cookie; the client must echo it back in the X-CSRF-Token header.
// Safe methods (GET/HEAD/OPTIONS/TRACE) are not validated.
func CSRFProtect() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := c.Cookie(CSRFCookieName)
		if err != nil || token == "" {
			token = generateCSRFToken()
			http.SetCookie(c.Writer, &http.Cookie{
				Name:     CSRFCookieName,
				Value:    token,
				Path:     "/",
				MaxAge:   86400 * 7,
				HttpOnly: false, // client JS must be able to read it
				Secure:   c.Request.TLS != nil,
				SameSite: http.SameSiteLaxMode,
			})
		}

		switch c.Request.Method {
		case http.MethodGet, http.MethodHead, http.MethodOptions, http.MethodTrace:
			c.Next()
			return
		}

		if headerToken := c.GetHeader("X-CSRF-Token"); headerToken == "" || headerToken != token {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "invalid CSRF token"})
			return
		}
		c.Next()
	}
}

func generateCSRFToken() string {
	b := make([]byte, 32)
	_, _ = rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}
