package handler

import (
	"context"
	"log"

	usersv1 "github.com/go-chat/users/pkg/api/users/v1"
)

// UpdateProfile updates an existing user profile
// Returns dummy updated profile until database integration is added
func (s *Server) UpdateProfile(ctx context.Context, req *usersv1.UpdateProfileRequest) (*usersv1.UpdateProfileResponse, error) {
	log.Printf("UpdateProfile called for user_id: %s", req.UserId)

	// TODO: Verify user_id matches authenticated user (extract from JWT in metadata)
	// TODO: Check if profile exists (NOT_FOUND error)
	// TODO: If nickname changed, validate format and uniqueness
	// TODO: Update profile in database (only non-nil fields)
	// TODO: Return ALREADY_EXISTS if new nickname is taken
	// TODO: Return NOT_FOUND if profile doesn't exist

	return &usersv1.UpdateProfileResponse{
		Profile: &usersv1.UserProfile{
			UserId:    req.UserId,
			Nickname:  req.Nickname,
			Bio:       req.Bio,
			AvatarUrl: req.AvatarUrl,
		},
	}, nil
}

