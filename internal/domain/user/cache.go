package user

import (
	"context"

	"github.com/google/uuid"
)

const (
	RedisUserPrefix    = "user:"
	RedisSessionPrefix = "session:"
)

type Cache interface {
	GetByID(ctx context.Context, userID uuid.UUID) (*Users, error)
	GetByEmail(ctx context.Context, email string) (*Users, error)
	Set(ctx context.Context, user *Users) error
	Delete(ctx context.Context, userID uuid.UUID) error
	Update(ctx context.Context, user *Users) error
}
