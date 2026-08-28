package services

import (
	"context"
	"errors"

	"github.com/icekidtech/Sportz44/backend/internal/external"
	"github.com/icekidtech/Sportz44/backend/internal/models"
	"github.com/icekidtech/Sportz44/backend/internal/repository"
)

// PlayerService orchestrates squad ingestion and player queries.
type PlayerService struct {
	players  *repository.PlayerRepo
	clubs    *repository.ClubRepo
	registry *external.Registry
}

// NewPlayerService creates a new PlayerService.
func NewPlayerService(players *repository.PlayerRepo, clubs *repository.ClubRepo, registry *external.Registry) *PlayerService {
	return &PlayerService{players: players, clubs: clubs, registry: registry}
}

// SyncSquad fetches a club's squad from the provider and upserts it.
func (s *PlayerService) SyncSquad(ctx context.Context, clubID uint) error {
	club, err := s.clubs.FindByID(ctx, clubID)
	if err != nil {
		return err
	}
	sp := s.registry.Squad(ctx)
	if sp == nil {
		return errors.New("no squad provider available")
	}
	squad, err := sp.GetSquad(ctx, club.ExternalID)
	if err != nil {
		return err
	}
	return s.players.UpsertPlayers(ctx, squad, clubID)
}

// ListByClub returns the roster for a club.
func (s *PlayerService) ListByClub(ctx context.Context, clubID uint) ([]models.Player, error) {
	return s.players.ListByClub(ctx, clubID)
}

// Get returns a single player.
func (s *PlayerService) Get(ctx context.Context, id uint) (*models.Player, error) {
	return s.players.FindByID(ctx, id)
}

// GetSeasonStats returns a player's season statistics.
func (s *PlayerService) GetSeasonStats(ctx context.Context, id uint) ([]models.PlayerSeasonStats, error) {
	return s.players.GetSeasonStats(ctx, id)
}

// GetForm returns a player's recent form.
func (s *PlayerService) GetForm(ctx context.Context, id uint, limit int) ([]repository.PlayerFormEntry, error) {
	return s.players.GetForm(ctx, id, limit)
}

// GetInjuries returns a player's injuries.
func (s *PlayerService) GetInjuries(ctx context.Context, id uint) ([]models.Injury, error) {
	return s.players.GetInjuries(ctx, id)
}
