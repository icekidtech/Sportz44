package services

import (
	"context"

	"github.com/icekidtech/Sportz44/backend/internal/repository"
)

// StandingsService computes league tables and top scorers.
type StandingsService struct {
	standings *repository.StandingsRepo
}

// NewStandingsService creates a new StandingsService.
func NewStandingsService(standings *repository.StandingsRepo) *StandingsService {
	return &StandingsService{standings: standings}
}

// GetStandings returns the league table for a competition, optionally filtered by season.
func (s *StandingsService) GetStandings(ctx context.Context, competitionID uint, season string) ([]repository.StandingRow, error) {
	return s.standings.GetStandings(ctx, competitionID, season)
}

// GetTopScorers returns the golden-boot standings for a competition, optionally filtered by season.
func (s *StandingsService) GetTopScorers(ctx context.Context, competitionID uint, season string, limit int) ([]repository.TopScorer, error) {
	return s.standings.GetTopScorers(ctx, competitionID, season, limit)
}
