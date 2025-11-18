package handler

import (
	"context"

	"github.com/go-chat/social/internal/domain"
	socialv1 "github.com/go-chat/social/pkg/api/social/v1"
)

// ListFriends lists all friends of a user
func (s *Server) ListFriends(ctx context.Context, req *socialv1.ListFriendsRequest) (*socialv1.ListFriendsResponse, error) {
	// Default limit if not provided
	limit := req.Limit
	if limit == 0 {
		limit = 20
	}

	// Delegate to service layer
	friendIDs, nextCursor, err := s.friendshipService.ListFriends(
		ctx,
		domain.NewUserID(req.UserId),
		req.Cursor,
		limit,
	)
	if err != nil {
		return nil, err // Middleware will map domain error to gRPC status
	}

	// Convert domain types to strings
	userIDs := make([]string, len(friendIDs))
	for i, id := range friendIDs {
		userIDs[i] = id.String()
	}

	return &socialv1.ListFriendsResponse{
		UserIds:    userIDs,
		NextCursor: nextCursor,
	}, nil
}

// RemoveFriend removes a user from friends list
func (s *Server) RemoveFriend(ctx context.Context, req *socialv1.RemoveFriendRequest) (*socialv1.RemoveFriendResponse, error) {
	// TODO: Extract authenticated user_id from JWT
	userID := domain.NewUserID("authenticated-user-id") // Placeholder

	// Delegate to service layer
	err := s.friendshipService.RemoveFriend(
		ctx,
		userID,
		domain.NewUserID(req.FriendUserId),
	)
	if err != nil {
		return nil, err // Middleware will map domain error to gRPC status
	}

	// Return empty response on success (as per proto)
	return &socialv1.RemoveFriendResponse{}, nil
}
