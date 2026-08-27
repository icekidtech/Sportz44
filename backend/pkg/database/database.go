package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/icekidtech/Sportz44/backend/internal/models"
)

// Connect opens a GORM connection to PostgreSQL.
func Connect(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, err
	}
	return db, nil
}

// AutoMigrate creates/updates all tables from the GORM models (source of truth).
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.Club{},
		&models.Player{},
		&models.PlayerSeasonStats{},
		&models.Match{},
		&models.MatchEvent{},
		&models.MatchLineup{},
		&models.User{},
		&models.UserClubSubscription{},
		&models.NotificationPreference{},
	)
}
