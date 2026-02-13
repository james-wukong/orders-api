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
	// restaurant info hash key format: "restaurant:id"
	key := fmt.Sprintf("%s%s", restaurant.RedisRestaurantPrefix, id)
	// HGetALl return nil if the key doesn't exist.
	err := r.redis.HGetAll(ctx, key).Scan(&res)
	if err != nil {
		r.log.Error().Err(err).Msg("Error scanning restaurant from cache")
		return nil, err
	} else if res.ID == uuid.Nil {
		return nil, nil // Cache miss, return nil without error
	}
	return &res, nil
}

// GetBySlug retrieves the restaurant ID from the slug mapping and then fetches the restaurant data
// if this fails, it will return to database to fetch the restaurant data and update the cache
func (r *restaurantCache) GetBySlug(ctx context.Context, slug string) (*restaurant.Restaurants, error) {
	var idStr string
	// Slug to id map key format: "restaurant:slug_to_id:slug"
	slugKey := fmt.Sprintf("%s%s%s", restaurant.RedisRestaurantPrefix, restaurant.RedisSlupToIDPrefix, slug)
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

// Set create a restaurant hash in redis
func (r *restaurantCache) Set(ctx context.Context, entity *restaurant.Restaurants) error {
	var jsonMap map[string]interface{}
	// restaurant info hash key format: "restaurant:id"
	key := fmt.Sprintf("%s%s", restaurant.RedisRestaurantPrefix, entity.ID)

	// struct to json string for HSet
	jsonStr, _ := json.Marshal(entity)
	// Unmarshal back from json string to map[string]interface{}
	json.Unmarshal(jsonStr, &jsonMap)
	r.redis.HSet(ctx, key, jsonMap)
	err := r.redis.HSet(ctx, key, jsonMap).Err()
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
	// restaurant info hash key format: "restaurant:id"
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
	// restaurant info hash key format: "restaurant:id"
	idKey := fmt.Sprintf("%s%s", restaurant.RedisRestaurantPrefix, entity.ID)
	// Slug to id map key format: "restaurant:slug_to_id:slug"
	slugKey := fmt.Sprintf("%s%s%s", restaurant.RedisRestaurantPrefix, restaurant.RedisSlupToIDPrefix, entity.Slug)

	pipe := r.redis.Pipeline()
	// 1. Save the actual data
	inrec, _ := json.Marshal(entity)
	json.Unmarshal(inrec, &jsonData)
	pipe.HSet(ctx, idKey, jsonData)
	pipe.Expire(ctx, idKey, 7*24*time.Hour) // Set TTL for the restaurant data

	// 2. Save the slug to id mapping
	pipe.Set(ctx, slugKey, entity.ID.String(), 7*24*time.Hour) // Set TTL for the slug mapping

	// 3. Execute the pipeline
	if _, err := pipe.Exec(ctx); err != nil {
		r.log.Error().Err(err).Msg("Error updating restaurant in cache")
		return err
	}

	return nil
}

func (r *restaurantCache) MapSlugToID(ctx context.Context, id uuid.UUID, slug string) error {
	// Slug to id map key format: "restaurant:slug_to_id:slug"
	key := fmt.Sprintf("%s%s%s", restaurant.RedisRestaurantPrefix, restaurant.RedisSlupToIDPrefix, slug)
	err := r.redis.Set(ctx, key, id.String(), 7*24*time.Hour).Err()
	if err != nil {
		r.log.Error().Err(err).Msg("Error mapping slug to ID in cache")
		return err
	}
	return nil
}

func (r *restaurantCache) DeleteSlugToID(ctx context.Context, slug string) error {
	// Slug to id map key format: "restaurant:slug_to_id:slug"
	key := fmt.Sprintf("%s%s%s", restaurant.RedisRestaurantPrefix, restaurant.RedisSlupToIDPrefix, slug)
	err := r.redis.Del(ctx, key).Err()
	if err != nil {
		r.log.Error().Err(err).Msg("Error deleting slug to ID mapping from cache")
		return err
	}
	return nil
}
