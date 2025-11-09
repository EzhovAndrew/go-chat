package domain

import "errors"

var (
	// ErrProfileAlreadyExists is returned when attempting to create a profile that already exists
	ErrProfileAlreadyExists = errors.New("profile already exists")

	// ErrProfileNotFound is returned when a profile cannot be found
	ErrProfileNotFound = errors.New("profile not found")

	// ErrNicknameAlreadyExists is returned when attempting to use a nickname that is already taken
	ErrNicknameAlreadyExists = errors.New("nickname already exists")

	// ErrInvalidNickname is returned when a nickname does not meet requirements
	ErrInvalidNickname = errors.New("invalid nickname format")
)
