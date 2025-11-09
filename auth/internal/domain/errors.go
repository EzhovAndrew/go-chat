package domain

import "errors"

var (
	// ErrEmailAlreadyExists is returned when attempting to register with an existing email
	ErrEmailAlreadyExists = errors.New("email already exists")

	// ErrInvalidCredentials is returned when login credentials are invalid
	ErrInvalidCredentials = errors.New("invalid credentials")

	// ErrInvalidToken is returned when a token is malformed or invalid
	ErrInvalidToken = errors.New("invalid token")

	// ErrTokenExpired is returned when a token has expired
	ErrTokenExpired = errors.New("token expired")

	// ErrTokenRevoked is returned when a token has been revoked
	ErrTokenRevoked = errors.New("token revoked")
)

