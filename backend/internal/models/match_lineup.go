package models

import "time"

// MatchLineup represents a player in a team's lineup for a specific match.
type MatchLineup struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	MatchID    uint      `gorm:"index" json:"match_id"`
	TeamID     uint      `json:"team_id"`
	PlayerID   uint      `json:"player_id"`
	PlayerName string    `gorm:"size:100" json:"player_name"`
	Position   string    `gorm:"size:20" json:"position"`
	Number     int       `json:"number"`
	IsStarter  bool      `gorm:"default:true" json:"is_starter"`
	CreatedAt  time.Time `json:"created_at"`
}
