// Package user defines the user domain, including entities, value objects, repositories, and services.
package user

import (
	"context"

	"github.com/google/uuid"
)

type Repository interface {
	// Create creates a new user in the repository and returns the created user or an error if the operation fails.
	Create(ctx context.Context, user *Users) error

	// GetByID retrieves a user by their unique identifier. It returns the user or an error if the user is not found.
	GetByID(ctx context.Context, id uuid.UUID) (*Users, error)

	// GetByEmail retrieves a user by their email address. It returns the user or an error if the user is not found.
	GetByEmail(ctx context.Context, email string) (*Users, error)
	LoginByEmail(ctx context.Context, email string) (*UserLogin, error)

	// Update updates an existing user's information in the repository. It returns the updated user or an error if the operation fails.
	Update(ctx context.Context, user *Users) error

	// Delete removes a user from the repository by their unique identifier. It returns an error if the operation fails.
	Delete(ctx context.Context, id uuid.UUID) error

	// List retrieves all users from the repository. It returns a slice of users or an error if the operation fails.
	List(ctx context.Context, filter *UserFilterEntity) ([]Users, error)
}
