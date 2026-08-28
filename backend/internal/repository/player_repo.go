package repository

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/icekidtech/Sportz44/backend/internal/external"
	"github.com/icekidtech/Sportz44/backend/internal/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// PlayerRepo handles DB operations for players and their season stats.
type PlayerRepo struct {
	db *gorm.DB
}

// NewPlayerRepo creates a new PlayerRepo.
func NewPlayerRepo(db *gorm.DB) *PlayerRepo {
	return &PlayerRepo{db: db}
}

// UpsertPlayers inserts or updates the given players for a club, keyed by
// (provider, external_id).
func (r *PlayerRepo) UpsertPlayers(ctx context.Context, players []external.Player, clubID uint) error {
	for _, p := range players {
		apiID, _ := strconv.Atoi(p.ProviderID)
		m := models.Player{
			Provider:     "apisports", // squads currently come from API-Sports
			ExternalID:   p.ProviderID,
			APISportsID:  apiID,
			ClubID:       clubID,
			Name:         p.Name,
			Position:     p.Position,
			JerseyNumber: p.JerseyNumber,
			Nationality:  p.Nationality,
			BirthDate:    p.BirthDate,
			PhotoURL:     p.PhotoURL,
			Rating:       p.Rating,
		}
		err := r.db.WithContext(ctx).
			Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "provider"}, {Name: "external_id"}},
				DoUpdates: clause.AssignmentColumns([]string{"club_id", "name", "position", "jersey_number", "nationality", "birth_date", "photo_url", "rating"}),
			}).
			Create(&m).Error
		if err != nil {
			return fmt.Errorf("upsert player %s/%s: %w", m.Provider, m.ExternalID, err)
		}
	}
	return nil
}

// ListByClub returns the squad roster for a club.
func (r *PlayerRepo) ListByClub(ctx context.Context, clubID uint) ([]models.Player, error) {
	var players []models.Player
	err := r.db.WithContext(ctx).
		Where("club_id = ?", clubID).
		Order("jersey_number ASC").
		Find(&players).Error
	return players, err
}

// FindByID returns a player by primary key, with their club preloaded.
func (r *PlayerRepo) FindByID(ctx context.Context, id uint) (*models.Player, error) {
	var p models.Player
	err := r.db.WithContext(ctx).Preload("Club").First(&p, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &p, nil
}

// GetSeasonStats returns a player's season statistics, newest season first.
func (r *PlayerRepo) GetSeasonStats(ctx context.Context, playerID uint) ([]models.PlayerSeasonStats, error) {
	var stats []models.PlayerSeasonStats
	err := r.db.WithContext(ctx).
		Preload("Competition").
		Where("player_id = ?", playerID).
		Order("season DESC").
		Find(&stats).Error
	return stats, err
}

// PlayerFormEntry is a single match in a player's recent form.
type PlayerFormEntry struct {
	MatchID     uint      `json:"match_id"`
	MatchDate   time.Time `json:"match_date"`
	Opponent    string    `json:"opponent"`
	HomeAway    string    `json:"home_away"` // home | away
	Result      string    `json:"result"`    // W | D | L
	Goals       int       `json:"goals"`
	Assists     int       `json:"assists"`
	YellowCards int       `json:"yellow_cards"`
	RedCards    int       `json:"red_cards"`
	Rating      float64   `json:"rating"`
}

// GetForm returns a player's last `limit` matches with their contributions,
// computed from match events (matched by provider player ID).
func (r *PlayerRepo) GetForm(ctx context.Context, playerID uint, limit int) ([]PlayerFormEntry, error) {
	p, err := r.FindByID(ctx, playerID)
	if err != nil {
		return nil, err
	}

	// Find the player's most recent matches where they have events.
	var matches []models.Match
	err = r.db.WithContext(ctx).
		Distinct("matches.*").
		Joins("JOIN match_events e ON e.match_id = matches.id").
		Where("e.player_id = ?", p.APISportsID).
		Order("matches.match_date DESC").
		Limit(limit).
		Find(&matches).Error
	if err != nil {
		return nil, err
	}

	form := make([]PlayerFormEntry, 0, len(matches))
	for _, m := range matches {
		// Player's events in this match.
		var events []models.MatchEvent
		if err := r.db.WithContext(ctx).
			Where("match_id = ? AND player_id = ?", m.ID, p.APISportsID).
			Find(&events).Error; err != nil {
			return nil, err
		}

		entry := PlayerFormEntry{
			MatchID:   m.ID,
			MatchDate: m.MatchDate,
		}

		// Determine home/away + opponent + result.
		if m.HomeClubID == p.ClubID {
			entry.HomeAway = "home"
			entry.Opponent = m.AwayClub.Name
			switch {
			case m.HomeScore > m.AwayScore:
				entry.Result = "W"
			case m.HomeScore < m.AwayScore:
				entry.Result = "L"
			default:
				entry.Result = "D"
			}
		} else {
			entry.HomeAway = "away"
			entry.Opponent = m.HomeClub.Name
			switch {
			case m.AwayScore > m.HomeScore:
				entry.Result = "W"
			case m.AwayScore < m.HomeScore:
				entry.Result = "L"
			default:
				entry.Result = "D"
			}
		}

		// Aggregate the player's contributions.
		for _, e := range events {
			switch e.EventType {
			case "goal":
				entry.Goals++
				// The player may also be credited as the assister on a goal
				// scored by a teammate in the same match.
				if e.AssistPlayerID == uint(p.APISportsID) {
					entry.Assists++
				}
			case "card":
				if e.Detail == "Red Card" {
					entry.RedCards++
				} else {
					entry.YellowCards++
				}
			}
		}
		form = append(form, entry)
	}
	return form, nil
}

// GetInjuries returns a player's injuries. Injury tracking is a near-term
// follow-up; the model does not exist yet, so this returns an empty list.
func (r *PlayerRepo) GetInjuries(ctx context.Context, playerID uint) ([]models.Injury, error) {
	return []models.Injury{}, nil
}
