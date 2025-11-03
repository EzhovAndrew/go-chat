package handler

import (
	"context"
	"log"

	usersv1 "github.com/go-chat/users/pkg/api/users/v1"
)

// SearchByNickname searches for profiles matching a query with cursor-based pagination
// Returns dummy search results until database integration is added
func (s *Server) SearchByNickname(ctx context.Context, req *usersv1.SearchByNicknameRequest) (*usersv1.SearchByNicknameResponse, error) {
	log.Printf("SearchByNickname called with query: %s, limit: %d", req.Query, req.Limit)

	// TODO: Implement full-text search or LIKE query on nickname
	// TODO: Use cursor-based pagination (cursor = last nickname + user_id)
	// TODO: Decode cursor to get starting position
	// TODO: Apply limit (with max cap, e.g., 100)
	// TODO: Query database with LIMIT + 1 to check if more results exist
	// TODO: If more results exist, encode next_cursor from last result
	// TODO: Return empty array if no matches (not an error)

	// Dummy response with 2 profiles
	return &usersv1.SearchByNicknameResponse{
		Profiles: []*usersv1.UserProfile{
			{
				UserId:    "550e8400-e29b-41d4-a716-446655440001",
				Nickname:  "user_one",
				Bio:       "First dummy user",
				AvatarUrl: "https://example.com/avatar1.jpg",
			},
			{
				UserId:    "550e8400-e29b-41d4-a716-446655440002",
				Nickname:  "user_two",
				Bio:       "Second dummy user",
				AvatarUrl: "https://example.com/avatar2.jpg",
			},
		},
		NextCursor: "", // Empty means no more results
	}, nil
}

