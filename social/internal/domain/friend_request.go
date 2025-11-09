package domain

import "time"

// FriendRequestStatus represents the status of a friend request
type FriendRequestStatus string

const (
	// FriendRequestStatusPending indicates the request is awaiting response
	FriendRequestStatusPending FriendRequestStatus = "pending"

	// FriendRequestStatusAccepted indicates the request was accepted
	FriendRequestStatusAccepted FriendRequestStatus = "accepted"

	// FriendRequestStatusDeclined indicates the request was declined
	FriendRequestStatusDeclined FriendRequestStatus = "declined"
)

// FriendRequest represents a friend request between two users
type FriendRequest struct {
	RequestID   RequestID
	RequesterID UserID
	TargetID    UserID
	Status      FriendRequestStatus
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

