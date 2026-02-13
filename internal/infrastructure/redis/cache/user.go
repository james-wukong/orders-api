// Package cache implements the user repository using Redis as the underlying storage mechanism.
package cache

import (
	"context"
	"encoding/json"
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
	var u user.Users
	// Use HGetAll to get all fields of the user hash.
	// If the key doesn't exist, it will return nil error.
	key := fmt.Sprintf("%s%s", user.RedisUserPrefix, userID)
	// HGetAll(key) is special when the key does not exist,
	// it returns an empty map and nil error. We need to check if the user ID is actually present in the cache.
	err := c.redis.HGetAll(ctx, key).Scan(&u)
	if err != nil {
		c.log.Error().Err(err).Msg("Error retrieving user from Redis cache")
		return nil, err
	} else if u.ID == uuid.Nil {
		return nil, nil
	}

	return &u, nil
}

func (c *userCache) GetByEmail(ctx context.Context, email string) (*user.Users, error) {
	// Implement logic to retrieve user by email from Redis cache
	// We can use a separate key to map email to user ID, then fetch the user by ID.
	emailKey := fmt.Sprintf("%s%s%s", user.RedisUserPrefix, user.RedisEmailToIDPrefix, email)
	idStr, err := c.redis.Get(ctx, emailKey).Result()
	if err == redis.Nil {
		return nil, nil // Cache miss, return nil without error
	} else if err != nil {
		c.log.Error().Err(err).Msg("Error retrieving user ID from email mapping in Redis cache")
		return nil, err
	}

	userID, err := uuid.Parse(idStr)
	if err != nil {
		c.log.Error().Err(err).Msg("Error parsing user ID from Redis cache")
		return nil, err
	}
	return c.GetByID(ctx, userID)
}

func (c *userCache) Set(ctx context.Context, entity *user.Users) error {
	// Implement logic to set user data in Redis cache
	var jsonMap map[string]interface{}
	key := fmt.Sprintf("%s%s", user.RedisUserPrefix, entity.ID)

	// struct to json string for HSet
	jsonStr, _ := json.Marshal(entity)
	// Unmarshal back from json string to map[string]interface{}
	json.Unmarshal(jsonStr, &jsonMap)
	err := c.redis.HSet(ctx, key, jsonMap).Err()
	if err != nil {
		c.log.Error().Err(err).Msg("Error setting user in Redis cache")
		return err
	}
	// Set TTL for the user data
	c.redis.Expire(ctx, key, 7*24*time.Hour)

	return nil
}

func (c *userCache) Delete(ctx context.Context, userID uuid.UUID) error {
	// Implement logic to delete user data from Redis cache
	key := fmt.Sprintf("%s%s", user.RedisUserPrefix, userID)
	return c.redis.Del(ctx, key).Err()
}

func (c *userCache) Update(ctx context.Context, entity *user.Users) error {
	var jsonData map[string]interface{}
	idKey := fmt.Sprintf("%s%s", user.RedisUserPrefix, entity.ID)
	emailKey := fmt.Sprintf("%s%s%s", user.RedisUserPrefix, user.RedisEmailToIDPrefix, entity.Email)

	pipe := c.redis.Pipeline()
	// 1. Save the actual data
	jsonStr, _ := json.Marshal(entity)
	json.Unmarshal(jsonStr, &jsonData)
	pipe.HSet(ctx, idKey, jsonData)
	pipe.Expire(ctx, idKey, 7*24*time.Hour)

	// 2. Save the email to id mapping
	pipe.Set(ctx, emailKey, entity.ID.String(), 7*24*time.Hour)

	// 3. Execute the pipeline
	if _, err := pipe.Exec(ctx); err != nil {
		c.log.Error().Err(err).Msg("Error updating user in cache")
		return err
	}

	return nil
}

func (c *userCache) MapEmailToID(ctx context.Context, email string, id uuid.UUID) error {
	emailKey := fmt.Sprintf("%s%s%s", user.RedisUserPrefix, user.RedisEmailToIDPrefix, email)
	err := c.redis.Set(ctx, emailKey, id.String(), 7*24*time.Hour).Err()
	if err != nil {
		c.log.Error().Err(err).Msg("Error mapping email to ID in cache")
		return err
	}
	return nil
}

func (c *userCache) DeleteEmailToID(ctx context.Context, email string) error {
	emailKey := fmt.Sprintf("%s%s%s", user.RedisUserPrefix, user.RedisEmailToIDPrefix, email)
	err := c.redis.Del(ctx, emailKey).Err()
	if err != nil {
		c.log.Error().Err(err).Msg("Error deleting email to ID mapping from cache")
		return err
	}
	return nil
}
