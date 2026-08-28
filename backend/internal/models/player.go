package models

import "time"

// Player represents a footballer belonging to a club.
type Player struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	APISportsID  int        `gorm:"uniqueIndex" json:"api_sports_id"`
	ClubID       uint       `gorm:"index" json:"club_id"`
	Name         string     `gorm:"size:100;not null" json:"name"`
	Position     string     `gorm:"size:50" json:"position"`
	JerseyNumber int        `json:"jersey_number"`
	Nationality  string     `gorm:"size:100" json:"nationality"`
	BirthDate    *time.Time `json:"birth_date"`
	PhotoURL     string     `gorm:"size:500" json:"photo_url"`
	Rating       float64    `json:"rating"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`

	Club Club `gorm:"foreignKey:ClubID" json:"club,omitempty"`
}
