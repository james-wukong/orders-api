// Package restaurant contains the use case for creating a restaurant.
// It orchestrates the business logic and interacts with the repository to persist the new restaurant entity.
package restaurant

import (
	"context"
	"errors"
	"fmt"

	"github.com/james-wukong/orders-api/internal/domain/restaurant"
)

// CreateRestaurantInput defines the data needed to create a restaurant
// This is often mapped from the DTO in the handler
type CreateRestaurantInput struct {
	Name        string
	Slug        string
	Description string
	Phone       string
	Email       string
	Address     string
	City        string
	State       string
}

type CreateRestaurantUseCase struct {
	repo restaurant.Repository
}

func NewCreateRestaurantUseCase(repo restaurant.Repository) *CreateRestaurantUseCase {
	return &CreateRestaurantUseCase{
		repo: repo,
	}
}

func (uc *CreateRestaurantUseCase) Execute(ctx context.Context, input CreateRestaurantInput) (*restaurant.Restaurant, error) {
	// 1. Business Rule: Check if a restaurant with the same slug already exists
	existing, err := uc.repo.GetBySlug(ctx, input.Slug)
	if err != nil {
		return nil, fmt.Errorf("error checking existing restaurant: %w", err)
	}
	if existing != nil {
		return nil, errors.New("a restaurant with this slug already exists")
	}

	// 2. Map input to Domain Entity
	// We use the NewRestaurant factory to ensure default values (like UUID) are set
	res := restaurant.NewRestaurant(input.Name, input.Slug)
	res.Description = input.Description
	res.Phone = input.Phone
	res.Email = input.Email
	res.Address = input.Address
	res.City = input.City
	res.State = input.State

	// 3. Persist the entity using the repository
	if err := uc.repo.Create(ctx, res); err != nil {
		return nil, fmt.Errorf("could not create restaurant: %w", err)
	}

	return res, nil
}
