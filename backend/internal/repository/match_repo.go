package repository

import (
	"context"
	"fmt"
	"strconv"

	"github.com/icekidtech/Sportz44/backend/internal/external"
	"github.com/icekidtech/Sportz44/backend/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// MatchRepo handles DB operations for matches.
type MatchRepo struct {
	db *gorm.DB
}

// NewMatchRepo creates a new MatchRepo.
func NewMatchRepo(db *gorm.DB) *MatchRepo {
	return &MatchRepo{db: db}
}

// UpsertMatches inserts or updates the given fixtures by (provider,
// external_id). The competitionID and club IDs are resolved from the provider
// data via the provided lookup functions, so matches get proper FKs.
func (r *MatchRepo) UpsertMatches(
	ctx context.Context,
	fixtures []external.Fixture,
	competitionID uint,
	clubIDFor func(provider, externalID string) uint,
) error {
	for _, f := range fixtures {
		apiID, _ := strconv.Atoi(f.ProviderID)
		homeID := clubIDFor(f.Provider, f.HomeTeamID)
		awayID := clubIDFor(f.Provider, f.AwayTeamID)
		m := models.Match{
			Provider:      f.Provider,
			ExternalID:    f.ProviderID,
			APISportsID:   apiID,
			CompetitionID: competitionID,
			Season:        f.Season,
			HomeClubID:    homeID,
			AwayClubID:    awayID,
			MatchDate:     f.MatchDate,
			Status:        f.Status,
			HomeScore:     f.HomeScore,
			AwayScore:     f.AwayScore,
			Venue:         f.Venue,
			Referee:       f.Referee,
			Minute:        f.Minute,
		}
		err := r.db.WithContext(ctx).
			Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "provider"}, {Name: "external_id"}},
				DoUpdates: clause.AssignmentColumns([]string{"status", "home_score", "away_score", "minute", "match_date", "competition_id", "home_club_id", "away_club_id"}),
			}).
			Create(&m).Error
		if err != nil {
			return fmt.Errorf("upsert match %s/%s: %w", f.Provider, f.ProviderID, err)
		}
	}
	return nil
}
