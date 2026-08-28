package services

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/icekidtech/Sportz44/backend/internal/models"
	"github.com/icekidtech/Sportz44/backend/internal/repository"
	"github.com/icekidtech/Sportz44/backend/pkg/jwt"
)

var (
	// ErrUserExists is returned when a username/email is already taken.
	ErrUserExists = errors.New("user already exists")
	// ErrInvalidCredentials is returned on a bad login.
	ErrInvalidCredentials = errors.New("invalid credentials")
)

// AuthService handles registration, login, and token issuance.
type AuthService struct {
	repo       *repository.UserRepo
	jwtSecret  string
	accessTTL  time.Duration
	refreshTTL time.Duration
}

// NewAuthService creates an AuthService.
func NewAuthService(repo *repository.UserRepo, jwtSecret string, accessTTL, refreshTTL time.Duration) *AuthService {
	return &AuthService{
		repo:       repo,
		jwtSecret:  jwtSecret,
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
	}
}

// RegisterInput is the data needed to create a new account.
type RegisterInput struct {
	Username string
	Email    string
	Password string
}

// Register creates a user and returns the user plus a fresh token pair.
func (s *AuthService) Register(in RegisterInput) (*models.User, string, string, error) {
	if _, err := s.repo.FindByEmail(in.Email); err == nil {
		return nil, "", "", ErrUserExists
	} else if !errors.Is(err, repository.ErrNotFound) {
		return nil, "", "", err
	}
	if _, err := s.repo.FindByUsername(in.Username); err == nil {
		return nil, "", "", ErrUserExists
	} else if !errors.Is(err, repository.ErrNotFound) {
		return nil, "", "", err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", "", err
	}

	user := &models.User{
		Username:     in.Username,
		Email:        in.Email,
		PasswordHash: string(hash),
		Role:         "user",
	}
	if err := s.repo.Create(user); err != nil {
		return nil, "", "", err
	}

	access, refresh, err := jwt.Generate(user.ID, user.Role, s.jwtSecret, s.accessTTL, s.refreshTTL)
	if err != nil {
		return nil, "", "", err
	}
	return user, access, refresh, nil
}

// Login verifies credentials (email or username) and returns a token pair.
func (s *AuthService) Login(identifier, password string) (*models.User, string, string, error) {
	user, err := s.repo.FindByEmail(identifier)
	if err != nil {
		user, err = s.repo.FindByUsername(identifier)
		if err != nil {
			return nil, "", "", ErrInvalidCredentials
		}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, "", "", ErrInvalidCredentials
	}

	now := time.Now()
	user.LastLoginAt = &now
	_ = s.repo.Update(user)

	access, refresh, err := jwt.Generate(user.ID, user.Role, s.jwtSecret, s.accessTTL, s.refreshTTL)
	if err != nil {
		return nil, "", "", err
	}
	return user, access, refresh, nil
}

// Refresh validates a refresh token and issues a new token pair.
func (s *AuthService) Refresh(refreshToken string) (*models.User, string, string, error) {
	claims, err := jwt.Parse(refreshToken, s.jwtSecret)
	if err != nil {
		return nil, "", "", ErrInvalidCredentials
	}
	user, err := s.repo.FindByID(claims.UserID)
	if err != nil {
		return nil, "", "", ErrInvalidCredentials
	}
	access, refresh, err := jwt.Generate(user.ID, user.Role, s.jwtSecret, s.accessTTL, s.refreshTTL)
	if err != nil {
		return nil, "", "", err
	}
	return user, access, refresh, nil
}

// GetUser returns a user by ID (used by the /me endpoint).
func (s *AuthService) GetUser(id uint) (*models.User, error) {
	return s.repo.FindByID(id)
}
