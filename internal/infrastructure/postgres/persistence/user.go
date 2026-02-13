// Package persistence implements the user repository using PostgreSQL as the database.
// It implements the UserRepository interface defined in the user domain, providing methods for creating, retrieving, updating, and deleting user records in a PostgreSQL database using GORM as the ORM.
package persistence

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/james-wukong/orders-api/internal/domain/user"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type userRepository struct {
	db  *gorm.DB
	log *zerolog.Logger
}

func NewUserRepository(db *gorm.DB, log *zerolog.Logger) user.Repository {
	return &userRepository{db: db, log: log}
}

// Implement the UserRepository interface methods here, using GORM to interact with the PostgreSQL database.
func (r *userRepository) Create(ctx context.Context, user *user.Users) error {
	// First, create the user in the database
	err := r.db.WithContext(ctx).Create(user).Error
	if err != nil {
		r.log.Error().Err(err).Msg("Error creating user in database")
		return err
	}
	return nil
}

// Implement GetByID method
func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*user.Users, error) {
	var record user.Users
	err := r.db.WithContext(ctx).First(&record, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil if not found
		}
		return nil, err
	}
	return &record, nil
}

// Implement GetByEmail method
func (r *userRepository) GetByEmail(ctx context.Context, email string) (*user.Users, error) {
	var record user.Users
	err := r.db.WithContext(ctx).First(&record, "email = ?", email).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Return domain error
			return nil, user.ErrUserNotFound
		}
		return nil, err
	}

	return &record, nil
}

// Implement Update method
func (r *userRepository) Update(ctx context.Context, user *user.Users) error {
	// Update the user record in the database
	err := r.db.WithContext(ctx).Save(user).Error
	if err != nil {
		r.log.Error().Err(err).Msg("Error updating user in database")
		return err
	}

	return nil
}

// Implement Delete method
func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	// Delete the user record from the database
	err := r.db.WithContext(ctx).Delete(&user.Users{}, "id = ?", id).Error
	if err != nil {
		r.log.Error().Err(err).Msg("Error deleting user from database")
		return err
	}

	return nil
}

// Implement List method
func (r *userRepository) List(ctx context.Context, filter *user.UserFilterEntity) ([]user.Users, error) {
	var users []user.Users
	query := r.db.WithContext(ctx).Model(&user.Users{})
	if filter != nil {
		if filter.Email != nil && *filter.Email != "" {
			query = query.Where("email = ?", filter.Email)
		}
		if len(filter.Roles) > 0 {
			query = query.Where("role IN ?", filter.Roles)
		}
		if filter.IsActive != nil {
			query = query.Where("is_active = ?", filter.IsActive)
		}
		if filter.FirstName != nil && *filter.FirstName != "" {
			query = query.Where("first_name ILIKE ?", "%"+*filter.FirstName+"%")
		}
		if filter.LastName != nil && *filter.LastName != "" {
			query = query.Where("last_name ILIKE ?", "%"+*filter.LastName+"%")
		}
	}

	err := query.Limit(filter.Limit).
		Offset((filter.Page - 1) * filter.Limit).
		Find(&users).
		Error
	return users, err
}
