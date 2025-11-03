package handler

import (
	"context"
	"log"

	usersv1 "github.com/go-chat/users/pkg/api/users/v1"
)

// GetProfileByID retrieves a profile by user ID
// Returns dummy profile until database integration is added
func (s *Server) GetProfileByID(ctx context.Context, req *usersv1.GetProfileByIDRequest) (*usersv1.GetProfileByIDResponse, error) {
	log.Printf("GetProfileByID called for user_id: %s", req.UserId)

	// TODO: Query profile from database by user_id
	// TODO: Return NOT_FOUND error if profile doesn't exist

	return &usersv1.GetProfileByIDResponse{
		Profile: &usersv1.UserProfile{
			UserId:    req.UserId,
			Nickname:  "dummy_user",
			Bio:       "This is a dummy bio",
			AvatarUrl: "https://example.com/avatar.jpg",
		},
	}, nil
}

// GetProfilesByIDs retrieves multiple profiles by user IDs (batch operation)
// Returns dummy profiles until database integration is added
func (s *Server) GetProfilesByIDs(ctx context.Context, req *usersv1.GetProfilesByIDsRequest) (*usersv1.GetProfilesByIDsResponse, error) {
	log.Printf("GetProfilesByIDs called with %d user_ids", len(req.UserIds))

	// TODO: Query profiles from database by user_ids (IN query or batch)
	// TODO: Return only found profiles (no error for missing ones)
	// TODO: Maintain order or return map for efficient lookup

	profiles := make([]*usersv1.UserProfile, 0, len(req.UserIds))
	for _, userID := range req.UserIds {
		profiles = append(profiles, &usersv1.UserProfile{
			UserId:    userID,
			Nickname:  "dummy_user",
			Bio:       "This is a dummy bio",
			AvatarUrl: "https://example.com/avatar.jpg",
		})
	}

	return &usersv1.GetProfilesByIDsResponse{
		Profiles: profiles,
	}, nil
}

// GetProfileByNickname retrieves a profile by nickname
// Returns dummy profile until database integration is added
func (s *Server) GetProfileByNickname(ctx context.Context, req *usersv1.GetProfileByNicknameRequest) (*usersv1.GetProfileByNicknameResponse, error) {
	log.Printf("GetProfileByNickname called for nickname: %s", req.Nickname)

	// TODO: Query profile from database by nickname (unique index)
	// TODO: Return NOT_FOUND error if nickname doesn't exist

	return &usersv1.GetProfileByNicknameResponse{
		Profile: &usersv1.UserProfile{
			UserId:    "550e8400-e29b-41d4-a716-446655440000",
			Nickname:  req.Nickname,
			Bio:       "This is a dummy bio",
			AvatarUrl: "https://example.com/avatar.jpg",
		},
	}, nil
}

