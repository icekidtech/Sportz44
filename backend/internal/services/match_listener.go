package services

import (
	"context"
	"time"

	"github.com/icekidtech/Sportz44/backend/internal/external"
	"github.com/icekidtech/Sportz44/backend/internal/repository"
	"github.com/icekidtech/Sportz44/backend/pkg/cache"
	"github.com/sirupsen/logrus"
)

// MatchListener polls the realtime provider for live matches, diffs events,
// upserts them to the database, and publishes updates to Redis so the
// WebSocket hub can broadcast them to clients.
//
// It runs as a background goroutine inside the API server process (no separate
// binary needed), sharing the same DB, Redis, and provider registry.
type MatchListener struct {
	registry  *external.Registry
	matchRepo *repository.MatchRepo
	rdb       *cache.Redis
	log       *logrus.Logger
	interval  time.Duration
}

// NewMatchListener creates a MatchListener with the given poll interval.
func NewMatchListener(registry *external.Registry, matchRepo *repository.MatchRepo, rdb *cache.Redis, log *logrus.Logger, interval time.Duration) *MatchListener {
	return &MatchListener{
		registry:  registry,
		matchRepo: matchRepo,
		rdb:       rdb,
		log:       log,
		interval:  interval,
	}
}

// Start launches the polling loop in a goroutine. It stops when ctx is
// cancelled (e.g. on graceful server shutdown).
func (l *MatchListener) Start(ctx context.Context) {
	go func() {
		l.log.Infof("match-listener started (polling every %s)", l.interval)
		ticker := time.NewTicker(l.interval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				l.log.Info("match-listener stopped")
				return
			case <-ticker.C:
				l.poll(ctx)
			}
		}
	}()
}

// poll fetches live matches and their events, upserts them, and publishes.
func (l *MatchListener) poll(ctx context.Context) {
	lp := l.registry.Live(ctx)
	if lp == nil {
		l.log.Warn("no healthy live provider")
		return
	}

	live, err := lp.GetLiveMatches(ctx)
	if err != nil {
		l.log.Warnf("get live matches: %v", err)
		return
	}
	if len(live) == 0 {
		return
	}
	l.log.Infof("found %d live matches", len(live))

	for _, f := range live {
		// Resolve the internal match by provider + external ID.
		m, err := l.matchRepo.FindByExternal(ctx, f.Provider, f.ProviderID)
		if err != nil {
			l.log.Warnf("match %s/%s not in DB (skipping): %v", f.Provider, f.ProviderID, err)
			continue
		}

		// Fetch and upsert events.
		events, err := lp.GetLiveEvents(ctx, f.ProviderID)
		if err != nil {
			l.log.Warnf("get live events for %s: %v", f.ProviderID, err)
			continue
		}
		if err := l.matchRepo.UpsertMatchEvents(ctx, m.ID, events); err != nil {
			l.log.Warnf("upsert events for match %d: %v", m.ID, err)
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
		if err := l.rdb.Publish(ctx, "match-events", update); err != nil {
			l.log.Warnf("publish update for match %d: %v", m.ID, err)
		}
	}
}
