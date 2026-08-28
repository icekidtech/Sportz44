package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"github.com/icekidtech/Sportz44/backend/internal/external"
	"github.com/icekidtech/Sportz44/backend/internal/repository"
	"github.com/icekidtech/Sportz44/backend/pkg/cache"
	"github.com/icekidtech/Sportz44/backend/pkg/config"
	"github.com/icekidtech/Sportz44/backend/pkg/database"
	"github.com/icekidtech/Sportz44/backend/pkg/logger"
	"github.com/sirupsen/logrus"
)

// match-listener polls the realtime provider for live matches, diffs events,
// upserts them to the database, and publishes updates to Redis so the
// WebSocket hub can broadcast them to clients.
//
// Run with:
//
//	go run ./cmd/match-listener
func main() {
	_ = godotenv.Load()
	cfg, err := config.Load()
	if err != nil {
		logger.New("").Fatalf("config: %v", err)
	}
	log := logger.New(cfg.Environment)

	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	rdb, err := cache.New(cfg.RedisURL)
	if err != nil {
		log.Fatalf("redis: %v", err)
	}
	if err := rdb.Ping(context.Background()); err != nil {
		log.Warnf("redis ping failed (continuing): %v", err)
	}

	// Providers & repos.
	apiSports := external.NewAPISportsProvider(cfg.APISportsKey, cfg.APISportsHost)
	flashScore := external.NewFlashScoreProvider(cfg.FlashScoreURL)
	reg := external.NewRegistry(apiSports, flashScore)
	matchRepo := repository.NewMatchRepo(db)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Graceful shutdown on SIGINT/SIGTERM.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-quit
		log.Info("shutting down match-listener...")
		cancel()
	}()

	log.Info("match-listener started (polling every 30s)")

	// Poll loop.
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			poll(ctx, reg, matchRepo, rdb, log)
		}
	}
}

// poll fetches live matches and their events, upserts them, and publishes.
func poll(ctx context.Context, reg *external.Registry, matchRepo *repository.MatchRepo, rdb *cache.Redis, log *logrus.Logger) {
	lp := reg.Live(ctx)
	if lp == nil {
		log.Warn("no healthy live provider")
		return
	}

	live, err := lp.GetLiveMatches(ctx)
	if err != nil {
		log.Warnf("get live matches: %v", err)
		return
	}
	if len(live) == 0 {
		return
	}
	log.Infof("found %d live matches", len(live))

	for _, f := range live {
		// Resolve the internal match by provider + external ID.
		m, err := matchRepo.FindByExternal(ctx, f.Provider, f.ProviderID)
		if err != nil {
			log.Warnf("match %s/%s not in DB (skipping): %v", f.Provider, f.ProviderID, err)
			continue
		}

		// Fetch and upsert events.
		events, err := lp.GetLiveEvents(ctx, f.ProviderID)
		if err != nil {
			log.Warnf("get live events for %s: %v", f.ProviderID, err)
			continue
		}
		if err := matchRepo.UpsertMatchEvents(ctx, m.ID, events); err != nil {
			log.Warnf("upsert events for match %d: %v", m.ID, err)
			continue
		}

		// Publish the update to Redis for the WebSocket hub.
		update := map[string]interface{}{
			"match_id":    m.ID,
			"provider":    f.Provider,
			"external_id": f.ProviderID,
			"status":      f.Status,
			"home_score":  f.HomeScore,
			"away_score":  f.AwayScore,
			"minute":      f.Minute,
			"events":      events,
			"timestamp":   time.Now().UTC(),
		}
		if err := rdb.Publish(ctx, "match-events", update); err != nil {
			log.Warnf("publish update for match %d: %v", m.ID, err)
		}
	}
}
