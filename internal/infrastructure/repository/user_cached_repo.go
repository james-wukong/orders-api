package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/james-wukong/orders-api/internal/domain/user"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type CachedUserRepository struct {
	repo      user.Repository
	redisRepo user.RedisCache
	log       *zerolog.Logger
}

func NewCachedUserRepository(repo user.Repository,
	redisRepo user.RedisCache,
	log *zerolog.Logger) *CachedUserRepository {
	return &CachedUserRepository{
		repo:      repo,
		redisRepo: redisRepo,
		log:       log,
	}
}

// Create a user in postgres database
// create user hash in redis
func (c *CachedUserRepository) Create(ctx context.Context, entity *user.Users) error {
	// First, create the user in the database
	err := c.repo.Create(ctx, entity)
	if err != nil {
		c.log.Error().Err(err).Msg("Error creating user in database")
		return err
	}
	// If successful, also set the user in cache
	newCtx := context.WithoutCancel(ctx)
	go c.redisRepo.Set(newCtx, entity)

	return nil
}

// GetByID tries to retrieve User info from redis, then database
func (c *CachedUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*user.Users, error) {
	// 1. Try the cache first (if implemented)
	record, err := c.redisRepo.GetByID(ctx, id)
	if err == nil && record != nil {
		return record, nil // Cache hit
	}
	// 2. Log the error but continue to fetch from DB
	c.log.Error().Err(err).Msg("Error fetching user from cache")
	record, err = c.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil if not found
		}
		return nil, err
	}

	// 3. Update cache with the fetched user data
	newCtx := context.WithoutCancel(ctx)
	go c.redisRepo.Set(newCtx, record)

	return record, nil
}

// GetByEmail tries to retrieve User info from redis, then database
func (c *CachedUserRepository) GetByEmail(ctx context.Context, email string) (*user.Users, error) {
	// 1. Try the cache first (if implemented)
	record, err := c.redisRepo.GetByEmail(ctx, email)
	if err == nil || record != nil {
		return record, nil // Cache hit
	}

	// 2. If not found in cache, query the database
	record, err = c.repo.GetByEmail(ctx, email)
	// 2.1 if not found in database
	if err != nil || record == nil {
		c.log.Error().Err(err).Msg("Error fetching restaurant from database")
		return nil, err
	}
	// 2.2 found in DB, update cache
	newCtx := context.WithoutCancel(ctx)
	go c.redisRepo.Set(newCtx, record)

	return record, nil
}

// Update updates user info in both database and redis
func (c *CachedUserRepository) Update(ctx context.Context, entity *user.Users) error {
	// Update the user record in the database
	err := c.repo.Update(ctx, entity)

	if err != nil {
		c.log.Error().Err(err).Msg("Error updating user in database")
		return err
	}
	// If successful, also update the user in cache
	newCtx := context.WithoutCancel(ctx)
	go c.redisRepo.Update(newCtx, entity)

	return nil
}

// Delete remove user info from both database and redis
// including email->userID mapping in redis
func (c *CachedUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	var email string
	// 1. get the user by ID to retrieve the email for cache deletion
	record, err := c.redisRepo.GetByID(ctx, id)
	if err != nil || record == nil {
		c.log.Error().Err(err).Msgf("Error find user with id: %s", id)
		return err
	}
	// 2.Delete the user record from the database
	err = c.repo.Delete(ctx, id)
	if err != nil {
		c.log.Error().Err(err).Msg("Error deleting user from database")
		return err
	}
	// 3. If successful, also delete the user from cache
	newCtx := context.WithoutCancel(ctx)
	go c.redisRepo.Delete(newCtx, id)

	// 4. Also delete the email to ID mapping from cache
	go c.redisRepo.DeleteEmailToID(newCtx, email)

	return nil
}

// List lists by filters
func (c *CachedUserRepository) List(ctx context.Context, filter *user.UserFilterEntity) ([]user.Users, error) {
	users, err := c.repo.List(ctx, filter)
	if err != nil {
		c.log.Error().Err(err).Msg("Error deleting user from database")
	}
	return users, err
}
