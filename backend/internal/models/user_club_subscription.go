package models

import "time"

// UserClubSubscription tracks which clubs a user follows.
type UserClubSubscription struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index" json:"user_id"`
	ClubID    uint      `gorm:"index" json:"club_id"`
	CreatedAt time.Time `json:"created_at"`

	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Club Club `gorm:"foreignKey:ClubID" json:"club,omitempty"`
}
