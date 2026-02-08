package user

import (
	"context"
	"fmt"

	"github.com/james-wukong/orders-api/internal/domain/user"
	"github.com/james-wukong/orders-api/internal/interfaces/http/dto"
)

type CreateUserUseCase struct {
	repo   user.UserRepository
	hasher user.PasswordHasher
}

func NewCreateUserUseCase(
	repo user.UserRepository,
	hasher user.PasswordHasher,
) *CreateUserUseCase {
	return &CreateUserUseCase{repo: repo, hasher: hasher}
}

func (uc *CreateUserUseCase) Execute(
	ctx context.Context,
	req dto.CreateUserRequest,
) (*user.UserEntity, error) {
	// 1. Validate if email is unique
	existing, err := uc.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, fmt.Errorf("user with email %s already exists", req.Email)
	}
	// 2. Hash Password
	hashedPassword, err := uc.hasher.Hash(req.Password)
	if err != nil {
		return nil, err
	}
	// 2. Initialize Entity using the Factory
	u := user.NewUserEntity(req.Email, hashedPassword)

	// 3. Map Basic Fields
	u.FirstName = req.FirstName
	u.LastName = req.LastName
	u.Phone = req.Phone

	// 4. Save to Repository
	if err := uc.repo.Create(ctx, u); err != nil {
		return nil, err
	}
	return u, nil
}
