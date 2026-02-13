// Package persistence implements the restaurant repository using GORM for PostgreSQL
// It implements the restaurant.Repository interface defined in the restaurant domain, providing methods for creating, retrieving, updating, and deleting restaurant records in a PostgreSQL database using GORM as the ORM.
package persistence

import (
	"context"
	"errors"

	"github.com/james-wukong/orders-api/internal/domain/restaurant"
	"github.com/rs/zerolog"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type restaurantRepository struct {
	db    *gorm.DB
	cache restaurant.RedisCache
	log   *zerolog.Logger
}

// NewRestaurantRepository creates a new instance of the GORM repository
func NewRestaurantRepository(db *gorm.DB,
	log *zerolog.Logger,
) *restaurantRepository {
	return &restaurantRepository{
		db:  db,
		log: log,
	}
}

// Create creates a new restaurant record in the database
func (r *restaurantRepository) Create(ctx context.Context, res *restaurant.Restaurants) error {
	// Create the restaurant record in the database
	err := r.db.WithContext(ctx).Create(res).Error
	if err != nil {
		r.log.Error().Err(err).Msg("Error creating restaurant in database")
		return err
	}

	return nil
}

// GetByID retrieved restaurant info using user id as key from database
func (r *restaurantRepository) GetByID(ctx context.Context, id uuid.UUID) (*restaurant.Restaurants, error) {
	var res *restaurant.Restaurants
	// First, Take, Last expects a single record, return ErrRecordNotFound when no record exists
	// Find expects a collection, returns error nil and empty collection when no record exists
	err := r.db.WithContext(ctx).First(res, "id = ?", id).Error
	if err != nil {
		r.log.Error().Err(err).Msg("Error fetching restaurant from database")
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil if not found
		}
		return nil, err
	}

	return res, nil
}

// GetBySlug retrieves restuarant info using slug as key from database
func (r *restaurantRepository) GetBySlug(ctx context.Context, slug string) (*restaurant.Restaurants, error) {
	var res *restaurant.Restaurants
	// Query the database
	err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&res).Error
	if err != nil {
		r.log.Error().Err(err).Msgf("Error fetching restaurant from database: %v", err)
		// If not found, return nil without error to indicate restaurant doesn't exist
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return res, nil
}

// List gets a list of restaurant based on filters
// TODO: Add filter parameter in this method
func (r *restaurantRepository) List(ctx context.Context, limit, offset int) ([]restaurant.Restaurants, error) {
	var restaurants []restaurant.Restaurants
	err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&restaurants).Error
	return restaurants, err
}

// Save updates restaurant info
func (r *restaurantRepository) Update(ctx context.Context, res *restaurant.Restaurants) error {
	// Updates current record, only updating non-zero fields
	// If you want to update all fields (including zeros), use .Save(res)
	err := r.db.WithContext(ctx).Save(res).Error
	if err != nil {
		r.log.Error().Err(err).Msg("Error updating restaurant in database")
		return err
	}

	return nil
}

// Delete removes a restaurant record from database
func (r *restaurantRepository) Delete(ctx context.Context, id uuid.UUID) error {
	err := r.db.WithContext(ctx).Delete(&restaurant.Restaurants{}, "id = ?", id).Error
	if err != nil {
		r.log.Error().Err(err).Msg("Error deleting restaurant from database")
		return err
	}

	return nil
}
