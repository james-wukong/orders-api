package repository

import (
	"context"

	"github.com/google/uuid"
	uSession "github.com/james-wukong/orders-api/internal/domain/user_session"
	"github.com/rs/zerolog"
)

type CachedUserSessionRepository struct {
	repo      uSession.Repository
	redisRepo uSession.RedisCache
	log       *zerolog.Logger
}

func NewCachedUserSessionRepository(
	repo uSession.Repository,
	redisRepo uSession.RedisCache,
	log *zerolog.Logger,
) *CachedUserSessionRepository {
	return &CachedUserSessionRepository{
		repo:      repo,
		redisRepo: redisRepo,
		log:       log,
	}
}

// Create creates a new user session in the database and redis
func (c *CachedUserSessionRepository) Create(
	ctx context.Context,
	entity *uSession.UserSessions,
) error {
	// 1. Create a record in database
	err := c.repo.Create(ctx, entity)
	if err != nil {
		c.log.Error().Err(err).Msg("Error creating user session in database")
		return err
	}
	// 2. If successful, search redis and database for the previous active session
	//
	//
	// also set the user session in redis
	newCtx := context.WithoutCancel(ctx)
	go c.redisRepo.Set(newCtx, entity)

	return nil
}

// Update user session in database and redis
func (c *CachedUserSessionRepository) Update(
	ctx context.Context,
	entity *uSession.UserSessions,
) error {
	// 1. update in database
	err := c.repo.Update(ctx, entity)
	if err != nil {
		c.log.Error().Err(err).Msg("Error update user session in database")
		return err
	}

	// 2. If successful, also update the user session in redis
	newCtx := context.WithoutCancel(ctx)
	go c.redisRepo.Update(newCtx, entity)

	return nil
}

// GetByToken retrieves a user session by its token from database
func (c *CachedUserSessionRepository) GetByToken(
	ctx context.Context,
	token string,
) (*uSession.UserSessions, error) {
	var session *uSession.UserSessions

	// First, Take, Last expects a single record, return ErrRecordNotFound when no record exists
	// Find expects a collection, returns error nil and empty collection when no record exists
	session, err := c.redisRepo.GetByToken(ctx, token)
	if err == nil && session != nil {
		return session, nil // Cache hit
	}

	// If not found in redis, try to search the postgre database
	session, err = c.repo.GetByToken(ctx, token)
	if err != nil || session == nil {
		c.log.Error().Err(err).Msg("Error fetching user session info from database")
		return nil, err
	}

	return session, nil
}

// GetByUserID retrieves all user sessions for a given user ID
func (c *CachedUserSessionRepository) GetByUserID(
	ctx context.Context,
	userID uuid.UUID,
) ([]uSession.UserSessions, error) {
	// 1. try to find them in cache
	sessions, err := c.redisRepo.GetByUserID(ctx, userID)
	if err == nil && len(sessions) != 0 {
		// cache hit
		return sessions, nil
	}

	// 2. Not found in cache, try database
	sessions, err = c.repo.GetByUserID(ctx, userID)
	if err != nil || len(sessions) == 0 {
		c.log.Error().Err(err).Msg("Error searching user session in database")
		return nil, err
	}
	return sessions, nil
}

// DeleteByToken deletes a user session by its token
// Delete from database first, then redis
func (c *CachedUserSessionRepository) DeleteByToken(
	ctx context.Context,
	token string,
) error {
	// 1. delete from database
	err := c.repo.DeleteByToken(ctx, token)
	if err != nil {
		c.log.Error().Err(err).Msg("Error deleting user session by token in database")
		return err
	}

	// 2. delete from redis
	newCtx := context.WithoutCancel(ctx)
	go c.redisRepo.DeleteByToken(newCtx, token)

	return nil
}

// DeleteByUserID deletes all user sessions for a given user ID
// Delete from database first, then redis
func (c *CachedUserSessionRepository) DeleteByUserID(
	ctx context.Context,
	userID uuid.UUID,
) error {
	// 1. Delete from database first
	err := c.repo.DeleteByUserID(ctx, userID)
	if err != nil {
		c.log.Error().Err(err).Msg("Error deleting user session by user id in database")
		return err
	}

	// 2. delete from redis
	newCtx := context.WithoutCancel(ctx)
	go c.redisRepo.DeleteByUserID(newCtx, userID)

	return nil
}
