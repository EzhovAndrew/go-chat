package domain

import "errors"

var (
	// ErrNotificationNotFound is returned when a notification cannot be found
	ErrNotificationNotFound = errors.New("notification not found")

	// ErrPermissionDenied is returned when a user lacks permission to access a notification
	ErrPermissionDenied = errors.New("permission denied")
)

