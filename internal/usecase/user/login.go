package user

import (
	"context"

	"github.com/james-wukong/orders-api/internal/domain/user"
	uSession "github.com/james-wukong/orders-api/internal/domain/user_session"
	"github.com/james-wukong/orders-api/internal/interfaces/http/dto"
)

type LoginUseCase struct {
	uRepo  user.Repository
	sRepo  uSession.Repository
	hasher user.PasswordHasher
}

func NewLoginUseCase(
	uRepo user.Repository,
	sRepo uSession.Repository,
	hasher user.PasswordHasher,
) *LoginUseCase {
	return &LoginUseCase{uRepo: uRepo, sRepo: sRepo, hasher: hasher}
}

func (uc *LoginUseCase) Execute(ctx context.Context, req *dto.LoginRequest,
) (*dto.LoginResponse, error) {
	// Fetch the user (including hash) from Postgres.

	// Compare the password using bcrypt.CompareHashAndPassword.

	// If successful, generate a random token.

	// Store the UserSession (without the hash) in Redis.

	// Return the token to the user.

	// // 1. Retrieve user by email
	// u, err := uc.uRepo.GetByEmail(ctx, req.Email)
	// if err != nil || u == nil {
	// 	return nil, user.ErrInvalidCredentials
	// }
	// // 2. Retrieve user session by user ID
	// // session, err := uc.sRepo.GetByUserID(ctx, u.ID)

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
