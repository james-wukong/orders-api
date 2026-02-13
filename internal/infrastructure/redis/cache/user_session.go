package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	userSession "github.com/james-wukong/orders-api/internal/domain/user_session"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

type UserSessionCache struct {
	redis *redis.Client
	log   *zerolog.Logger
}

func NewUserSessionCache(redis *redis.Client, log *zerolog.Logger) *UserSessionCache {
	return &UserSessionCache{
		redis: redis,
		log:   log,
	}
}

// GetByUserID returns a list of tokens that belong to a user id
func (c *UserSessionCache) GetByUserID(ctx context.Context, userID uuid.UUID) ([]userSession.UserSessions, error) {
	var sessions []userSession.UserSessions
	var session userSession.UserSessions
	// session token set key format: "user:session:userid"
	key := fmt.Sprintf("%s%s%s", userSession.RedisUserPrefix, userSession.RedisSessionPrefix, userID)
	// 1. Get all tokens belong to a user by using UserID from token set
	// SMembers return error nil and empty list if the key doesn't exist
	allMembers, err := c.redis.SMembers(ctx, key).Result()
	if err != nil {
		c.log.Error().Err(err).Msg("Error getting token set members from cache")
		return nil, err
	} else if len(allMembers) == 0 {
		return nil, nil // Cache miss, return nil without error
	}
	// 2. Get all tokens from Hash
	for _, v := range allMembers {
		// session token hash key format: "user:session:token:token"
		tokenKey := fmt.Sprintf("%s%s%s%s", userSession.RedisUserPrefix,
			userSession.RedisSessionPrefix,
			userSession.RedisTokenPrefix,
			v)
		err = c.redis.HGetAll(ctx, tokenKey).Scan(&session)
		if err != nil {
			c.log.Error().Err(err).Msgf("Error getting token: %s from HGetAll cache", v)
			return nil, err
		} else if session.ID == uuid.Nil {
			return nil, nil // Cache miss, return nil without error
		}
		sessions = append(sessions, session)
	}
	return sessions, nil
}
func (c *UserSessionCache) GetByToken(ctx context.Context, token string) (*userSession.UserSessions, error) {
	var session userSession.UserSessions
	tokenKey := fmt.Sprintf("%s%s%s%s", userSession.RedisUserPrefix,
		userSession.RedisSessionPrefix,
		userSession.RedisTokenPrefix,
		token)
	err := c.redis.HGetAll(ctx, tokenKey).Scan(&session)
	if err != nil {
		c.log.Error().Err(err).Msgf("Error getting token: %s from HGetAll cache", token)
		return nil, err
	} else if session.ID == uuid.Nil {
		return nil, nil // Cache miss, return nil without error
	}

	return &session, nil
}

