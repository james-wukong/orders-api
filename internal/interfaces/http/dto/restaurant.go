package dto

import (
	"time"

	"github.com/james-wukong/orders-api/internal/domain/restaurant"
)

// CreateRestaurantRequest is what the client sends (POST /register)
type CreateRestaurantRequest struct {
	Name        string  `json:"name" binding:"required,min=3"`
	Slug        string  `json:"slug" binding:"required"`
	Description string  `json:"description"`
	Phone       string  `json:"phone"`
	Email       string  `json:"email" binding:"omitempty,email"`
	Address     string  `json:"address"`
	City        string  `json:"city"`
	State       string  `json:"state"`
	PostalCode  string  `json:"postal_code"`
	Country     string  `json:"country"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`

	// We use pointers for fields with DB defaults so we don't
	// accidentally overwrite the default with a zero-value.
	OpeningTime           *string  `json:"opening_time" binding:"omitempty,datetime=15:04"`
	ClosingTime           *string  `json:"closing_time" binding:"omitempty,datetime=15:04"`
	DeliveryFee           *float64 `json:"delivery_fee"`
	MinimumOrder          *float64 `json:"minimum_order"`
	EstimatedDeliveryTime *int     `json:"estimated_delivery_time"`

	LogoURL   string `json:"logo_url"`
	BannerURL string `json:"banner_url"`
}

// RestaurantResponse is what we send back to the client
type RestaurantResponse struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	Slug         string  `json:"slug"`
	Description  string  `json:"description,omitempty"`
	Phone        string  `json:"phone,omitempty"`
	Email        string  `json:"email,omitempty"`
	Rating       float64 `json:"rating"`
	TotalReviews int     `json:"total_reviews"`
	IsOpen       bool    `json:"is_open"`
	OpeningTime  string  `json:"opening_time"`
	ClosingTime  string  `json:"closing_time"`
	Address      string  `json:"address,omitempty"`
	City         string  `json:"city,omitempty"`
	State        string  `json:"state,omitempty"`
	CreatedAt    string  `json:"created_at"`
}

func MapToRestaurantResponse(entity *restaurant.Restaurant) RestaurantResponse {
	return RestaurantResponse{
		ID:           entity.ID.String(),
		Name:         entity.Name,
		Slug:         entity.Slug,
		Description:  entity.Description,
		Phone:        entity.Phone,
		Email:        entity.Email,
		TotalReviews: entity.TotalReviews,
		IsOpen:       entity.IsOpen,
		OpeningTime:  entity.OpeningTime,
		ClosingTime:  entity.ClosingTime,
		Rating:       entity.Rating,
		Address:      entity.Address,
		City:         entity.City,
		State:        entity.State,
		CreatedAt:    entity.CreatedAt.Format(time.RFC3339),
	}
}
