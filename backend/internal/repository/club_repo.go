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

// ClubRepo handles DB operations for clubs.
type ClubRepo struct {
	db *gorm.DB
}

// NewClubRepo creates a new ClubRepo.
func NewClubRepo(db *gorm.DB) *ClubRepo {
	return &ClubRepo{db: db}
}

// UpsertClubs inserts or updates the given clubs, keyed by (provider,
// external_id). Returns a map of providerID -> internal ID.
func (r *ClubRepo) UpsertClubs(ctx context.Context, clubs []external.Club) (map[string]uint, error) {
	ids := make(map[string]uint, len(clubs))
	for _, c := range clubs {
		apiID, _ := strconv.Atoi(c.ProviderID)
		m := models.Club{
			Provider:    c.Provider,
			ExternalID:  c.ProviderID,
			APISportsID: apiID,
			Name:        c.Name,
			ShortName:   c.ShortName,
			Country:     c.Country,
			LogoURL:     c.LogoURL,
			Stadium:     c.Stadium,
			Colors:      c.Colors,
			Founded:     c.Founded,
		}
		err := r.db.WithContext(ctx).
			Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "provider"}, {Name: "external_id"}},
				DoUpdates: clause.AssignmentColumns([]string{"name", "short_name", "country", "logo_url", "stadium", "colors", "founded"}),
			}).
			Create(&m).Error
		if err != nil {
			return nil, fmt.Errorf("upsert club %s/%s: %w", c.Provider, c.ProviderID, err)
		}
		ids[c.ProviderID] = m.ID
	}
	return ids, nil
}

// FindByProviderExternal returns the internal ID for a (provider, externalID)
// pair, or 0 if not found.
func (r *ClubRepo) FindByProviderExternal(ctx context.Context, provider, externalID string) (uint, error) {
	var c models.Club
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
