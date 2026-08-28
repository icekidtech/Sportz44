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

// CompetitionRepo handles DB operations for competitions.
type CompetitionRepo struct {
	db *gorm.DB
}

// NewCompetitionRepo creates a new CompetitionRepo.
func NewCompetitionRepo(db *gorm.DB) *CompetitionRepo {
	return &CompetitionRepo{db: db}
}

// UpsertCompetitions inserts or updates the given competitions, keyed by
// (provider, external_id). Returns a map of providerID -> internal ID.
func (r *CompetitionRepo) UpsertCompetitions(ctx context.Context, comps []external.Competition) (map[string]uint, error) {
	ids := make(map[string]uint, len(comps))
	for _, c := range comps {
		apiID, _ := strconv.Atoi(c.ProviderID)
		m := models.Competition{
			Provider:    c.Provider,
			ExternalID:  c.ProviderID,
			APISportsID: apiID,
			Name:        c.Name,
			Type:        c.Type,
			Country:     c.Country,
			LogoURL:     c.LogoURL,
			Season:      c.Season,
		}
		err := r.db.WithContext(ctx).
			Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "provider"}, {Name: "external_id"}},
				DoUpdates: clause.AssignmentColumns([]string{"name", "type", "country", "logo_url", "season"}),
			}).
			Create(&m).Error
		if err != nil {
			return nil, fmt.Errorf("upsert competition %s/%s: %w", c.Provider, c.ProviderID, err)
		}
		ids[c.ProviderID] = m.ID
	}
	return ids, nil
}

// FindByProviderExternal returns the internal ID for a (provider, externalID)
// pair, or 0 if not found.
func (r *CompetitionRepo) FindByProviderExternal(ctx context.Context, provider, externalID string) (uint, error) {
	var c models.Competition
	err := r.db.WithContext(ctx).
		Where("provider = ? AND external_id = ?", provider, externalID).
		First(&c).Error
	if err == gorm.ErrRecordNotFound {
		return 0, nil
	}
	if err != nil {
		return 0, err
	}
	return c.ID, nil
}
