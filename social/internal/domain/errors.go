package domain

import "errors"

var (
	// ErrRequestNotFound is returned when a friend request cannot be found
	ErrRequestNotFound = errors.New("friend request not found")

	// ErrRequestAlreadyExists is returned when attempting to create a duplicate friend request
	ErrRequestAlreadyExists = errors.New("friend request already exists")

	// ErrAlreadyFriends is returned when attempting to send a friend request to an existing friend
	ErrAlreadyFriends = errors.New("users are already friends")

	// ErrNotFriends is returned when attempting an action that requires friendship
	ErrNotFriends = errors.New("users are not friends")

	// ErrUserBlocked is returned when attempting to interact with a blocked user
	ErrUserBlocked = errors.New("user is blocked")

	// ErrPermissionDenied is returned when a user lacks permission to perform an action
	ErrPermissionDenied = errors.New("permission denied")

	// ErrSelfAction is returned when attempting to perform an action on oneself
	ErrSelfAction = errors.New("cannot perform this action on yourself")
)

