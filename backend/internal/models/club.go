package models

import "time"

// Club represents a football club aggregated from API-Sports.
type Club struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	Provider      string    `gorm:"size:50;uniqueIndex:idx_club_provider_external,priority:1" json:"provider"`
	ExternalID    string    `gorm:"size:50;uniqueIndex:idx_club_provider_external,priority:2" json:"external_id"`
	APISportsID   int       `gorm:"uniqueIndex" json:"api_sports_id"`
	Name          string    `gorm:"size:100;not null" json:"name"`
	ShortName     string    `gorm:"size:50" json:"short_name"`
	Country       string    `gorm:"size:100" json:"country"`
	CompetitionID uint      `gorm:"index" json:"competition_id"` // primary competition
	LogoURL       string    `gorm:"size:500" json:"logo_url"`
	Stadium       string    `gorm:"size:200" json:"stadium"`
	Colors        string    `gorm:"size:100" json:"colors"`
	Founded       int       `json:"founded"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`

	Competition Competition `gorm:"foreignKey:CompetitionID" json:"competition,omitempty"`
}
