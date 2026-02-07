// Package restaurant defines the Restaurant entity and related value objects.
// It represents how data looks in your database or business rules.
package restaurant

import (
	"time"

	"github.com/google/uuid"
)

type Restaurant struct {
	ID                    uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name                  string    `gorm:"not null"`
	Slug                  string    `gorm:"uniqueIndex;not null"`
	Description           string
	Phone                 string
	Email                 string
	Address               string
	City                  string
	State                 string
	PostalCode            string
	Country               string  `gorm:"default:'USA'"`
	Latitude              float64 `gorm:"type:decimal(10,8)"`
	Longitude             float64 `gorm:"type:decimal(11,8)"`
	OpeningTime           string  `gorm:"type:time"`
	ClosingTime           string  `gorm:"type:time"`
	IsOpen                bool    `gorm:"default:true"`
	DeliveryFee           float64 `gorm:"type:decimal(10,2);default:0.00"`
	MinimumOrder          float64 `gorm:"type:decimal(10,2);default:0.00"`
	EstimatedDeliveryTime int     `gorm:"default:30"`
	LogoURL               string
	BannerURL             string
	Rating                float64 `gorm:"type:decimal(3,2);default:0.00"`
	TotalReviews          int     `gorm:"default:0"`
	CreatedAt             time.Time
	UpdatedAt             time.Time
}

// NewRestaurant is a factory function for creating a new restaurant instance
func NewRestaurant(name, slug string) *Restaurant {
	return &Restaurant{
		ID:        uuid.New(),
		Name:      name,
		Slug:      slug,
		IsOpen:    true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
