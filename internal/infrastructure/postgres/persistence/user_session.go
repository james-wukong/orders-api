package persistence

import (
	"context"

	"github.com/google/uuid"
	session "github.com/james-wukong/orders-api/internal/domain/user_session"
	"gorm.io/gorm"
)

type UserSessionRepository struct {
	db *gorm.DB
}

func NewUserSessionRepository(db *gorm.DB) *UserSessionRepository {
	return &UserSessionRepository{db: db}
}

// Create creates a new user session in the database
func (r *UserSessionRepository) Create(
	ctx context.Context,
	entity *session.UserSessions,
) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

// GetByToken retrieves a user session by its token
func (r *UserSessionRepository) GetByToken(
	ctx context.Context,
	token string,
) (*session.UserSessions, error) {
	var userSession session.UserSessions
	err := r.db.WithContext(ctx).First(&userSession, "token = ?", token).Error
	if err != nil {
		return nil, err
	}
	return &userSession, nil
}

// GetByUserID retrieves all user sessions for a given user ID
func (r *UserSessionRepository) GetByUserID(
	ctx context.Context,
	id string,
) ([]session.UserSessions, error) {
	var userSessions []session.UserSessions
	err := r.db.WithContext(ctx).Where("user_id = ?", id).Find(&userSessions).Error
	if err != nil {
		return nil, err
	}
	return userSessions, nil
}

// DeleteByToken deletes a user session by its token
func (r *UserSessionRepository) DeleteByToken(
	ctx context.Context,
	token string,
) error {
	return r.db.WithContext(ctx).Where("token = ?", token).Delete(&session.UserSessions{}).Error
}

// DeleteByUserID deletes all user sessions for a given user ID
func (r *UserSessionRepository) DeleteByUserID(
	ctx context.Context,
	id uuid.UUID,
) error {
	return r.db.WithContext(ctx).Where("user_id = ?", id).Delete(&session.UserSessions{}).Error
}
