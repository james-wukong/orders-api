// Package persistence implements the restaurant repository using GORM for PostgreSQL
// It implements the restaurant.Repository interface defined in the restaurant domain, providing methods for creating, retrieving, updating, and deleting restaurant records in a PostgreSQL database using GORM as the ORM.
package persistence

import (
	"context"
	"errors"

	"github.com/james-wukong/orders-api/internal/domain/restaurant"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type restaurantRepository struct {
	db    *gorm.DB
	cache restaurant.Cache
	log   *zerolog.Logger
}

// NewRestaurantRepository creates a new instance of the GORM repository
func NewRestaurantRepository(db *gorm.DB, cache restaurant.Cache, log *zerolog.Logger) restaurant.Repository {
	return &restaurantRepository{
		db:    db,
		cache: cache,
		log:   log,
	}
}

// Create creates a new restaurant record in the database
// it also updates the cache with the new restaurant data and slug to id mapping
func (r *restaurantRepository) Create(ctx context.Context, res *restaurant.Restaurants) error {
	// Create the restaurant record in the database
	err := r.db.WithContext(ctx).Create(res).Error
	if err != nil {
		r.log.Error().Err(err).Msg("Error creating restaurant in database")
		return err
	}
	// Update cache with new restaurant data and slug to id mapping
	newCtx := context.WithoutCancel(ctx)
	go r.cache.Update(newCtx, res)

	return nil
}

func (r *restaurantRepository) GetByID(ctx context.Context, id uuid.UUID) (*restaurant.Restaurants, error) {
	var res *restaurant.Restaurants
	// 1. Try the cache first (if implemented)
	res, err := r.cache.GetByID(ctx, id)
	if err == nil {
		return res, nil // Cache hit
	}
	if err != redis.Nil {
		// Log the error but continue to fetch from DB
		r.log.Error().Err(err).Msg("Error fetching restaurant from cache")
	}
	// 2. If not found in cache, query the database
	err = r.db.WithContext(ctx).First(res, "id = ?", id).Error
	if err != nil {
		r.log.Error().Err(err).Msg("Error fetching restaurant from database")
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil if not found
		}
		return nil, err
	}
	// 3. found in DB, update cache
	if res != nil {
		newCtx := context.WithoutCancel(ctx)
		go r.cache.Update(newCtx, res)
	}

	return res, nil
}

func (r *restaurantRepository) GetBySlug(ctx context.Context, slug string) (*restaurant.Restaurants, error) {
	// 1. Try the cache first (if implemented)
	res, err := r.cache.GetBySlug(ctx, slug)
	if err == nil {
		return res, nil // Cache hit
	}
	if err != redis.Nil {
		// Log the error but continue to fetch from DB
		r.log.Error().Err(err).Msg("Error fetching slug id mapping from cache")
	}
	// 2. If not found in cache, query the database
	err = r.db.WithContext(ctx).Where("slug = ?", slug).First(&res).Error
	if err != nil {
		r.log.Error().Err(err).Msgf("Error fetching restaurant from database: %v", err)
		// If not found, return nil without error to indicate restaurant doesn't exist
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	// 3. found in DB, update cache
	r.log.Info().Msgf("Found restaurant: %v in database, updating cache", res)
	if res != nil {
		newCtx := context.WithoutCancel(ctx)
		go r.cache.Update(newCtx, res)
	}

	return res, nil
}

func (r *restaurantRepository) List(ctx context.Context, limit, offset int) ([]*restaurant.Restaurants, error) {
	var restaurants []*restaurant.Restaurants
	err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&restaurants).Error
	return restaurants, err
}

func (r *restaurantRepository) Update(ctx context.Context, res *restaurant.Restaurants) error {
	// Updates current record, only updating non-zero fields
	// If you want to update all fields (including zeros), use .Save(res)
	err := r.db.WithContext(ctx).Save(res).Error
	if err != nil {
		r.log.Error().Err(err).Msg("Error updating restaurant in database")
		return err
	}
	// Update cache with new restaurant data and slug to id mapping
	newCtx := context.WithoutCancel(ctx)
	go r.cache.Update(newCtx, res)

	return nil
}

func (r *restaurantRepository) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.db.WithContext(ctx).Delete(&restaurant.Restaurants{}, "id = ?", id).Error
	if err != nil {
		r.log.Error().Err(err).Msg("Error deleting restaurant from database")
		return err
	}
	// Delete from cache as well

	newCtx := context.WithoutCancel(ctx)
	go r.cache.Delete(newCtx, id)

	return nil
}
