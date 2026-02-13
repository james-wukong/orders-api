package user

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const (
	RedisUserPrefix      = "user:"
	RedisEmailToIDPrefix = "email_to_id:"
	RedisSessionPrefix   = "session:"
	RedisUerTTL          = 2 * 24 * time.Hour
)

type RedisCache interface {
	GetByID(ctx context.Context, userID uuid.UUID) (*Users, error)
	GetByEmail(ctx context.Context, email string) (*Users, error)
	Set(ctx context.Context, user *Users) error
	Delete(ctx context.Context, userID uuid.UUID) error
	Update(ctx context.Context, user *Users) error
	MapEmailToID(ctx context.Context, email string, id uuid.UUID) error
	DeleteEmailToID(ctx context.Context, email string) error
}
