// Package postgres implements the user repository using PostgreSQL as the database.
// It implements the UserRepository interface defined in the user domain, providing methods for creating, retrieving, updating, and deleting user records in a PostgreSQL database using GORM as the ORM.
package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/james-wukong/orders-api/internal/domain/user"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) user.Repository {
	return &userRepository{db: db}
}

// Implement the UserRepository interface methods here, using GORM to interact with the PostgreSQL database.
func (u *userRepository) Create(ctx context.Context, user *user.Users) error {
	return u.db.WithContext(ctx).Create(user).Error
}

// Implement GetByID method
func (u *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*user.Users, error) {
	var user user.Users
	err := u.db.WithContext(ctx).First(&user, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil if not found
		}
		return nil, err
	}
	return &user, nil
}

// Implement GetByEmail method
func (u *userRepository) GetByEmail(ctx context.Context, email string) (*user.Users, error) {
	var record user.Users
	err := u.db.WithContext(ctx).First(&record, "email = ?", email).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Return your own domain error
			return nil, user.ErrUserNotFound
		}
		return nil, err
	}
	return &record, nil
}

// Implement Update method
func (u *userRepository) Update(ctx context.Context, user *user.Users) error {
	return u.db.WithContext(ctx).Save(user).Error
}

// Implement Delete method
func (u *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return u.db.WithContext(ctx).Delete(&user.Users{}, "id = ?", id).Error
}

// Implement List method
func (u *userRepository) List(ctx context.Context, filter *user.UserFilterEntity) ([]*user.Users, error) {
	var users []*user.Users
	query := u.db.WithContext(ctx).Model(&user.Users{})
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
