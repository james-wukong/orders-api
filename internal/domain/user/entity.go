// Package user defines the User entity and related types.
package user

import (
	"time"

	"github.com/google/uuid"
)

// Role defines a custom type for the user_role_enum
type Role string

const (
	RoleCustomer         Role = "customer"
	RoleAdmin            Role = "admin"
	RoleKitchen          Role = "kitchen"
	RoleDelivery         Role = "delivery"
	RoleInventoryManager Role = "inventory_manager"
)

type UserEntity struct {
	ID                     uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Email                  string     `gorm:"size:255;unique;not null"`
	PasswordHash           string     `gorm:"size:255;not null"`
	FirstName              string     `gorm:"size:100"`
	LastName               string     `gorm:"size:100"`
	Phone                  string     `gorm:"size:20"`
	Role                   Role       `gorm:"type:user_role_enum;default:'customer'"`
	IsActive               bool       `gorm:"default:true"`
	EmailVerified          bool       `gorm:"default:false"`
	EmailVerificationToken *string    `gorm:"size:255"` // Pointer: can be null if verified
	PasswordResetToken     *string    `gorm:"size:255"` // Pointer: can be null
	PasswordResetExpires   *time.Time `gorm:"type:timestamp"`
	LastLogin              *time.Time `gorm:"type:timestamp"`
	CreatedAt              time.Time  `gorm:"autoCreateTime"`
	UpdatedAt              time.Time  `gorm:"autoUpdateTime"`
}

type UserFilterEntity struct {
	Email     *string
	Roles     []Role
	IsActive  *bool
	FirstName *string
	LastName  *string
	Page      int
	Limit     int
}

// NewUserEntity is a factory function to initialize a new user entity with defaults
func NewUserEntity(email, passwordHash string) *UserEntity {
	return &UserEntity{
		ID:            uuid.New(),
		Email:         email,
		PasswordHash:  passwordHash,
		Role:          RoleCustomer,
		IsActive:      true,
		EmailVerified: false,
	}
}
