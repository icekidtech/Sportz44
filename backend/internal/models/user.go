package models

import "time"

// User represents an application user.
type User struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	Username     string     `gorm:"size:50;uniqueIndex" json:"username"`
	Email        string     `gorm:"size:150;uniqueIndex" json:"email"`
	PasswordHash string     `gorm:"size:255" json:"-"`
	Role         string     `gorm:"size:20;default:user" json:"role"` // user | admin | moderator | analyst
	IsPremium    bool       `gorm:"default:false" json:"is_premium"`
	AvatarURL    string     `gorm:"size:500" json:"avatar_url"`
	FanLevel     int        `gorm:"default:1" json:"fan_level"`
	LastLoginAt  *time.Time `json:"last_login_at"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}
