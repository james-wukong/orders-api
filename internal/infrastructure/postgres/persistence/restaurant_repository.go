// Package postgres implements the restaurant repository using GORM for PostgreSQL
// It implements the restaurant.Repository interface defined in the restaurant domain, providing methods for creating, retrieving, updating, and deleting restaurant records in a PostgreSQL database using GORM as the ORM.
package postgres

import (
	"context"
	"errors"

	"github.com/james-wukong/orders-api/internal/domain/restaurant"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type restaurantRepository struct {
	db *gorm.DB
}

// NewRestaurantRepository creates a new instance of the GORM repository
func NewRestaurantRepository(db *gorm.DB) restaurant.Repository {
	return &restaurantRepository{db: db}
}

func (r *restaurantRepository) Create(ctx context.Context, res *restaurant.Restaurants) error {
	return r.db.WithContext(ctx).Create(res).Error
}

func (r *restaurantRepository) GetByID(ctx context.Context, id uuid.UUID) (*restaurant.Restaurants, error) {
	var res restaurant.Restaurants
	err := r.db.WithContext(ctx).First(&res, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil if not found
		}
		return nil, err
	}
	return &res, nil
}

func (r *restaurantRepository) GetBySlug(ctx context.Context, slug string) (*restaurant.Restaurants, error) {
	var res restaurant.Restaurants
	err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&res).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &res, nil
}

func (r *restaurantRepository) List(ctx context.Context, limit, offset int) ([]*restaurant.Restaurants, error) {
	var restaurants []*restaurant.Restaurants
	err := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&restaurants).Error
	return restaurants, err
}

func (r *restaurantRepository) Update(ctx context.Context, res *restaurant.Restaurants) error {
	// Updates current record, only updating non-zero fields
	// If you want to update all fields (including zeros), use .Save(res)
	return r.db.WithContext(ctx).Save(res).Error
}

func (r *restaurantRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&restaurant.Restaurants{}, "id = ?", id).Error
}
