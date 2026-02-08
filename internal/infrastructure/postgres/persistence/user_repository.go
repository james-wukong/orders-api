// Package postgres implements the user repository using PostgreSQL as the database.
// It implements the UserRepository interface defined in the user domain, providing methods for creating, retrieving, updating, and deleting user records in a PostgreSQL database using GORM as the ORM.
package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/james-wukong/orders-api/internal/domain/user"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{db: db}
}

// Implement the UserRepository interface methods here, using GORM to interact with the PostgreSQL database.
func (u *userRepository) Create(ctx context.Context, user *user.UserEntity) error {
	return u.db.WithContext(ctx).Create(user).Error
}

// Implement GetByID method
func (u *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*user.UserEntity, error) {
	var user user.UserEntity
	err := u.db.WithContext(ctx).First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Implement GetByEmail method
func (u *userRepository) GetByEmail(ctx context.Context, email string) (*user.UserEntity, error) {
	var user user.UserEntity
	err := u.db.WithContext(ctx).First(&user, "email = ?", email).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Implement Update method
func (u *userRepository) Update(ctx context.Context, user *user.UserEntity) error {
	return u.db.WithContext(ctx).Save(user).Error
}

// Implement Delete method
func (u *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return u.db.WithContext(ctx).Delete(&user.UserEntity{}, "id = ?", id).Error
}

// Implement List method
func (u *userRepository) List(ctx context.Context, filter *user.UserFilterEntity) ([]*user.UserEntity, error) {
	var users []*user.UserEntity
	query := u.db.WithContext(ctx).Model(&user.UserEntity{})
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
