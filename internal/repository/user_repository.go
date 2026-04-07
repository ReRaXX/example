package repository

import (
	"user-api/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserRepository defines the interface for user database operations
type UserRepository interface {
	Create(user *entity.User) error
	FindByID(id uuid.UUID) (*entity.User, error)
	FindByIDWithPassword(id uuid.UUID) (*entity.User, error)
	FindByEmail(email string) (*entity.User, error)
	FindByEmailWithPassword(email string) (*entity.User, error)
	FindAll() ([]*entity.User, error)
	Update(id uuid.UUID, updates map[string]interface{}) (*entity.User, error)
	Delete(id uuid.UUID) error
}

// userRepository implements UserRepository
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new instance of UserRepository
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// Create creates a new user in the database
func (r *userRepository) Create(user *entity.User) error {
	return r.db.Create(user).Error
}

// FindByID finds a user by ID (without password)
func (r *userRepository) FindByID(id uuid.UUID) (*entity.User, error) {
	var user entity.User
	err := r.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByIDWithPassword finds a user by ID (with password)
func (r *userRepository) FindByIDWithPassword(id uuid.UUID) (*entity.User, error) {
	var user entity.User
	err := r.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByEmail finds a user by email (without password)
func (r *userRepository) FindByEmail(email string) (*entity.User, error) {
	var user entity.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByEmailWithPassword finds a user by email (with password)
func (r *userRepository) FindByEmailWithPassword(email string) (*entity.User, error) {
	var user entity.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindAll finds all users ordered by creation date (newest first)
func (r *userRepository) FindAll() ([]*entity.User, error) {
	var users []*entity.User
	err := r.db.Order("created_at DESC").Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

// Update updates a user with the given updates
func (r *userRepository) Update(id uuid.UUID, updates map[string]interface{}) (*entity.User, error) {
	var user entity.User
	err := r.db.Model(&user).Where("id = ?", id).Updates(updates).Error
	if err != nil {
		return nil, err
	}
	
	// Return the updated user
	return r.FindByID(id)
}

// Delete soft deletes a user
func (r *userRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&entity.User{}, id).Error
}
