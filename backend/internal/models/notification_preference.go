package models

import "time"

// NotificationPreference stores a user's notification toggles.
type NotificationPreference struct {
	ID                uint      `gorm:"primaryKey" json:"id"`
	UserID            uint      `gorm:"uniqueIndex" json:"user_id"`
	KickoffReminder   bool      `gorm:"default:true" json:"kickoff_reminder"`
	GoalAlerts        bool      `gorm:"default:true" json:"goal_alerts"`
	RedCardAlerts     bool      `gorm:"default:true" json:"red_card_alerts"`
	FullTimeResult    bool      `gorm:"default:true" json:"full_time_result"`
	NewsAlerts        bool      `gorm:"default:true" json:"news_alerts"`
	WhatsAppEnabled   bool      `gorm:"default:false" json:"whatsapp_enabled"`
	PushEnabled       bool      `gorm:"default:true" json:"push_enabled"`
	DoNotDisturbStart string    `gorm:"size:5" json:"dnd_start"`
	DoNotDisturbEnd   string    `gorm:"size:5" json:"dnd_end"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}
