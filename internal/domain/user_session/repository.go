package usersession

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, session *UserSessions) error
	Update(ctx context.Context, session *UserSessions) error
	GetByToken(ctx context.Context, token string) (*UserSessions, error)
	GetByUserID(ctx context.Context, id uuid.UUID) ([]UserSessions, error)
	DeleteByToken(ctx context.Context, token string) error
	DeleteByUserID(ctx context.Context, id uuid.UUID) error
}
