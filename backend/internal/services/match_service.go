package services

import (
	"context"
	"fmt"

	"github.com/icekidtech/Sportz44/backend/internal/external"
	"github.com/icekidtech/Sportz44/backend/internal/repository"
)

// MatchService orchestrates match ingestion.
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
				Provider:   provider,
				ProviderID: f.HomeTeamID,
				Name:       f.HomeTeamName,
			})
		}
		if !seen[f.AwayTeamID] {
			seen[f.AwayTeamID] = true
			clubs = append(clubs, external.Club{
				Provider:   provider,
				ProviderID: f.AwayTeamID,
				Name:       f.AwayTeamName,
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
