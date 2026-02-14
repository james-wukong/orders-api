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
	token  user.TokenManager
}

func NewLoginUseCase(
	uRepo user.Repository,
	sRepo uSession.Repository,
	hasher user.PasswordHasher,
	token user.TokenManager,
) *LoginUseCase {
	return &LoginUseCase{uRepo: uRepo, sRepo: sRepo, hasher: hasher, token: token}
}

func (uc *LoginUseCase) Execute(ctx context.Context,
	req *dto.LoginRequest,
	ip string,
	os string,
	device string,
	userAgent string,
	deviceInfo string,
) (*uSession.UserSessions, error) {
	// 1. Fetch the user (including hash) from Postgres.
	u, err := uc.uRepo.LoginByEmail(ctx, req.Email)
	if err != nil || u == nil {
		return nil, user.ErrInvalidCredentials
	}
	// 2. Compare the password using bcrypt.CompareHashAndPassword.
	if !uc.hasher.Compare(u.PasswordHash, req.Password) {
		return nil, user.ErrInvalidCredentials
	}
	// 3. If successful, generate a jwt token.
	token, err := uc.token.GenerateToken(u.ID)
	if err != nil {
		return nil, err
	}

	session := uSession.NewUserSession(u.ID,
		token, deviceInfo, ip, os, device, userAgent)
	session.Token = token

	// 4. Store the UserSession (without the hash) in Postgres and Redis.
	err = uc.sRepo.Create(ctx, session)
	if err != nil {
		return nil, err
	}
	// 5. Return the token to the user.
	return session, nil
}
