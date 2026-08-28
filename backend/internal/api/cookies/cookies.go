package cookies

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Cookie names used for the auth session.
const (
	AccessCookieName  = "sportz44_access"
	RefreshCookieName = "sportz44_refresh"
)

// CookieConfig controls how auth cookies are written.
type CookieConfig struct {
	Secure   bool // true in production (HTTPS only)
	Domain   string
	SameSite http.SameSite
}

// SetAuthCookies writes the access + refresh tokens as HttpOnly cookies.
//
// Security notes:
//   - HttpOnly: tokens are invisible to JavaScript (XSS-safe).
//   - The refresh cookie is scoped to /api/auth so it is only sent to
//     auth endpoints, reducing the surface for token theft.
//   - Secure is enabled in production so cookies only travel over HTTPS.
func SetAuthCookies(c *gin.Context, cfg CookieConfig, access, refresh string, accessTTL, refreshTTL time.Duration) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     AccessCookieName,
		Value:    access,
		Path:     "/",
		MaxAge:   int(accessTTL.Seconds()),
		HttpOnly: true,
		Secure:   cfg.Secure,
		SameSite: cfg.SameSite,
		Domain:   cfg.Domain,
	})
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     RefreshCookieName,
		Value:    refresh,
		Path:     "/api/auth",
		MaxAge:   int(refreshTTL.Seconds()),
		HttpOnly: true,
		Secure:   cfg.Secure,
		SameSite: cfg.SameSite,
		Domain:   cfg.Domain,
	})
}

// ClearAuthCookies expires both auth cookies (used on logout).
func ClearAuthCookies(c *gin.Context, cfg CookieConfig) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     AccessCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   cfg.Secure,
		SameSite: cfg.SameSite,
		Domain:   cfg.Domain,
	})
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     RefreshCookieName,
		Value:    "",
		Path:     "/api/auth",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   cfg.Secure,
		SameSite: cfg.SameSite,
		Domain:   cfg.Domain,
	})
}
