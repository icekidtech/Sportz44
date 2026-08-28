package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"github.com/icekidtech/Sportz44/backend/internal/models"
)

// ErrNotFound is returned when a record does not exist.
var ErrNotFound = errors.New("record not found")

// UserRepo provides database access for users.
type UserRepo struct {
	db *gorm.DB
}

// NewUserRepo creates a UserRepo.
func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db: db}
}

// Create inserts a new user.
func (r *UserRepo) Create(u *models.User) error {
	return r.db.Create(u).Error
}

// FindByEmail looks up a user by email.
func (r *UserRepo) FindByEmail(email string) (*models.User, error) {
	var u models.User
	if err := r.db.Where("email = ?", email).First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &u, nil
}

// FindByUsername looks up a user by username.
func (r *UserRepo) FindByUsername(username string) (*models.User, error) {
	var u models.User
	if err := r.db.Where("username = ?", username).First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &u, nil
}

// FindByID looks up a user by primary key.
func (r *UserRepo) FindByID(id uint) (*models.User, error) {
	var u models.User
	if err := r.db.First(&u, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &u, nil
}

// Update persists changes to a user.
func (r *UserRepo) Update(u *models.User) error {
	return r.db.Save(u).Error
}

// ListSubscriptions returns the clubs a user follows.
func (r *UserRepo) ListSubscriptions(ctx context.Context, userID uint) ([]models.UserClubSubscription, error) {
	var subs []models.UserClubSubscription
	err := r.db.WithContext(ctx).
		Preload("Club").
		Where("user_id = ?", userID).
		Find(&subs).Error
	return subs, err
}

// AddSubscription follows a club for a user. It is idempotent.
func (r *UserRepo) AddSubscription(ctx context.Context, userID, clubID uint) error {
	sub := models.UserClubSubscription{UserID: userID, ClubID: clubID}
	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(&sub).Error
}

// RemoveSubscription unfollows a club for a user.
func (r *UserRepo) RemoveSubscription(ctx context.Context, userID, clubID uint) error {
	return r.db.WithContext(ctx).
		Where("user_id = ? AND club_id = ?", userID, clubID).
		Delete(&models.UserClubSubscription{}).Error
}

// GetNotificationPrefs returns a user's notification preferences, creating
// defaults on first access.
func (r *UserRepo) GetNotificationPrefs(ctx context.Context, userID uint) (*models.NotificationPreference, error) {
	var prefs models.NotificationPreference
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&prefs).Error
	if err == nil {
		return &prefs, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	// Create defaults.
	prefs = models.NotificationPreference{UserID: userID}
	if err := r.db.WithContext(ctx).Create(&prefs).Error; err != nil {
		return nil, err
	}
	return &prefs, nil
}

// UpdateNotificationPrefs persists a user's notification preferences.
func (r *UserRepo) UpdateNotificationPrefs(ctx context.Context, prefs *models.NotificationPreference) error {
	return r.db.WithContext(ctx).Save(prefs).Error
}
