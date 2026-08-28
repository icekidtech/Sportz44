package models

import "time"

// Injury tracks a player's injury status. This is a near-term follow-up
// within squad analytics; the model exists so the API surface is stable, but
// ingestion is not wired yet (endpoints return an empty list until then).
type Injury struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	PlayerID    uint      `gorm:"index" json:"player_id"`
	Type        string    `gorm:"size:100" json:"type"` // e.g. "Ankle Sprain"
	Severity    string    `gorm:"size:50" json:"severity"` // minor | moderate | severe
	Status      string    `gorm:"size:20;index" json:"status"` // injured | recovering | fit
	StartDate   *time.Time `json:"start_date"`
	EndDate     *time.Time `json:"end_date"`
	Description string    `gorm:"size:500" json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	Player Player `gorm:"foreignKey:PlayerID" json:"player,omitempty"`
}