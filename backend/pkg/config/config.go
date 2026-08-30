package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all runtime configuration loaded from environment / .env.
type Config struct {
	Environment       string
	HTTPPort          string
	DatabaseURL       string
	RedisURL          string
	APISportsKey      string
	APISportsHost     string
	ESPNURL           string
	FootballDataKey   string
	FootballDataURL   string
	JWTSecret         string
	JWTExpiry         time.Duration
	RefreshExpiry     time.Duration
	AllowedOrigins    []string
	TrustedProxies    []string
	CookieDomain      string
	AdminEmail        string
	AdminPassword     string
	RateLimitRequests int
	RateLimitWindow   time.Duration
}

// Load reads configuration from the environment (and an optional .env file).
func Load() (*Config, error) {
	_ = godotenv.Load()

	c := &Config{
		Environment:     getEnv("ENVIRONMENT", ""),
		HTTPPort:        getEnv("HTTP_PORT", ""),
		DatabaseURL:     getEnv("DATABASE_URL", ""),
		RedisURL:        getEnv("REDIS_URL", ""),
		APISportsKey:    os.Getenv("API_SPORTS_KEY"),
		APISportsHost:   getEnv("API_SPORTS_HOST", ""),
		ESPNURL:         getEnv("ESPN_URL", ""),
		FootballDataKey: os.Getenv("FOOTBALL_DATA_KEY"),
		FootballDataURL: getEnv("FOOTBALL_DATA_URL", ""),
		JWTSecret:       os.Getenv("JWT_SECRET"),
		AllowedOrigins:  parseOrigins(getEnv("ALLOWED_ORIGINS", "*")),
		TrustedProxies:  parseProxies(getEnv("TRUSTED_PROXIES", "")),
		CookieDomain:    getEnv("COOKIE_DOMAIN", ""),
		AdminEmail:      os.Getenv("ADMIN_EMAIL"),
		AdminPassword:   os.Getenv("ADMIN_PASSWORD"),
	}

	exp, err := time.ParseDuration(getEnv("JWT_EXPIRY", "15m"))
	if err != nil {
		return nil, fmt.Errorf("invalid JWT_EXPIRY: %w", err)
	}
	c.JWTExpiry = exp

	rexp, err := time.ParseDuration(getEnv("REFRESH_TOKEN_EXPIRY", "168h"))
	if err != nil {
		return nil, fmt.Errorf("invalid REFRESH_TOKEN_EXPIRY: %w", err)
	}
	c.RefreshExpiry = rexp

	rlWindow, err := time.ParseDuration(getEnv("RATE_LIMIT_WINDOW", "1m"))
	if err != nil {
		return nil, fmt.Errorf("invalid RATE_LIMIT_WINDOW: %w", err)
	}
	c.RateLimitWindow = rlWindow
	c.RateLimitRequests = getEnvInt("RATE_LIMIT_REQUESTS", 100)

	if c.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
	}
	if c.APISportsKey == "" {
		return nil, fmt.Errorf("API_SPORTS_KEY is required")
	}
	return c, nil
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func getEnvInt(key string, def int) int {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	n := 0
	if _, err := fmt.Sscanf(v, "%d", &n); err != nil {
		return def
	}
	return n
}

func parseOrigins(s string) []string {
	if s == "*" {
		return []string{"*"}
	}
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

// parseProxies parses a comma-separated list of trusted proxy CIDRs/IPs.
// An empty value means "trust no proxy" (Gin will use the TCP peer address
// directly, ignoring X-Forwarded-For).
func parseProxies(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}
