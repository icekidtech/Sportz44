package models

import "time"

// Competition represents a league or cup (e.g. La Liga, Copa del Rey,
// Champions League) aggregated from API-Sports.
type Competition struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Provider    string    `gorm:"size:50;uniqueIndex:idx_provider_external,priority:1" json:"provider"`   // "apisports", "footballdata"
	ExternalID  string    `gorm:"size:50;uniqueIndex:idx_provider_external,priority:2" json:"external_id"` // provider-specific ID
	APISportsID int       `gorm:"uniqueIndex" json:"api_sports_id"`                                        // legacy compat
	Name        string    `gorm:"size:150;not null" json:"name"`
	Type        string    `gorm:"size:20" json:"type"` // League | Cup
	Country     string    `gorm:"size:100" json:"country"`
	LogoURL     string    `gorm:"size:500" json:"logo_url"`
	Season      string    `gorm:"size:20" json:"season"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
