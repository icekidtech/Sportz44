package models

import "time"

// PlayerSeasonStats holds a player's statistics for a single season/competition.
type PlayerSeasonStats struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	PlayerID      uint      `gorm:"index:idx_player_stats_player,priority:1" json:"player_id"`
	Season        string    `gorm:"size:20" json:"season"`
	CompetitionID uint      `gorm:"index" json:"competition_id"`
	Goals         int       `json:"goals"`
	Assists       int       `json:"assists"`
	YellowCards   int       `json:"yellow_cards"`
	RedCards      int       `json:"red_cards"`
	MinutesPlayed int       `json:"minutes_played"`
	Appearances   int       `json:"appearances"`
	PassAccuracy  float64   `json:"pass_accuracy"`
	Tackles       int       `json:"tackles"`
	Interceptions int       `json:"interceptions"`
	Rating        float64   `json:"rating"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`

	Player      Player      `gorm:"foreignKey:PlayerID" json:"player,omitempty"`
	Competition Competition `gorm:"foreignKey:CompetitionID" json:"competition,omitempty"`
}
