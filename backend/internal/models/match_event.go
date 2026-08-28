package models

import "time"

// MatchEvent represents a single in-match event (goal, card, sub, etc).
type MatchEvent struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	MatchID          uint      `gorm:"index:idx_match_events_match,priority:1;uniqueIndex:idx_match_events_unique,priority:1" json:"match_id"`
	Minute           int       `gorm:"index:idx_match_events_match,priority:2;uniqueIndex:idx_match_events_unique,priority:2" json:"minute"`
	EventType        string    `gorm:"size:30;index;uniqueIndex:idx_match_events_unique,priority:3" json:"event_type"` // goal | card | substitution | injury | own_goal | penalty
	TeamID           uint      `json:"team_id"`
	PlayerID         uint      `json:"player_id"`
	PlayerName       string    `gorm:"size:100;uniqueIndex:idx_match_events_unique,priority:4" json:"player_name"`
	AssistPlayerID   uint      `json:"assist_player_id"`
	AssistPlayerName string    `gorm:"size:100" json:"assist_player_name"`
	Detail           string    `gorm:"size:100" json:"detail"` // e.g. "Yellow Card", "Normal Goal"
	Comment          string    `gorm:"size:500" json:"comment"`
	CreatedAt        time.Time `json:"created_at"`
}
