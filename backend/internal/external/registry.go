package external

import (
	"context"
	"errors"
	"sync"
)

// Registry holds all configured providers and provides failover-aware access
// by role (bulk, realtime, media). It lets the ingestion service pick the best
// available provider for a task and fall back when one is rate-limited or down.
type Registry struct {
	mu        sync.RWMutex
	providers []Provider
}

// NewRegistry creates a registry from the given providers.
func NewRegistry(providers ...Provider) *Registry {
	return &Registry{providers: providers}
}

// All returns every registered provider.
func (r *Registry) All() []Provider {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]Provider, len(r.providers))
	copy(out, r.providers)
	return out
}

// ByKind returns the first healthy provider of the given kind, or nil.
func (r *Registry) ByKind(ctx context.Context, kind ProviderKind) Provider {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, p := range r.providers {
		if p.Kind() == kind && p.Healthy(ctx) {
			return p
		}
	}
	return nil
}

// Fixtures returns a FixtureProvider of the given kind (falling back to any
// healthy fixture provider if the preferred kind is unavailable).
func (r *Registry) Fixtures(ctx context.Context, kind ProviderKind) FixtureProvider {
	if p := r.ByKind(ctx, kind); p != nil {
		if fp, ok := p.(FixtureProvider); ok {
			return fp
		}
	}
	// Fallback: any healthy fixture provider.
	for _, p := range r.All() {
		if fp, ok := p.(FixtureProvider); ok && p.Healthy(ctx) {
			return fp
		}
	}
	return nil
}

// FixturesForSeason returns a healthy FixtureProvider that can serve the given
// season, preferring providers that explicitly support it (e.g. Football-Data
// for the current season, API-Sports for 2022-2024). Falls back to any healthy
// fixture provider if none declares season support.
func (r *Registry) FixturesForSeason(ctx context.Context, season string) FixtureProvider {
	for _, p := range r.All() {
		if sp, ok := p.(SeasonProvider); ok && sp.SupportsSeason(season) {
			if fp, ok := p.(FixtureProvider); ok && p.Healthy(ctx) {
				return fp
			}
		}
	}
	return r.Fixtures(ctx, ProviderBulk)
}

// CompetitionsForSeason returns a healthy CompetitionProvider that can serve
// the given season, preferring providers that explicitly support it.
func (r *Registry) CompetitionsForSeason(ctx context.Context, season string) CompetitionProvider {
	for _, p := range r.All() {
		if sp, ok := p.(SeasonProvider); ok && sp.SupportsSeason(season) {
			if cp, ok := p.(CompetitionProvider); ok && p.Healthy(ctx) {
				return cp
			}
		}
	}
	return r.Competitions(ctx)
}

// Live returns a healthy realtime LiveProvider, or nil.
func (r *Registry) Live(ctx context.Context) LiveProvider {
	if p := r.ByKind(ctx, ProviderRealtime); p != nil {
		if lp, ok := p.(LiveProvider); ok {
			return lp
		}
	}
	return nil
}

// Squad returns a healthy SquadProvider, or nil.
func (r *Registry) Squad(ctx context.Context) SquadProvider {
	if p := r.ByKind(ctx, ProviderBulk); p != nil {
		if sp, ok := p.(SquadProvider); ok {
			return sp
		}
	}
	return nil
}

// Stats returns a healthy StatsProvider, or nil.
func (r *Registry) Stats(ctx context.Context) StatsProvider {
	for _, p := range r.All() {
		if sp, ok := p.(StatsProvider); ok && p.Healthy(ctx) {
			return sp
		}
	}
	return nil
}

// Competitions returns a healthy CompetitionProvider, or nil.
func (r *Registry) Competitions(ctx context.Context) CompetitionProvider {
	if p := r.ByKind(ctx, ProviderBulk); p != nil {
		if cp, ok := p.(CompetitionProvider); ok {
			return cp
		}
	}
	return nil
}

// Media returns a healthy MediaProvider, or nil.
func (r *Registry) Media(ctx context.Context) MediaProvider {
	if p := r.ByKind(ctx, ProviderMedia); p != nil {
		if mp, ok := p.(MediaProvider); ok {
			return mp
		}
	}
	return nil
}

// IsRateLimited reports whether an error indicates a provider quota/rate limit.
func IsRateLimited(err error) bool {
	return errors.Is(err, ErrRateLimited)
}
