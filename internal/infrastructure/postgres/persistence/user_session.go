package persistence

import (
	"context"
	"errors"

	"github.com/google/uuid"
	session "github.com/james-wukong/orders-api/internal/domain/user_session"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type UserSessionRepository struct {
	db  *gorm.DB
	log *zerolog.Logger
}

func NewUserSessionRepository(
	db *gorm.DB,
	log *zerolog.Logger) *UserSessionRepository {
	return &UserSessionRepository{db: db, log: log}
}

// Create creates a new user session in the database
func (r *UserSessionRepository) Create(
	ctx context.Context,
	entity *session.UserSessions,
) error {
	// 1. Create a record in database
	err := r.db.WithContext(ctx).Create(entity).Error
	if err != nil {
		r.log.Error().Err(err).Msg("Error creating user session in database")
		return err
	}
	return nil
}

// Update user session in database
func (r *UserSessionRepository) Update(
	ctx context.Context,
	entity *session.UserSessions,
) error {
	// 1. update in database
	err := r.db.WithContext(ctx).Save(entity).Error
	if err != nil {
		r.log.Error().Err(err).Msg("Error update user session in database")
		return err
	}

	return nil
}

// GetByToken retrieves a user session by its token from database
func (r *UserSessionRepository) GetByToken(
	ctx context.Context,
	token string,
) (*session.UserSessions, error) {
	var session *session.UserSessions

	// First, Take, Last expects a single record, return ErrRecordNotFound when no record exists
	// Find expects a collection, returns error nil and empty collection when no record exists
	err := r.db.WithContext(ctx).First(&session, "token = ?", token).Error
	if err != nil {
		r.log.Error().Err(err).Msg("Error fetching user session from database")
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil if not found
		}
		return nil, err
	}

	return session, nil
}

// GetByUserID retrieves all user sessions for a given user ID
func (r *UserSessionRepository) GetByUserID(
	ctx context.Context,
	userID uuid.UUID,
) ([]session.UserSessions, error) {
	var userSessions []session.UserSessions

	// search database with user id
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(userSessions).Error
	if err != nil || len(userSessions) == 0 {
		r.log.Error().Err(err).Msg("Error searching user session in database")
		return nil, err
	}
	return userSessions, nil
}

// DeleteByToken deletes a user session by its token
// Delete from database first, then redis
func (r *UserSessionRepository) DeleteByToken(
	ctx context.Context,
	token string,
) error {
	// 1. delete from database
	err := r.db.WithContext(ctx).Where("token = ?", token).Delete(&session.UserSessions{}).Error
	if err != nil {
		r.log.Error().Err(err).Msg("Error deleting user session by token in database")
		return err
	}

	return nil
}

// DeleteByUserID deletes all user sessions for a given user ID
// Delete from database first, then redis
func (r *UserSessionRepository) DeleteByUserID(
	ctx context.Context,
	userID uuid.UUID,
) error {
	// 1. Delete from database first
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&session.UserSessions{}).Error
	if err != nil {
		r.log.Error().Err(err).Msg("Error deleting user session by user id in database")
		return err
	}

	return nil
}
