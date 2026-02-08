// Package dto defines the Data Transfer Objects for the User entity.
// DTOs represent how data looks in an HTTP request or response.
package dto

import (
	"time"

	"github.com/james-wukong/orders-api/internal/domain/user"
)

// CreateUserRequest is used for registration/onboarding
type CreateUserRequest struct {
	Email     string `json:"email" binding:"required,email"`
	Password  string `json:"password" binding:"required,min=8"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Phone     string `json:"phone"`
	Role      string `json:"role" binding:"oneof=customer admin kitchen delivery inventory_manager"`
}

// UpdateUserRequest uses pointers to allow partial updates (PATCH)
type UpdateUserRequest struct {
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
	Phone     *string `json:"phone"`
	IsActive  *bool   `json:"is_active"`
}

type UserResponse struct {
	ID            string     `json:"id"`
	Email         string     `json:"email"`
	FirstName     string     `json:"first_name"`
	LastName      string     `json:"last_name"`
	Phone         string     `json:"phone,omitempty"`
	Role          string     `json:"role"`
	IsActive      bool       `json:"is_active"`
	EmailVerified bool       `json:"email_verified"`
	LastLogin     *time.Time `json:"last_login,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

// MapToUserResponse maps domain entity to response DTO
func MapToUserResponse(user *user.UserEntity) UserResponse {
	return UserResponse{
		ID:            user.ID.String(),
		Email:         user.Email,
		FirstName:     user.FirstName,
		LastName:      user.LastName,
		Phone:         user.Phone,
		Role:          string(user.Role),
		IsActive:      user.IsActive,
		EmailVerified: user.EmailVerified,
		LastLogin:     user.LastLogin,
		CreatedAt:     user.CreatedAt,
		UpdatedAt:     user.UpdatedAt,
	}
}

// QueryUserRequest is used for filtering and paginating user queries.
// Gin can bind both, but json is more consistent with the rest of the DTOs.
// TODO: decide to use json or form tags for query parameters.
type QueryUserRequest struct {
	Email     *string  `json:"email"`
	Roles     []string `json:"role"`
	FirstName *string  `json:"first_name"`
	LastName  *string  `json:"last_name"`
	IsActive  *bool    `json:"is_active"`
	Page      int      `json:"page"`
	Limit     int      `json:"limit"`
}
