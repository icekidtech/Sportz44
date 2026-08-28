package models

import "time"

// MatchStat holds a per-team statistic for a match (possession, shots, xG,
// corners, fouls, cards, etc). Values are stored as strings because providers
// return them in mixed formats (e.g. "58%", "12", "1.34").
type MatchStat struct {
	ID       uint   `gorm:"primaryKey" json:"id"`
	MatchID  uint   `gorm:"index:idx_match_stats_match,priority:1" json:"match_id"`
	TeamID   uint   `gorm:"index:idx_match_stats_match,priority:2" json:"team_id"`
	StatType string `gorm:"size:50;index:idx_match_stats_match,priority:3" json:"stat_type"` // possession | shots | shots_on_target | corners | fouls | yellow_cards | red_cards | xg | ...
	Value    string `gorm:"size:50" json:"value"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}