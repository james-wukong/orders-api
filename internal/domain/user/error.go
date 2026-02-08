// Package user defines error types related to user operations, such as authentication, validation, and repository errors.
package user

import "errors"

var (
	ErrUserNotFound     = errors.New("user not found")
	ErrEmailAlreadyUsed = errors.New("email already used")
)
