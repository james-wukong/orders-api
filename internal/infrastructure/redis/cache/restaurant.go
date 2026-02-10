package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/james-wukong/orders-api/internal/domain/restaurant"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

type restaurantCache struct {
	redis *redis.Client
	log   *zerolog.Logger
}

func NewRestaurantCache(redis *redis.Client, log *zerolog.Logger) *restaurantCache {
	return &restaurantCache{redis: redis, log: log}
}

// GetByID returns nil if the restaurant is not found in cache, and returns error only if there is an actual error fetching from cache
func (r *restaurantCache) GetByID(ctx context.Context, id uuid.UUID) (*restaurant.Restaurants, error) {
	var res restaurant.Restaurants
	key := fmt.Sprintf("%s%s", restaurant.RedisRestaurantPrefix, id)
	// HGetALl return nil if the key doesn't exist.
	err := r.redis.HGetAll(ctx, key).Scan(&res)
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Cache miss, return nil without error
		}
		r.log.Error().Err(err).Msg("Error scanning restaurant from cache")
		return nil, err
	}
	if res.ID == uuid.Nil {
		return nil, nil // Cache miss, return nil without error
	}
	return &res, nil
}

// GetBySlug retrieves the restaurant ID from the slug mapping and then fetches the restaurant data
// if this fails, it will return to database to fetch the restaurant data and update the cache
func (r *restaurantCache) GetBySlug(ctx context.Context, slug string) (*restaurant.Restaurants, error) {
	var idStr string
	slugKey := fmt.Sprintf("%s%s", restaurant.RedisRestaurantPrefix, slug)
	// First, get the restaurant ID from the slug mapping
	idStr, err := r.redis.Get(ctx, slugKey).Result()
	if err == redis.Nil {
		return nil, nil // Cache miss, return nil without error
	} else if err != nil {
		r.log.Error().Err(err).Msg("Error scanning slug id mapping from cache")
		return nil, err
	}
	// If we got the ID, fetch the restaurant data
	id, err := uuid.Parse(idStr)
	if err != nil {
		r.log.Error().Err(err).Msg("Error parsing restaurant ID from cache")
		return nil, err
	}
	return r.GetByID(ctx, id)

}

func (r *restaurantCache) Set(ctx context.Context, entity *restaurant.Restaurants) error {
	var jsonData map[string]interface{}
	key := fmt.Sprintf("%s%s", restaurant.RedisRestaurantPrefix, entity.ID)

	inrec, _ := json.Marshal(entity)
	json.Unmarshal(inrec, &jsonData)
	r.redis.HSet(ctx, key, jsonData)
	err := r.redis.HSet(ctx, key, jsonData).Err()
	if err != nil {
		r.log.Error().Err(err).Msg("Error setting restaurant in cache")
		return err
	}
	r.redis.Expire(ctx, key, 7*24*time.Hour) // Set TTL for the restaurant data

	return nil
}

// Delete removes the restaurant data and the slug to id mapping from cache
func (r *restaurantCache) Delete(ctx context.Context, id uuid.UUID) error {
	var slug string
	key := fmt.Sprintf("%s%s", restaurant.RedisRestaurantPrefix, id)
	// Get the slug
	err := r.redis.HGet(ctx, key, "slug").Scan(&slug)
	if err == redis.Nil {
		return nil // Cache miss, nothing to delete, return nil without error
	}
	if err != nil {
		r.log.Error().Err(err).Msg("Error scanning restaurant ID from cache for deletion")
		return err
	}
	// Delete slug to id mapping
	r.DeleteSlugToID(ctx, slug)
	// Delete the restaurant data
	err = r.redis.Del(ctx, key).Err()
	if err != nil {
		r.log.Error().Err(err).Msg("Error deleting restaurant from cache")
		return err
	}
	return nil
}

// Update updates the restaurant data in cache and also updates the slug to id mapping
func (r *restaurantCache) Update(ctx context.Context, entity *restaurant.Restaurants) error {
	var jsonData map[string]interface{}
	idKey := fmt.Sprintf("%s%s", restaurant.RedisRestaurantPrefix, entity.ID)
	slugKey := fmt.Sprintf("%s%s", restaurant.RedisRestaurantPrefix, entity.Slug)

	pipe := r.redis.Pipeline()
	// Save the actual data
	inrec, _ := json.Marshal(entity)
	json.Unmarshal(inrec, &jsonData)
	pipe.HSet(ctx, idKey, jsonData)
	pipe.Expire(ctx, idKey, 7*24*time.Hour) // Set TTL for the restaurant data

	// Save the slug to id mapping
	pipe.Set(ctx, slugKey, entity.ID.String(), 7*24*time.Hour) // Set TTL for the slug mapping

	if _, err := pipe.Exec(ctx); err != nil {
		r.log.Error().Err(err).Msg("Error updating restaurant in cache")
		return err
	}

	return nil
}

func (r *restaurantCache) MapSlugToID(ctx context.Context, id uuid.UUID, slug string) error {
	key := fmt.Sprintf("%s%s", restaurant.RedisRestaurantPrefix, slug)
	err := r.redis.Set(ctx, key, id.String(), 7*24*time.Hour).Err()
	if err != nil {
		r.log.Error().Err(err).Msg("Error mapping slug to ID in cache")
		return err
	}
	return nil
}

func (r *restaurantCache) DeleteSlugToID(ctx context.Context, slug string) error {
	key := fmt.Sprintf("%s%s", restaurant.RedisRestaurantPrefix, slug)
	err := r.redis.Del(ctx, key).Err()
	if err != nil {
		r.log.Error().Err(err).Msg("Error deleting slug to ID mapping from cache")
		return err
	}
	return nil
}
