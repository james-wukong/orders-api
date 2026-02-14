// Package user defines the user domain services, which contain business logic related to user operations that may involve multiple entities or value objects.
// A Domain Service contains business logic that does not naturally belong to a single entity or value object,
// but is still pure domain logic.
// Create a domain service when ALL are true:
// ✔️ Logic is business-related
// ✔️ Logic involves multiple entities
// ✔️ Logic doesn’t fit one entity naturally
// ✔️ Logic must be reusable
// ✔️ Logic must be testable without DB
package user

import (
	"github.com/google/uuid"
	"github.com/james-wukong/orders-api/internal/infrastructure/security"
)

// PasswordHasher defines the contract for password security
type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(hash, password string) bool
}

// TokenManager define the contract for token related functions
type TokenManager interface {
	GenerateToken(userID uuid.UUID) (string, error)
	VerifyToken(token string) (*security.Claims, error)
}
