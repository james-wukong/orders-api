package user

import (
	"context"

	"github.com/james-wukong/orders-api/internal/domain/user"
	uSession "github.com/james-wukong/orders-api/internal/domain/user_session"
	"github.com/james-wukong/orders-api/internal/interfaces/http/dto"
	"github.com/redis/go-redis/v9"
)

type LoginUseCase struct {
	uRepo  user.Repository
	sRepo  uSession.Repository
	hasher user.PasswordHasher
	redis  *redis.Client
}

func NewLoginUseCase(
	uRepo user.Repository,
	sRepo uSession.Repository,
	hasher user.PasswordHasher,
	redis *redis.Client,
) *LoginUseCase {
	return &LoginUseCase{uRepo: uRepo, sRepo: sRepo, hasher: hasher, redis: redis}
}

func (uc *LoginUseCase) Execute(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error) {
	// // 1. Retrieve user by email from the redis cache
	// u, err := uc.uRepo.GetByEmail(ctx, req.Email)
	// if err != nil {
	// 	return nil, err
	// }
	// // 2. Retrieve user session by user ID from the database
	// session, err := uc.sRepo.GetByID(ctx, u.ID)

	// // 2. Verify password
	// if !uc.hasher.Compare(u.PasswordHash, req.Password) {
	// 	return nil, user.ErrInvalidCredentials
	// }

	// // 3. Return login response
	// return &dto.LoginResponse{
	// 	Token: u.Token,
	// }, nil
	return nil, nil
}
