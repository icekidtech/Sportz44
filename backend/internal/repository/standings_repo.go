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
	CleanSheets  int    `json:"clean_sheets"`
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

		// Clean sheets: a team keeps one when it concedes zero goals.
		if m.AwayScore == 0 {
			home.CleanSheets++
		}
		if m.HomeScore == 0 {
			away.CleanSheets++
		}

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
// goal events, sorted by goals, then assists. Assists are read from the
// assist fields embedded in each goal event, so assist-only players appear.
func (r *StandingsRepo) GetTopScorers(ctx context.Context, competitionID uint, limit int) ([]TopScorer, error) {
	type goalEvent struct {
		PlayerID         uint
		PlayerName       string
		ClubID           uint
		ClubName         string
		AssistPlayerID   uint
		AssistPlayerName string
		MatchID          uint
	}
	var events []goalEvent
	err := r.db.WithContext(ctx).
		Model(&models.MatchEvent{}).
		Select("e.player_id, e.player_name, c.id AS club_id, c.name AS club_name, e.assist_player_id, e.assist_player_name, e.match_id").
		Table("match_events e").
		Joins("JOIN matches m ON m.id = e.match_id").
		Joins("LEFT JOIN clubs c ON c.external_id = CAST(e.team_id AS TEXT)").
		Where("m.competition_id = ? AND e.event_type = 'goal'", competitionID).
		Scan(&events).Error
	if err != nil {
		return nil, err
	}

	players := map[uint]*TopScorer{}
	appearances := map[uint]map[uint]bool{}
	for _, e := range events {
		// Scorer.
		scorer := players[e.PlayerID]
		if scorer == nil {
			scorer = &TopScorer{PlayerID: e.PlayerID, PlayerName: e.PlayerName, ClubID: e.ClubID, ClubName: e.ClubName}
			players[e.PlayerID] = scorer
		}
		scorer.Goals++
		if appearances[e.PlayerID] == nil {
			appearances[e.PlayerID] = map[uint]bool{}
		}
		appearances[e.PlayerID][e.MatchID] = true

		// Assister (if any).
		if e.AssistPlayerID != 0 {
			assister := players[e.AssistPlayerID]
			if assister == nil {
				assister = &TopScorer{PlayerID: e.AssistPlayerID, PlayerName: e.AssistPlayerName, ClubID: e.ClubID, ClubName: e.ClubName}
				players[e.AssistPlayerID] = assister
			}
			assister.Assists++
			if appearances[e.AssistPlayerID] == nil {
				appearances[e.AssistPlayerID] = map[uint]bool{}
			}
			appearances[e.AssistPlayerID][e.MatchID] = true
		}
	}

	out := make([]TopScorer, 0, len(players))
	for _, p := range players {
		p.Appearances = len(appearances[p.PlayerID])
		out = append(out, *p)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Goals != out[j].Goals {
			return out[i].Goals > out[j].Goals
		}
		if out[i].Assists != out[j].Assists {
			return out[i].Assists > out[j].Assists
		}
		return out[i].PlayerName < out[j].PlayerName
	})
	if limit > 0 && len(out) > limit {
		out = out[:limit]
	}
	return out, nil
}
