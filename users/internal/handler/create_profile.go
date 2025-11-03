package handler

import (
	"context"
	"log"

	usersv1 "github.com/go-chat/users/pkg/api/users/v1"
)

// CreateProfile creates a new user profile
// Returns dummy profile until database integration is added
func (s *Server) CreateProfile(ctx context.Context, req *usersv1.CreateProfileRequest) (*usersv1.CreateProfileResponse, error) {
	log.Printf("CreateProfile called for user_id: %s, nickname: %s", req.UserId, req.Nickname)

	// TODO: Validate nickname format (^[a-z0-9_]{3,20}$)
	// TODO: Check if profile already exists for user_id (ALREADY_EXISTS error)
	// TODO: Check if nickname is already taken (ALREADY_EXISTS error)
	// TODO: Store profile in database
	// TODO: Return INVALID_ARGUMENT for invalid inputs

	return &usersv1.CreateProfileResponse{
		Profile: &usersv1.UserProfile{
			UserId:    req.UserId,
			Nickname:  req.Nickname,
			Bio:       req.Bio,
			AvatarUrl: req.AvatarUrl,
		},
	}, nil
}

