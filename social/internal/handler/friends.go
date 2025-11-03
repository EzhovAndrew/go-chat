package handler

import (
	"context"
	"log"

	socialv1 "github.com/go-chat/social/pkg/api/social/v1"
)

// RemoveFriend removes a user from friends list
// Returns empty response until database integration is added
func (s *Server) RemoveFriend(ctx context.Context, req *socialv1.RemoveFriendRequest) (*socialv1.RemoveFriendResponse, error) {
	log.Printf("RemoveFriend called with friend_user_id: %s", req.FriendUserId)

	// TODO: Extract authenticated user_id from JWT
	// TODO: Query friendship record from database
	// TODO: Return NOT_FOUND if friendship doesn't exist
	// TODO: Delete bidirectional friendship records
	// TODO: Return empty response

	return &socialv1.RemoveFriendResponse{}, nil
}

// ListFriends lists all friends of a user with cursor-based pagination
// Returns dummy friend IDs until database integration is added
func (s *Server) ListFriends(ctx context.Context, req *socialv1.ListFriendsRequest) (*socialv1.ListFriendsResponse, error) {
	log.Printf("ListFriends called for user_id: %s, limit: %d", req.UserId, req.Limit)

	// TODO: Extract authenticated user_id from JWT
	// TODO: Verify req.UserId matches authenticated user (or allow public access)
	// TODO: Decode cursor for pagination
	// TODO: Query friendships from database
	// TODO: Return user_ids of friends
	// TODO: Apply limit + 1 to check for more results
	// TODO: Encode next_cursor if more results exist

	return &socialv1.ListFriendsResponse{
		UserIds: []string{
			"550e8400-e29b-41d4-a716-446655440001",
			"550e8400-e29b-41d4-a716-446655440002",
			"550e8400-e29b-41d4-a716-446655440003",
		},
		NextCursor: "", // Empty means no more results
	}, nil
}

