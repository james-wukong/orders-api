// Package usersession contains use cases related to user session management, such as creating, validating, and deleting sessions.
package usersession

import (
	"context"

	uSession "github.com/james-wukong/orders-api/internal/domain/user_session"
)

type CreateSessionUseCase struct {
	repo uSession.Repository
}

func NewCreateSessionUseCase(repo uSession.Repository) *CreateSessionUseCase {
	return &CreateSessionUseCase{repo: repo}
}

func (uc *CreateSessionUseCase) Execute(
	ctx context.Context,
	session *uSession.UserSessions,
) error {
	return uc.repo.Create(ctx, session)
}
