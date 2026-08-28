package repository

import (
	"errors"

	"gorm.io/gorm"

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
