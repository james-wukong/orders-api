// Package restaurant defines the Restaurant entity and related value objects.
// It represents how data looks in your database or business rules.
package restaurant

import (
	"time"

	"github.com/google/uuid"
)

type Restaurant struct {
	ID                    uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name                  string    `gorm:"size:255;not null"`
	Slug                  string    `gorm:"size:255;unique;not null"`
	Description           string    `gorm:"type:text"`
	Phone                 string    `gorm:"size:20"`
	Email                 string    `gorm:"size:255"`
	Address               string    `gorm:"type:text"`
	City                  string    `gorm:"size:100"`
	State                 string    `gorm:"size:100"`
	PostalCode            string    `gorm:"size:20"`
	Country               string    `gorm:"size:100;default:'USA'"`
	Latitude              float64   `gorm:"type:decimal(10,8)"`
	Longitude             float64   `gorm:"type:decimal(11,8)"`
	OpeningTime           string    `gorm:"type:time;default:'10:00'"`
	ClosingTime           string    `gorm:"type:time;default:'22:00'"`
	IsOpen                bool      `gorm:"default:true"`
	DeliveryFee           float64   `gorm:"type:decimal(10,2);default:0.00"`
	MinimumOrder          float64   `gorm:"type:decimal(10,2);default:0.00"`
	EstimatedDeliveryTime int       `gorm:"default:30"`
	LogoURL               string    `gorm:"type:text"`
	BannerURL             string    `gorm:"type:text"`
	Rating                float64   `gorm:"type:decimal(3,2);default:0.00"`
	TotalReviews          int       `gorm:"default:0"`
	CreatedAt             time.Time `gorm:"autoCreateTime"`
	UpdatedAt             time.Time `gorm:"autoUpdateTime"`
}

// NewRestaurant is a Factory Function that ensures a Restaurant
// is always created with a valid ID and default business state.
func NewRestaurant(name, slug string) *Restaurant {
	return &Restaurant{
		ID:                    uuid.New(),
		Name:                  name,
		Slug:                  slug,
		IsOpen:                true,
		Country:               "USA",
		OpeningTime:           "10:00",
		ClosingTime:           "22:00",
		DeliveryFee:           0.00,
		MinimumOrder:          0.00,
		EstimatedDeliveryTime: 30,
	}
}
