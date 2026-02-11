// Package security provides password hashing and verification using bcrypt.
package security

import (
	"golang.org/x/crypto/bcrypt"
)

const (
	// DefaultBcryptCost is kept here because it's specific to this implementation
	DefaultBcryptCost = 12
)

type bcryptHasher struct {
	cost int
}

// NewBcryptHasher creates a new instance with a specific complexity cost
func NewBcryptHasher(cost int) *bcryptHasher {
	if cost == 0 {
		cost = bcrypt.DefaultCost
	}
	return &bcryptHasher{cost: cost}
}

func (h *bcryptHasher) Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), h.cost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (h *bcryptHasher) Compare(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
