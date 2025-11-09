package domain

import "errors"

var (
	// ErrChatNotFound is returned when a chat cannot be found
	ErrChatNotFound = errors.New("chat not found")

	// ErrChatAlreadyExists is returned when attempting to create a chat that already exists
	ErrChatAlreadyExists = errors.New("chat already exists")

	// ErrPermissionDenied is returned when a user lacks permission to access a resource
	ErrPermissionDenied = errors.New("permission denied")

	// ErrInvalidMessage is returned when a message is invalid (e.g., empty text)
	ErrInvalidMessage = errors.New("invalid message")

	// ErrInvalidChatID is returned when a chat ID is invalid
	ErrInvalidChatID = errors.New("invalid chat ID")

	// ErrUsersNotFriends is returned when attempting to create a chat with non-friends
	ErrUsersNotFriends = errors.New("users are not friends")

	// ErrUserBlocked is returned when attempting to interact with a blocked user
	ErrUserBlocked = errors.New("user is blocked")
)

