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
	Environment    string
	HTTPPort       string
	DatabaseURL    string
	RedisURL       string
	APISportsKey   string
	APISportsHost  string
	JWTSecret      string
	JWTExpiry      time.Duration
	RefreshExpiry  time.Duration
	AllowedOrigins []string
	AdminEmail     string
	AdminPassword  string
}

// Load reads configuration from the environment (and an optional .env file).
func Load() (*Config, error) {
	_ = godotenv.Load()

	c := &Config{
		Environment:    getEnv("ENVIRONMENT", ""),
		HTTPPort:       getEnv("HTTP_PORT", ""),
		DatabaseURL:    getEnv("DATABASE_URL", ""),
		RedisURL:       getEnv("REDIS_URL", ""),
		APISportsKey:   os.Getenv("API_SPORTS_KEY"),
		APISportsHost:  getEnv("API_SPORTS_HOST", ""),
		JWTSecret:      os.Getenv("JWT_SECRET"),
		AllowedOrigins: parseOrigins(getEnv("ALLOWED_ORIGINS", "*")),
		AdminEmail:     os.Getenv("ADMIN_EMAIL"),
		AdminPassword:  os.Getenv("ADMIN_PASSWORD"),
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
