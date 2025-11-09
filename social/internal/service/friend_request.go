package service

import (
	"context"

	"github.com/go-chat/social/internal/domain"
)

// FriendRequestService handles friend request operations
type FriendRequestService interface {
	// SendFriendRequest sends a friend request to another user
	SendFriendRequest(ctx context.Context, requesterID, targetID domain.UserID) (*domain.FriendRequest, error)

	// ListRequests lists pending friend requests for a user with pagination
	ListRequests(ctx context.Context, userID domain.UserID, cursor string, limit int32) ([]*domain.FriendRequest, string, error)

	// AcceptFriendRequest accepts a pending friend request
	AcceptFriendRequest(ctx context.Context, requestID domain.RequestID, userID domain.UserID) (*domain.FriendRequest, error)

	// DeclineFriendRequest declines a pending friend request
	DeclineFriendRequest(ctx context.Context, requestID domain.RequestID, userID domain.UserID) error
}

