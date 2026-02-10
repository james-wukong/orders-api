package restaurant

import (
	"context"

	"github.com/google/uuid"
)

const (
	RedisRestaurantPrefix = "restaurant:"
)

type Cache interface {
	GetByID(ctx context.Context, id uuid.UUID) (*Restaurants, error)
	GetBySlug(ctx context.Context, slug string) (*Restaurants, error)
	Set(ctx context.Context, entity *Restaurants) error
	Delete(ctx context.Context, id uuid.UUID) error
	Update(ctx context.Context, entity *Restaurants) error
	MapSlugToID(ctx context.Context, id uuid.UUID, slug string) error
	DeleteSlugToID(ctx context.Context, slug string) error
}
