// Package restaurant defines the domain model and repository interface for managing restaurants.
// It represents the core business logic and data structures related to restaurants.
package restaurant

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, restaurant *Restaurants) error
	GetByID(ctx context.Context, id uuid.UUID) (*Restaurants, error)
	GetBySlug(ctx context.Context, slug string) (*Restaurants, error)
	List(ctx context.Context, limit, offset int) ([]*Restaurants, error)
	Update(ctx context.Context, restaurant *Restaurants) error
	Delete(ctx context.Context, id uuid.UUID) error
}
