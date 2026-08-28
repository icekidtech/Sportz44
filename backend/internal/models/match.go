package models

import "time"

// Match represents a football fixture (scheduled, live, or finished).
type Match struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	APISportsID   int       `gorm:"uniqueIndex" json:"api_sports_id"`
	CompetitionID uint      `gorm:"index" json:"competition_id"`
	Season        string    `gorm:"size:20" json:"season"`
	HomeClubID    uint      `gorm:"index:idx_matches_club,priority:1" json:"home_club_id"`
	AwayClubID    uint      `gorm:"index:idx_matches_club,priority:1" json:"away_club_id"`
	MatchDate     time.Time `gorm:"index:idx_matches_club,priority:2" json:"match_date"`
	Status        string    `gorm:"size:20;index" json:"status"` // scheduled | live | finished | postponed
	HomeScore     int       `json:"home_score"`
	AwayScore     int       `json:"away_score"`
	Venue         string    `gorm:"size:200" json:"venue"`
	Referee       string    `gorm:"size:100" json:"referee"`
	Minute        int       `json:"minute"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`

	Competition Competition `gorm:"foreignKey:CompetitionID" json:"competition,omitempty"`
	HomeClub    Club        `gorm:"foreignKey:HomeClubID" json:"home_club,omitempty"`
	AwayClub    Club        `gorm:"foreignKey:AwayClubID" json:"away_club,omitempty"`
}
