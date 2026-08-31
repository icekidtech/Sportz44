package services

import (
	"context"
	"fmt"

	"github.com/icekidtech/Sportz44/backend/internal/external"
	"github.com/icekidtech/Sportz44/backend/internal/models"
	"github.com/icekidtech/Sportz44/backend/internal/repository"
)

// MatchService orchestrates match ingestion and querying.
type MatchService struct {
	matches      *repository.MatchRepo
	competitions *repository.CompetitionRepo
	clubs        *repository.ClubRepo
	registry     *external.Registry
}

// NewMatchService creates a new MatchService.
func NewMatchService(
	matches *repository.MatchRepo,
	competitions *repository.CompetitionRepo,
	clubs *repository.ClubRepo,
	registry *external.Registry,
) *MatchService {
	return &MatchService{
		matches:      matches,
		competitions: competitions,
		clubs:        clubs,
		registry:     registry,
	}
}

// SyncCompetition ingests a competition's fixtures for a season end-to-end:
//  1. Upsert the competition metadata.
//  2. Fetch fixtures from the season-appropriate provider.
//  3. Upsert the clubs referenced by those fixtures.
//  4. Upsert the matches with resolved foreign keys.
func (s *MatchService) SyncCompetition(ctx context.Context, competitionID, season string) error {
	// 1. Pick the provider for this season and fetch fixtures.
	fp := s.registry.FixturesForSeason(ctx, season)
	if fp == nil {
		return fmt.Errorf("no healthy fixture provider for season %s", season)
	}
	fixtures, err := fp.GetFixtures(ctx, competitionID, season, "")
	if err != nil {
		return fmt.Errorf("fetch fixtures: %w", err)
	}
	if len(fixtures) == 0 {
		return fmt.Errorf("no fixtures returned for competition %s season %s", competitionID, season)
	}
	// Stamp the season onto each fixture so matches carry it (providers may
	// not populate it themselves).
	for i := range fixtures {
		fixtures[i].Season = season
	}

	// 2. Upsert the competition metadata and resolve its internal ID.
	provider := fixtures[0].Provider
	comp := external.Competition{
		Provider:   provider,
		ProviderID: competitionID,
		Name:       competitionID,
		Season:     season,
	}
	compIDs, err := s.competitions.UpsertCompetitions(ctx, []external.Competition{comp})
	if err != nil {
		return fmt.Errorf("upsert competition: %w", err)
	}
	internalCompID := compIDs[competitionID]

	// 3. Upsert the clubs referenced by the fixtures.
	clubs := make([]external.Club, 0, len(fixtures)*2)
	seen := map[string]bool{}
	for _, f := range fixtures {
		if !seen[f.HomeTeamID] {
			seen[f.HomeTeamID] = true
			clubs = append(clubs, external.Club{
				Provider:      provider,
				ProviderID:    f.HomeTeamID,
				Name:          f.HomeTeamName,
				LogoURL:       f.HomeTeamLogo,
				CompetitionID: int(internalCompID),
			})
		}
		if !seen[f.AwayTeamID] {
			seen[f.AwayTeamID] = true
			clubs = append(clubs, external.Club{
				Provider:      provider,
				ProviderID:    f.AwayTeamID,
				Name:          f.AwayTeamName,
				LogoURL:       f.AwayTeamLogo,
				CompetitionID: int(internalCompID),
			})
		}
	}
	clubIDs, err := s.clubs.UpsertClubs(ctx, clubs)
	if err != nil {
		return fmt.Errorf("upsert clubs: %w", err)
	}

	// 4. Upsert matches with resolved FKs.
	clubIDFor := func(prov, externalID string) uint {
		return clubIDs[externalID]
	}
	if err := s.matches.UpsertMatches(ctx, fixtures, internalCompID, clubIDFor); err != nil {
		return fmt.Errorf("upsert matches: %w", err)
	}
	return nil
}

// ---- Query methods ----

// ListMatches returns paginated matches matching the given filters.
func (s *MatchService) ListMatches(ctx context.Context, f repository.MatchFilters) ([]models.Match, int64, error) {
	return s.matches.ListMatches(ctx, f)
}

// GetMatch returns a single match by ID.
func (s *MatchService) GetMatch(ctx context.Context, id uint) (*models.Match, error) {
	return s.matches.GetMatch(ctx, id)
}

// GetMatchEvents returns events for a match.
func (s *MatchService) GetMatchEvents(ctx context.Context, matchID uint) ([]models.MatchEvent, error) {
	return s.matches.GetMatchEvents(ctx, matchID)
}

