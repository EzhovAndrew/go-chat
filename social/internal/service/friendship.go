package service

import (
	"context"

	"github.com/go-chat/social/internal/domain"
)

// FriendshipService handles friendship management
type FriendshipService interface {
	// ListFriends lists all friends of a user with pagination
	ListFriends(ctx context.Context, userID domain.UserID, cursor string, limit int32) ([]domain.UserID, string, error)

	// RemoveFriend removes a friendship between two users
	RemoveFriend(ctx context.Context, userID, friendID domain.UserID) error
}

