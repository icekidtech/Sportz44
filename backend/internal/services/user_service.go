package services

import (
	"context"

	"github.com/icekidtech/Sportz44/backend/internal/models"
	"github.com/icekidtech/Sportz44/backend/internal/repository"
)

// UserService orchestrates user subscriptions and notification preferences.
type UserService struct {
	users *repository.UserRepo
}

// NewUserService creates a new UserService.
func NewUserService(users *repository.UserRepo) *UserService {
	return &UserService{users: users}
}

// ListSubscriptions returns the clubs a user follows.
func (s *UserService) ListSubscriptions(ctx context.Context, userID uint) ([]models.UserClubSubscription, error) {
	return s.users.ListSubscriptions(ctx, userID)
}

// AddSubscription follows a club for a user.
func (s *UserService) AddSubscription(ctx context.Context, userID, clubID uint) error {
	return s.users.AddSubscription(ctx, userID, clubID)
}

// RemoveSubscription unfollows a club for a user.
func (s *UserService) RemoveSubscription(ctx context.Context, userID, clubID uint) error {
	return s.users.RemoveSubscription(ctx, userID, clubID)
}

// GetNotificationPrefs returns a user's notification preferences.
func (s *UserService) GetNotificationPrefs(ctx context.Context, userID uint) (*models.NotificationPreference, error) {
	return s.users.GetNotificationPrefs(ctx, userID)
}

// UpdateNotificationPrefs persists a user's notification preferences.
func (s *UserService) UpdateNotificationPrefs(ctx context.Context, prefs *models.NotificationPreference) error {
	return s.users.UpdateNotificationPrefs(ctx, prefs)
}