// GetMatchLineup returns the lineup for a match.
func (s *MatchService) GetMatchLineup(ctx context.Context, matchID uint) ([]models.MatchLineup, error) {
	return s.matches.GetMatchLineup(ctx, matchID)
}

// GetMatchStats returns the statistics for a match.
func (s *MatchService) GetMatchStats(ctx context.Context, matchID uint) ([]models.MatchStat, error) {
	return s.matches.GetMatchStats(ctx, matchID)
}

// SyncMatchEvents fetches and upserts the events for a single match. This is
// used to backfill finished matches — the live listener only polls matches
// that are currently live. Events are fetched from the same provider that
// ingested the match, so the provider-specific external ID is always valid.
func (s *MatchService) SyncMatchEvents(ctx context.Context, matchID uint) error {
	m, err := s.matches.GetMatch(ctx, matchID)
	if err != nil {
		return err
	}
	// Use the provider that ingested the match (no health probe needed — we
	// already know it's the right one, and a per-match /status call would be
	// wasteful and flaky during bulk backfill).
	p := s.registry.ByName(m.Provider)
	lp, ok := p.(external.EventsProvider)
	if !ok || lp == nil {
		return fmt.Errorf("provider %q does not support events or is unavailable", m.Provider)
	}
	events, err := lp.GetLiveEvents(ctx, m.ExternalID)
	if err != nil {
		return fmt.Errorf("fetch events: %w", err)
	}
	return s.matches.UpsertMatchEvents(ctx, m.ID, events)
}

// SyncMatchLineup fetches and upserts the lineup for a single match.
func (s *MatchService) SyncMatchLineup(ctx context.Context, matchID uint) error {
	m, err := s.matches.GetMatch(ctx, matchID)
	if err != nil {
		return err
	}
	p := s.registry.ByName(m.Provider)
	lp, ok := p.(external.LineupProvider)
	if !ok || lp == nil {
		return fmt.Errorf("provider %q does not support lineups or is unavailable", m.Provider)
	}
	lineup, err := lp.GetMatchLineup(ctx, m.ExternalID)
	if err != nil {
		return fmt.Errorf("fetch lineup: %w", err)
	}
	return s.matches.UpsertMatchLineup(ctx, m.ID, lineup)
}

// SyncMatchStats fetches and upserts the statistics for a single match.
func (s *MatchService) SyncMatchStats(ctx context.Context, matchID uint) error {
	m, err := s.matches.GetMatch(ctx, matchID)
	if err != nil {
		return err
	}
	p := s.registry.ByName(m.Provider)
	sp, ok := p.(external.StatsProvider)
	if !ok || sp == nil {
		return fmt.Errorf("provider %q does not support stats or is unavailable", m.Provider)
	}
	stats, err := sp.GetMatchStats(ctx, m.ExternalID)
	if err != nil {
		return fmt.Errorf("fetch stats: %w", err)
	}
	return s.matches.UpsertMatchStats(ctx, m.ID, stats)
}

// SyncMatchDetails fetches and upserts events, lineup, and stats for a single
// match in one call. Each sub-fetch is best-effort — a failure in one does
// not abort the others.
func (s *MatchService) SyncMatchDetails(ctx context.Context, matchID uint) error {
	m, err := s.matches.GetMatch(ctx, matchID)
	if err != nil {
		return err
	}
	p := s.registry.ByName(m.Provider)
	if p == nil {
		return fmt.Errorf("provider %q not found", m.Provider)
	}
	var lastErr error
	if ep, ok := p.(external.EventsProvider); ok {
		if events, err := ep.GetLiveEvents(ctx, m.ExternalID); err == nil {
			if err := s.matches.UpsertMatchEvents(ctx, m.ID, events); err != nil {
				lastErr = err
			}
		} else {
			lastErr = fmt.Errorf("fetch events: %w", err)
		}
	}
	if lp, ok := p.(external.LineupProvider); ok {
		if lineup, err := lp.GetMatchLineup(ctx, m.ExternalID); err == nil {
			if err := s.matches.UpsertMatchLineup(ctx, m.ID, lineup); err != nil {
				lastErr = err
			}
		} else {
			lastErr = fmt.Errorf("fetch lineup: %w", err)
		}
	}
	if sp, ok := p.(external.StatsProvider); ok {
		if stats, err := sp.GetMatchStats(ctx, m.ExternalID); err == nil {
			if err := s.matches.UpsertMatchStats(ctx, m.ID, stats); err != nil {
				lastErr = err
			}
		} else {
			lastErr = fmt.Errorf("fetch stats: %w", err)
		}
	}
	return lastErr
}
