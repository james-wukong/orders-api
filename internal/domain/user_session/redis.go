package usersession

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const (
	RedisUserPrefix    = "user:"
	RedisSessionPrefix = "session:"
	RedisTokenPrefix   = "session:"
	RedisSessionTTL    = 2 * 24 * time.Hour
)

type RedisCache interface {
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]UserSessions, error)
	GetByToken(ctx context.Context, token string) (*UserSessions, error)
	Set(ctx context.Context, entity *UserSessions) error
	DeleteByUserID(ctx context.Context, userID uuid.UUID) error
	DeleteByToken(ctx context.Context, token string) error
	Update(ctx context.Context, entity *UserSessions) error
}
