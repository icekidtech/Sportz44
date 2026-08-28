package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/icekidtech/Sportz44/backend/internal/api"
	"github.com/icekidtech/Sportz44/backend/internal/api/cookies"
	"github.com/icekidtech/Sportz44/backend/internal/api/handlers"
	"github.com/icekidtech/Sportz44/backend/internal/api/middleware"
	"github.com/icekidtech/Sportz44/backend/internal/repository"
	"github.com/icekidtech/Sportz44/backend/internal/services"
	"github.com/icekidtech/Sportz44/backend/pkg/cache"
	"github.com/icekidtech/Sportz44/backend/pkg/config"
	"github.com/icekidtech/Sportz44/backend/pkg/database"
	"github.com/icekidtech/Sportz44/backend/pkg/logger"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		logger.New("").Fatalf("config: %v", err)
	}
	log := logger.New(cfg.Environment)

	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	if err := database.AutoMigrate(db); err != nil {
		log.Fatalf("migrate: %v", err)
	}
	log.Info("database connected and migrated")

	rdb, err := cache.New(cfg.RedisURL)
	if err != nil {
		log.Fatalf("redis: %v", err)
	}
	if err := rdb.Ping(context.Background()); err != nil {
		log.Warnf("redis ping failed (continuing): %v", err)
	}

	// Repositories & services.
	userRepo := repository.NewUserRepo(db)
	authSvc := services.NewAuthService(userRepo, cfg.JWTSecret, cfg.JWTExpiry, cfg.RefreshExpiry)

	cookieCfg := cookies.CookieConfig{
		Secure:   cfg.Environment == "production",
		Domain:   cfg.CookieDomain,
		SameSite: http.SameSiteLaxMode,
	}
	authHandler := handlers.NewAuthHandler(authSvc, cookieCfg, cfg.JWTExpiry, cfg.RefreshExpiry)
	healthHandler := handlers.NewHealthHandler(db, rdb)

	// Router.
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(middleware.Logging(log), gin.Recovery(), middleware.CORS(cfg.AllowedOrigins))
	api.RegisterRoutes(r, authHandler, healthHandler, cfg.JWTSecret)

	srv := &http.Server{
		Addr:    ":" + cfg.HTTPPort,
		Handler: r,
	}

	go func() {
		log.Infof("Sportz44 API listening on :%s", cfg.HTTPPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server: %v", err)
		}
	}()

	// Graceful shutdown.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Errorf("shutdown: %v", err)
	}
}
