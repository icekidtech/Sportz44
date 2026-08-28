package services

import (
	"context"
	"github.com/icekidtech/Sportz44/backend/internal/external"
	"github.com/icekidtech/Sportz44/backend/internal/repository"
)

// MatchService orchestrates match ingestion.
type MatchService struct {
	repo     *repository.MatchRepo
	registry *external.Registry
}

// NewMatchService creates a new MatchService.
func NewMatchService(repo *repository.MatchRepo, registry *external.Registry) *MatchService {
	return &MatchService{repo: repo, registry: registry}
}

// SyncCompetition fetches fixtures for the given competition and season,
// and upserts them into the database.
func (s *MatchService) SyncCompetition(ctx context.Context, competitionID, season string) error {
	fp := s.registry.FixturesForSeason(ctx, season)
	fixtures, err := fp.GetFixtures(ctx, competitionID, season, "")
	if err != nil {
		return err
	}
	return s.repo.UpsertMatches(ctx, fixtures)
}