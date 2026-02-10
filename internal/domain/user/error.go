// Package user defines error types related to user operations, such as authentication, validation, and repository errors.
package user

import "errors"

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidEmailFormat = errors.New("invalid email format")
	ErrPasswordTooWeak    = errors.New("password does not meet complexity requirements")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrEmailAlreadyUsed   = errors.New("email already used")
)
