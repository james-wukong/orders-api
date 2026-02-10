// Package cache implements the user repository using Redis as the underlying storage mechanism.
package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/james-wukong/orders-api/internal/domain/user"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

type userCache struct {
	redis *redis.Client
	log   *zerolog.Logger
}

func NewUserCache(redis *redis.Client, log *zerolog.Logger) *userCache {
	return &userCache{redis: redis, log: log}
}

func (c *userCache) GetByID(ctx context.Context, userID uuid.UUID) (*user.Users, error) {
	// Implement logic to retrieve user by ID from Redis cache
	// 1. Fetch all fields from the Hash
	key := fmt.Sprintf("%s%s", user.RedisUserPrefix, userID)
	result, err := c.redis.HGetAll(ctx, key).Result()
	if err != nil {
		c.log.Error().Err(err).Msg("Error retrieving user from Redis cache")
		return nil, err
	}

	// 2. Check if the hash actually exists (HGetAll returns empty map if key is missing)
	if len(result) == 0 {
		return nil, user.ErrUserNotFound
	}

	// 3. Map the result back to the struct
	// go-redis can scan a map directly into a struct using 'redis' tags
	var u user.Users
	err = c.redis.HGetAll(ctx, key).Scan(&u)
	if err != nil {
		return nil, fmt.Errorf("failed to scan redis hash: %w", err)
	}

	return &u, nil
}

func (c *userCache) GetByEmail(ctx context.Context, email string) (*user.Users, error) {
	// Implement logic to retrieve user by email from Redis cache
	return nil, nil
}

func (c *userCache) Set(ctx context.Context, user *user.Users) error {
	// Implement logic to set user data in Redis cache
	key := fmt.Sprintf("user_hash:%s", user.ID)

	// We use a map to represent the Hash fields
	fields := map[string]interface{}{
		"email":      user.Email,
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"role":       user.Role,
		"is_active":  user.IsActive,
		"last_login": user.LastLogin.Format(time.RFC3339), // Store as ISO string
		"created_at": user.CreatedAt.Format(time.RFC3339),
		"updated_at": user.UpdatedAt.Format(time.RFC3339),
	}

	return c.redis.HSet(ctx, key, fields).Err()
}

func (c *userCache) Delete(ctx context.Context, userID uuid.UUID) error {
	// Implement logic to delete user data from Redis cache
	key := fmt.Sprintf("user_hash:%s", userID)
	return c.redis.Del(ctx, key).Err()
}

func (c *userCache) Update(ctx context.Context, user *user.Users) error {
	// Implement logic to update user data in Redis cache

	return nil
}
