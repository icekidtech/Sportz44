package handlers

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/icekidtech/Sportz44/backend/internal/api/cookies"
	"github.com/icekidtech/Sportz44/backend/internal/api/middleware"
	"github.com/icekidtech/Sportz44/backend/internal/services"
)

// AuthHandler exposes the authentication endpoints.
type AuthHandler struct {
	svc        *services.AuthService
	cookieCfg  cookies.CookieConfig
	accessTTL  time.Duration
	refreshTTL time.Duration
}

// NewAuthHandler creates an AuthHandler.
func NewAuthHandler(svc *services.AuthService, cookieCfg cookies.CookieConfig, accessTTL, refreshTTL time.Duration) *AuthHandler {
	return &AuthHandler{
		svc:        svc,
		cookieCfg:  cookieCfg,
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
	}
}

type registerRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=72"`
}

// Register creates a new account and starts a session via HttpOnly cookies.
func (h *AuthHandler) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, access, refresh, err := h.svc.Register(services.RegisterInput{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		switch {
		case errors.Is(err, services.ErrUserExists):
			c.JSON(http.StatusConflict, gin.H{"error": "username or email already in use"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to register"})
		}
		return
	}

	cookies.SetAuthCookies(c, h.cookieCfg, access, refresh, h.accessTTL, h.refreshTTL)
	c.JSON(http.StatusCreated, gin.H{"user": user})
}

type loginRequest struct {
	Identifier string `json:"identifier" binding:"required"`
	Password   string `json:"password" binding:"required"`
}

// Login verifies credentials and starts a session via HttpOnly cookies.
func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, access, refresh, err := h.svc.Login(req.Identifier, req.Password)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to login"})
		return
	}

	cookies.SetAuthCookies(c, h.cookieCfg, access, refresh, h.accessTTL, h.refreshTTL)
	c.JSON(http.StatusOK, gin.H{"user": user})
}

// Refresh rotates the session using the HttpOnly refresh cookie.
func (h *AuthHandler) Refresh(c *gin.Context) {
	refreshToken, err := c.Cookie(cookies.RefreshCookieName)
	if err != nil || refreshToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing refresh token"})
		return
	}

	user, access, refresh, err := h.svc.Refresh(refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired refresh token"})
		return
	}

	cookies.SetAuthCookies(c, h.cookieCfg, access, refresh, h.accessTTL, h.refreshTTL)
	c.JSON(http.StatusOK, gin.H{"user": user})
}

// Logout clears the session cookies.
func (h *AuthHandler) Logout(c *gin.Context) {
	cookies.ClearAuthCookies(c, h.cookieCfg)
	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
}

// Me returns the currently authenticated user (requires auth middleware).
func (h *AuthHandler) Me(c *gin.Context) {
	userID, ok := c.Get(middleware.ContextUserID)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	user, err := h.svc.GetUser(userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": user})
}
