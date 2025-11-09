package handler

import (
	"context"

	"github.com/go-chat/users/internal/domain"
	usersv1 "github.com/go-chat/users/pkg/api/users/v1"
)

// CreateProfile creates a new user profile
func (s *Server) CreateProfile(ctx context.Context, req *usersv1.CreateProfileRequest) (*usersv1.CreateProfileResponse, error) {
	// Handle optional avatar_url field
	avatarURL := ""
	if req.AvatarUrl != nil {
		avatarURL = *req.AvatarUrl
	}

	// Delegate to service layer
	profile, err := s.userService.CreateProfile(
		ctx,
		domain.NewUserID(req.UserId),
		req.Nickname,
		req.Bio,
		avatarURL,
	)
	if err != nil {
		return nil, err // Middleware will map domain error to gRPC status
	}

	// Convert domain model to proto message
	return &usersv1.CreateProfileResponse{
		Profile: &usersv1.UserProfile{
			UserId:    profile.UserID.String(),
			Nickname:  profile.Nickname,
			Bio:       profile.Bio,
			AvatarUrl: profile.AvatarURL,
		},
	}, nil
}
