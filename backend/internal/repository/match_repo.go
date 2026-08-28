package repository

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/icekidtech/Sportz44/backend/internal/external"
	"github.com/icekidtech/Sportz44/backend/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// MatchFilters defines query parameters for listing matches.
type MatchFilters struct {
	ClubID       uint   // filter by home or away club
	CompetitionID uint  // filter by competition
	Status       string // scheduled | live | finished
	Season       string
	DateFrom     *time.Time
	DateTo       *time.Time
	Page         int
	Limit        int
}

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

// ---- Query methods ----

// ListMatches returns paginated matches matching the given filters.
func (r *MatchRepo) ListMatches(ctx context.Context, f MatchFilters) ([]models.Match, int64, error) {
	q := r.db.WithContext(ctx).Model(&models.Match{})
	if f.ClubID > 0 {
		q = q.Where("home_club_id = ? OR away_club_id = ?", f.ClubID, f.ClubID)
	}
	if f.CompetitionID > 0 {
		q = q.Where("competition_id = ?", f.CompetitionID)
	}
	if f.Status != "" {
		q = q.Where("status = ?", f.Status)
	}
	if f.Season != "" {
		q = q.Where("season = ?", f.Season)
	}
	if f.DateFrom != nil {
		q = q.Where("match_date >= ?", *f.DateFrom)
	}
	if f.DateTo != nil {
		q = q.Where("match_date <= ?", *f.DateTo)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	page, limit := f.Page, f.Limit
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	offset := (page - 1) * limit
	var matches []models.Match
	err := q.Preload("HomeClub").Preload("AwayClub").Preload("Competition").
		Order("match_date DESC").
		Offset(offset).Limit(limit).
		Find(&matches).Error
	return matches, total, err
}

// GetMatch returns a single match by ID with relations loaded.
func (r *MatchRepo) GetMatch(ctx context.Context, id uint) (*models.Match, error) {
	var m models.Match
	err := r.db.WithContext(ctx).
		Preload("HomeClub").Preload("AwayClub").Preload("Competition").
		First(&m, id).Error
	if err != nil {
		return nil, err
	}
	return &m, nil
}

// GetMatchEvents returns all events for a match, sorted by minute.
func (r *MatchRepo) GetMatchEvents(ctx context.Context, matchID uint) ([]models.MatchEvent, error) {
	var events []models.MatchEvent
	err := r.db.WithContext(ctx).
		Where("match_id = ?", matchID).
		Order("minute ASC").
		Find(&events).Error
	return events, err
}

// GetMatchLineup returns the lineup for a match, sorted by number.
func (r *MatchRepo) GetMatchLineup(ctx context.Context, matchID uint) ([]models.MatchLineup, error) {
	var lineup []models.MatchLineup
	err := r.db.WithContext(ctx).
		Where("match_id = ?", matchID).
		Order("number ASC").
		Find(&lineup).Error
	return lineup, err
}

// FindByExternal returns a match by (provider, external_id).
func (r *MatchRepo) FindByExternal(ctx context.Context, provider, externalID string) (*models.Match, error) {
	var m models.Match
	err := r.db.WithContext(ctx).
		Where("provider = ? AND external_id = ?", provider, externalID).
		First(&m).Error
	if err != nil {
		return nil, err
	}
	return &m, nil
}

// UpsertMatchEvents inserts or updates events for a match. Events are keyed
// by (match_id, minute, event_type, player_name) to avoid duplicates on
// repeated polling.
func (r *MatchRepo) UpsertMatchEvents(ctx context.Context, matchID uint, events []external.MatchEvent) error {
	for _, e := range events {
		teamID, _ := strconv.Atoi(e.TeamID)
		playerID, _ := strconv.Atoi(e.PlayerID)
		ev := models.MatchEvent{
			MatchID:    matchID,
			Minute:     e.Minute,
			EventType:  e.EventType,
			TeamID:     uint(teamID),
			PlayerID:   uint(playerID),
			PlayerName: e.PlayerName,
			Detail:     e.Detail,
			Comment:    e.Comment,
		}
		err := r.db.WithContext(ctx).
			Clauses(clause.OnConflict{
				Columns: []clause.Column{
					{Name: "match_id"},
					{Name: "minute"},
					{Name: "event_type"},
					{Name: "player_name"},
				},
				DoUpdates: clause.AssignmentColumns([]string{"detail", "comment"}),
			}).
			Create(&ev).Error
		if err != nil {
			return fmt.Errorf("upsert event for match %d: %w", matchID, err)
		}
	}
	return nil
}

// UpsertMatchStats inserts or updates per-team match statistics, keyed by
// (match_id, team_id, stat_type).
func (r *MatchRepo) UpsertMatchStats(ctx context.Context, matchID uint, stats []external.MatchStat) error {
	for _, s := range stats {
		teamID, _ := strconv.Atoi(s.TeamID)
		ms := models.MatchStat{
			MatchID:  matchID,
			TeamID:   uint(teamID),
			StatType: s.StatType,
			Value:    s.Value,
		}
		err := r.db.WithContext(ctx).
			Clauses(clause.OnConflict{
				Columns: []clause.Column{
					{Name: "match_id"},
					{Name: "team_id"},
					{Name: "stat_type"},
				},
				DoUpdates: clause.AssignmentColumns([]string{"value"}),
			}).
			Create(&ms).Error
		if err != nil {
			return fmt.Errorf("upsert stat for match %d: %w", matchID, err)
		}
	}
	return nil
}

// GetMatchStats returns the statistics for a match, grouped by team.
func (r *MatchRepo) GetMatchStats(ctx context.Context, matchID uint) ([]models.MatchStat, error) {
	var stats []models.MatchStat
	err := r.db.WithContext(ctx).
		Where("match_id = ?", matchID).
		Order("team_id ASC, id ASC").
		Find(&stats).Error
	return stats, err
}