// Set add session information by using token user:session:token:token as key
// and add its token to a set with key user:session:user_id
func (c *UserSessionCache) Set(ctx context.Context, entity *userSession.UserSessions) error {
	var jsonMap map[string]interface{}
	// token session hash key format: "user:session:token:token"
	tokenKey := fmt.Sprintf("%s%s%s%s", userSession.RedisUserPrefix,
		userSession.RedisSessionPrefix,
		userSession.RedisTokenPrefix,
		entity.Token)

	// struct to json string for HSet
	jsonStr, _ := json.Marshal(entity)
	// Unmarshal back from json string to map[string]interface{}
	json.Unmarshal(jsonStr, &jsonMap)
	// 1. Create user session in cache
	c.redis.HSet(ctx, tokenKey, jsonMap)
	err := c.redis.HSet(ctx, tokenKey, jsonMap).Err()
	if err != nil {
		c.log.Error().Err(err).Msg("Error setting user session in cache")
		return err
	}
	c.redis.Expire(ctx, tokenKey, 2*24*time.Hour)

	// 2. Create ID to tokens mapping set in redis
	// token set key format: "user:session:userid"
	userIDKey := fmt.Sprintf("%s%s%s", userSession.RedisUserPrefix,
		userSession.RedisSessionPrefix,
		entity.UserID)
	err = c.redis.SAdd(ctx, userIDKey, entity.Token).Err()
	if err != nil {
		c.log.Error().Err(err).Msg("Error add a token to user session in cache")
	}

	return nil
}
func (c *UserSessionCache) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	userIDKey := fmt.Sprintf("%s%s%s", userSession.RedisUserPrefix,
		userSession.RedisSessionPrefix,
		userID)
	// Retrieve tokens by user id
	tokens, err := c.redis.SMembers(ctx, userIDKey).Result()
	if err != nil {
		c.log.Error().Err(err).Msg("Error getting token set members from cache")
		return err
	} else if len(tokens) == 0 {
		return nil // Cache miss, return nil without error
	}
	pipe := c.redis.Pipeline()
	// 1. Delete Token set
	pipe.Del(ctx, userIDKey)
	// 2. Delete Session Token Hashes
	for _, v := range tokens {
		tokenKey := fmt.Sprintf("%s%s%s%s", userSession.RedisUserPrefix,
			userSession.RedisSessionPrefix,
			userSession.RedisTokenPrefix,
			v)
		pipe.Del(ctx, tokenKey)
	}
	if _, err := pipe.Exec(ctx); err != nil {
		// Check if the error is just a cancellation
		if err == context.Canceled {
			c.log.Warn().Msg("Session deletion aborted: context canceled")
			return err
		}
		c.log.Error().Err(err).Msg("Error updating user in cache")
		return err
	}

	return nil
}
func (c *UserSessionCache) DeleteByToken(ctx context.Context, token string) error {
	tokenKey := fmt.Sprintf("%s%s%s%s", userSession.RedisUserPrefix,
		userSession.RedisSessionPrefix,
		userSession.RedisTokenPrefix,
		token)
	// 1. get user id that owns this token
	userID, err := c.redis.HGet(ctx, tokenKey, "UserID").Result()
	if err != nil {
		if err == redis.Nil {
			return nil
		}
		c.log.Error().Err(err).Msg("Error retrieving user id from session hash in cache")
		return err
	}
	userIDKey := fmt.Sprintf("%s%s%s", userSession.RedisUserPrefix,
		userSession.RedisSessionPrefix,
		userID)
	pipe := c.redis.Pipeline()
	// 2. Delete token in session hash
	pipe.Del(ctx, tokenKey)
	// 3. Delete token in user token set
	pipe.SRem(ctx, userIDKey, token)
	if _, err := pipe.Exec(ctx); err != nil {
		// Check if the error is just a cancellation
		if err == context.Canceled {
			c.log.Warn().Msg("Session deletion aborted: context canceled")
			return err
		}
		c.log.Error().Err(err).Msg("Error updating user in cache")
		return err
	}

	return nil
}

// Update user session info
func (c *UserSessionCache) Update(ctx context.Context, entity *userSession.UserSessions) error {
	tokenKey := fmt.Sprintf("%s%s%s%s", userSession.RedisUserPrefix,
		userSession.RedisSessionPrefix,
		userSession.RedisTokenPrefix,
		entity.Token)
	session, err := c.redis.HGetAll(ctx, tokenKey).Result()
	if err != nil {
		if err == redis.Nil {
			return nil
		}
		c.log.Error().Err(err).Msg("Error retrieving user id from session hash in cache")
		return err
	}

	// convert entity struct to map
	var entityMap map[string]interface{}
	// 1. Convert struct to JSON bytes
	bytes, err := json.Marshal(entity)
	if err != nil {
		c.log.Error().Err(err).Msg("Error converting struct to JSON bytes in cache")
		return err
	}
	// 2. Convert JSON bytes to Map
	json.Unmarshal(bytes, &entityMap)
	for k, v := range session {
		if session[k] != entityMap[k] {
			// update the field in cache
			c.redis.HSet(ctx, tokenKey, k, v)
		}
	}

	return nil
}
