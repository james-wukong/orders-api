// Package repository is cache decorator, it implements read-through caching
package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/james-wukong/orders-api/internal/domain/restaurant"
	"github.com/rs/zerolog"
)

type CachedRestaurantRepository struct {
	repo      restaurant.Repository
	redisRepo restaurant.RedisCache
	log       *zerolog.Logger
}

func NewCachedRestaurantRepository(repo restaurant.Repository,
	redisRepo restaurant.RedisCache,
	log *zerolog.Logger) *CachedRestaurantRepository {
	return &CachedRestaurantRepository{
		repo:      repo,
		redisRepo: redisRepo,
		log:       log,
	}
}

// Create creates a new restaurant record in the database
// it also create redis string for slug to id mapping
func (c *CachedRestaurantRepository) Create(ctx context.Context, entity *restaurant.Restaurants) error {
	// Create the restaurant record in the database
	err := c.repo.Create(ctx, entity)
	if err != nil {
		return err
	}

	// create or insert slug to id mapping to the set
	newCtx := context.WithoutCancel(ctx)
	go c.redisRepo.Update(newCtx, entity)

	return nil
}

// GetByID retrieves restaurant information from repository, it tries the cache first
// If not found in cache, it goes to database. It caches the result if find a catch in database
func (c *CachedRestaurantRepository) GetByID(ctx context.Context, id uuid.UUID) (*restaurant.Restaurants, error) {
	var res *restaurant.Restaurants
	// 1. Try the cache first (if implemented)
	res, err := c.redisRepo.GetByID(ctx, id)
	if err == nil && res.ID != uuid.Nil {
		return res, nil // Cache hit
	}

	// 2. If not found in cache, query the database
	res, err = c.repo.GetByID(ctx, id)
	// 2.1 if not found in database
	if err != nil || res == nil {
		c.log.Error().Err(err).Msg("Error fetching restaurant from database")
		return nil, err
	}
	// 2.2 found in DB, update cache
	if res != nil {
		newCtx := context.WithoutCancel(ctx)
		go c.redisRepo.Update(newCtx, res)
	}

	return res, nil
}

// GetBySlug retrieves restaurant information from repository, it tries the cache first
// If not found in cache, it goes to database. It caches the result if find a catch in database
func (c *CachedRestaurantRepository) GetBySlug(
	ctx context.Context,
	slug string,
) (*restaurant.Restaurants, error) {
	// 1. Try the cache first (if implemented)
	res, err := c.redisRepo.GetBySlug(ctx, slug)
	if err == nil && res != nil {
		return res, nil // Cache hit
	}

	// 2. If not found in cache, query the database
	res, err = c.repo.GetBySlug(ctx, slug)
	// 2.1 If not found in database
	if err != nil || res == nil {
		c.log.Error().Err(err).Msg("Error fetching restaurant from database")
		return nil, err
	}
	// 3. found in DB, update cache
	if res != nil {
		newCtx := context.WithoutCancel(ctx)
		go c.redisRepo.Update(newCtx, res)
	}

	return res, nil
}

// List return a list of restaurant based on filters
func (c *CachedRestaurantRepository) List(ctx context.Context, limit, offset int) ([]restaurant.Restaurants, error) {
	// var restaurants []*restaurant.Restaurants

	return c.repo.List(ctx, limit, offset)
}

// Update saves restaurant info in both database and redis
func (c *CachedRestaurantRepository) Update(ctx context.Context, entity *restaurant.Restaurants) error {
	// Updates current record in database
	err := c.repo.Update(ctx, entity)
	if err != nil {
		c.log.Error().Err(err).Msg("Error updating restaurant in database")
		return err
	}
	// Update cache with new restaurant data and slug to id mapping
	newCtx := context.WithoutCancel(ctx)
	go c.redisRepo.Update(newCtx, entity)

	return nil
}

// Delete will remove a restaurant info from both database and redis by using user id as key
func (c *CachedRestaurantRepository) Delete(ctx context.Context, id uuid.UUID) error {
	// 1. Removes it from database
	err := c.repo.Delete(ctx, id)
	if err != nil {
		c.log.Error().Err(err).Msg("Error deleting restaurant from database")
		return err
	}
	// 2. Delete from cache as well
	newCtx := context.WithoutCancel(ctx)
	go c.redisRepo.Delete(newCtx, id)

	return nil
}
