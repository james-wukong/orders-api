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

// PasswordHasher defines the contract for password security
type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(hash, password string) bool
}
