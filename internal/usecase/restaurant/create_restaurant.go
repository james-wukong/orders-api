// Package restaurant contains the use case for creating a restaurant.
// It orchestrates the business logic and interacts with the repository to persist the new restaurant entity.
package restaurant

import (
	"context"
	"fmt"

	"github.com/james-wukong/orders-api/internal/domain/restaurant"
	"github.com/james-wukong/orders-api/internal/interfaces/http/dto"
)

type CreateRestaurantUseCase struct {
	repo restaurant.Repository
}

func NewCreateRestaurantUseCase(repo restaurant.Repository) *CreateRestaurantUseCase {
	return &CreateRestaurantUseCase{
		repo: repo,
	}
}

func (uc *CreateRestaurantUseCase) Execute(ctx context.Context, input dto.CreateRestaurantRequest) (*restaurant.Restaurant, error) {
	// 1. Validate if slug is unique
	existing, err := uc.repo.GetBySlug(ctx, input.Slug)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, fmt.Errorf("restaurant with slug %s already exists", input.Slug)
	}

	// 2. Initialize Entity using the Factory
	res := restaurant.NewRestaurant(input.Name, input.Slug)

	// 3. Map Basic Fields
	res.Description = input.Description
	res.Phone = input.Phone
	res.Email = input.Email
	res.Address = input.Address
	res.City = input.City
	res.State = input.State
	res.PostalCode = input.PostalCode
	res.Country = input.Country
	res.Latitude = input.Latitude
	res.Longitude = input.Longitude
	res.LogoURL = input.LogoURL
	res.BannerURL = input.BannerURL

	// 4. Map Optional Fields with Pointer Logic
	// This prevents the "" error by only assigning if a value was actually provided
	if input.OpeningTime != nil && *input.OpeningTime != "" {
		res.OpeningTime = *input.OpeningTime
	}
	if input.ClosingTime != nil && *input.ClosingTime != "" {
		res.ClosingTime = *input.ClosingTime
	}
	if input.DeliveryFee != nil {
		res.DeliveryFee = *input.DeliveryFee
	}
	if input.MinimumOrder != nil {
		res.MinimumOrder = *input.MinimumOrder
	}
	if input.EstimatedDeliveryTime != nil {
		res.EstimatedDeliveryTime = *input.EstimatedDeliveryTime
	}

	// 5. Save to Repository
	if err := uc.repo.Create(ctx, res); err != nil {
		return nil, fmt.Errorf("failed to save restaurant: %w", err)
	}

	return res, nil
}
