// Package user contains the use case for creating a new user in the system.
// It validates the input, checks for email uniqueness, hashes the password, and saves the new user to the repository.
package user

import (
	"context"
	"errors"
	"fmt"

	"github.com/james-wukong/orders-api/internal/domain/user"
	"github.com/james-wukong/orders-api/internal/interfaces/http/dto"
)

type CreateUserUseCase struct {
	repo   user.Repository
	hasher user.PasswordHasher
}

func NewCreateUserUseCase(
	repo user.Repository,
	hasher user.PasswordHasher,
) *CreateUserUseCase {
	return &CreateUserUseCase{repo: repo, hasher: hasher}
}

func (uc *CreateUserUseCase) Execute(
	ctx context.Context,
	req dto.CreateUserRequest,
) (*user.Users, error) {
	// 1. Validate if email is unique
	existing, err := uc.repo.GetByEmail(ctx, req.Email)
	if err != nil {
		if !errors.Is(err, user.ErrUserNotFound) {
			return nil, err
		}
	}
	if existing != nil {
		return nil, fmt.Errorf("user with email %s already exists", req.Email)
	}
	// 2. Hash Password
	hashedPassword, err := uc.hasher.Hash(req.Password)
	if err != nil {
		return nil, err
	}
	// 3. Initialize Entity using the Factory
	u := user.NewUsers(req.Email, hashedPassword)

	// 4. Map Fields
	if req.FirstName != nil && *req.FirstName != "" {
		u.FirstName = *req.FirstName
	}
	if req.LastName != nil && *req.LastName != "" {
		u.LastName = *req.LastName
	}
	if req.Phone != nil && *req.Phone != "" {
		u.Phone = *req.Phone
	}
	if req.Role != "" {
		req.Role = string(user.RoleCustomer)
	}

	u.Role = user.Role(req.Role)

	// 5. Save to Repository
	if err := uc.repo.Create(ctx, u); err != nil {
		return nil, err
	}
	return u, nil
}
