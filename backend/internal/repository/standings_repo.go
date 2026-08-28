package repository

import (
	"context"
	"sort"

	"github.com/icekidtech/Sportz44/backend/internal/models"
	"gorm.io/gorm"
)

// StandingRow is a single row in a league table, computed from finished
// matches (no separate standings table needed).
type StandingRow struct {
	ClubID       uint   `json:"club_id"`
	ClubName     string `json:"club_name"`
	LogoURL      string `json:"logo_url"`
	Played       int    `json:"played"`
	Won          int    `json:"won"`
	Drawn        int    `json:"drawn"`
	Lost         int    `json:"lost"`
	GoalsFor     int    `json:"goals_for"`
	GoalsAgainst int    `json:"goals_against"`
	GoalDiff     int    `json:"goal_diff"`
	Points       int    `json:"points"`
}

// TopScorer is a row in the golden-boot standings.
type TopScorer struct {
	PlayerID    uint   `json:"player_id"`
	PlayerName  string `json:"player_name"`
	ClubID      uint   `json:"club_id"`
	ClubName    string `json:"club_name"`
	Goals       int    `json:"goals"`
	Assists     int    `json:"assists"`
	Appearances int    `json:"appearances"`
}

// StandingsRepo computes league tables and top scorers from match data.
type StandingsRepo struct {
	db *gorm.DB
}

// NewStandingsRepo creates a new StandingsRepo.
func NewStandingsRepo(db *gorm.DB) *StandingsRepo {
	return &StandingsRepo{db: db}
}

// GetStandings computes the league table for a competition from its finished
// matches, sorted by points, then goal difference, then goals for.
func (r *StandingsRepo) GetStandings(ctx context.Context, competitionID uint) ([]StandingRow, error) {
	var matches []models.Match
	err := r.db.WithContext(ctx).
		Preload("HomeClub").
		Preload("AwayClub").
		Where("competition_id = ? AND status = ?", competitionID, "finished").
		Find(&matches).Error
	if err != nil {
		return nil, err
	}

	rows := map[uint]*StandingRow{}
	for _, m := range matches {
		home := rows[m.HomeClubID]
		if home == nil {
			home = &StandingRow{ClubID: m.HomeClubID, ClubName: m.HomeClub.Name, LogoURL: m.HomeClub.LogoURL}
			rows[m.HomeClubID] = home
		}
		away := rows[m.AwayClubID]
		if away == nil {
			away = &StandingRow{ClubID: m.AwayClubID, ClubName: m.AwayClub.Name, LogoURL: m.AwayClub.LogoURL}
			rows[m.AwayClubID] = away
		}

		home.Played++
		away.Played++
		home.GoalsFor += m.HomeScore
		home.GoalsAgainst += m.AwayScore
		away.GoalsFor += m.AwayScore
		away.GoalsAgainst += m.HomeScore

		switch {
		case m.HomeScore > m.AwayScore:
			home.Won++
			away.Lost++
			home.Points += 3
		case m.HomeScore < m.AwayScore:
			away.Won++
			home.Lost++
			away.Points += 3
		default:
			home.Drawn++
			away.Drawn++
			home.Points++
			away.Points++
		}
	}

	out := make([]StandingRow, 0, len(rows))
	for _, row := range rows {
		row.GoalDiff = row.GoalsFor - row.GoalsAgainst
		out = append(out, *row)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Points != out[j].Points {
			return out[i].Points > out[j].Points
		}
		if out[i].GoalDiff != out[j].GoalDiff {
			return out[i].GoalDiff > out[j].GoalDiff
		}
		return out[i].GoalsFor > out[j].GoalsFor
	})
	return out, nil
}

// GetTopScorers computes the golden-boot standings for a competition from
// goal events, sorted by goals, then assists.
func (r *StandingsRepo) GetTopScorers(ctx context.Context, competitionID uint, limit int) ([]TopScorer, error) {
	type agg struct {
		PlayerID    uint
		PlayerName  string
		ClubID      uint
		Goals       int
		Assists     int
		Appearances int
	}
	var rows []agg
	err := r.db.WithContext(ctx).
		Model(&models.MatchEvent{}).
		Select(`e.player_id, e.player_name, e.team_id AS club_id,
		        SUM(CASE WHEN e.event_type = 'goal' THEN 1 ELSE 0 END) AS goals,
		        SUM(CASE WHEN e.event_type = 'assist' THEN 1 ELSE 0 END) AS assists,
		        COUNT(DISTINCT e.match_id) AS appearances`).
		Table("match_events e").
		Joins("JOIN matches m ON m.id = e.match_id").
		Where("m.competition_id = ? AND e.event_type IN ('goal','assist')", competitionID).
		Group("e.player_id, e.player_name, e.team_id").
		Order("goals DESC, assists DESC").
		Limit(limit).
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	out := make([]TopScorer, 0, len(rows))
	for _, rw := range rows {
		out = append(out, TopScorer{
			PlayerID:    rw.PlayerID,
			PlayerName:  rw.PlayerName,
			ClubID:      rw.ClubID,
			Goals:       rw.Goals,
			Assists:     rw.Assists,
			Appearances: rw.Appearances,
		})
	}
	return out, nil
}
