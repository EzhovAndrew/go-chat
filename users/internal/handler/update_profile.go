package handler

import (
	"context"

	"github.com/go-chat/users/internal/domain"
	"github.com/go-chat/users/internal/dto"
	usersv1 "github.com/go-chat/users/pkg/api/users/v1"
)

// UpdateProfile updates an existing user profile
func (s *Server) UpdateProfile(ctx context.Context, req *usersv1.UpdateProfileRequest) (*usersv1.UpdateProfileResponse, error) {
	// Delegate to service layer
	profile, err := s.userService.UpdateProfile(
		ctx,
		domain.NewUserID(req.UserId),
		req.Nickname,
		req.Bio,
		req.AvatarUrl,
	)
	if err != nil {
		return nil, err // Middleware will map domain error to gRPC status
	}

	// Convert domain model to proto message
	return &usersv1.UpdateProfileResponse{
		Profile: dto.ToProtoProfile(profile),
	}, nil
}
