package handler

import (
	"context"

	"github.com/go-chat/users/internal/domain"
	"github.com/go-chat/users/internal/dto"
	usersv1 "github.com/go-chat/users/pkg/api/users/v1"
)

// GetProfileByID retrieves a profile by user ID
func (s *Server) GetProfileByID(ctx context.Context, req *usersv1.GetProfileByIDRequest) (*usersv1.GetProfileByIDResponse, error) {
	// Delegate to service layer
	profile, err := s.userService.GetProfileByID(ctx, domain.NewUserID(req.UserId))
	if err != nil {
		return nil, err // Middleware will map domain error to gRPC status
	}

	// Convert domain model to proto message
	return &usersv1.GetProfileByIDResponse{
		Profile: dto.ToProtoProfile(profile),
	}, nil
}

// GetProfilesByIDs retrieves multiple profiles by user IDs (batch operation)
func (s *Server) GetProfilesByIDs(ctx context.Context, req *usersv1.GetProfilesByIDsRequest) (*usersv1.GetProfilesByIDsResponse, error) {
	// Convert request user IDs to domain types
	userIDs := make([]domain.UserID, len(req.UserIds))
	for i, id := range req.UserIds {
		userIDs[i] = domain.NewUserID(id)
	}

	// Delegate to service layer
	profiles, err := s.userService.GetProfilesByIDs(ctx, userIDs)
	if err != nil {
		return nil, err // Middleware will map domain error to gRPC status
	}

	// Convert domain models to proto messages
	return &usersv1.GetProfilesByIDsResponse{
		Profiles: dto.ToProtoProfiles(profiles),
	}, nil
}

// GetProfileByNickname retrieves a profile by nickname
func (s *Server) GetProfileByNickname(ctx context.Context, req *usersv1.GetProfileByNicknameRequest) (*usersv1.GetProfileByNicknameResponse, error) {
	// Delegate to service layer
	profile, err := s.userService.GetProfileByNickname(ctx, req.Nickname)
	if err != nil {
		return nil, err // Middleware will map domain error to gRPC status
	}

	// Convert domain model to proto message
	return &usersv1.GetProfileByNicknameResponse{
		Profile: dto.ToProtoProfile(profile),
	}, nil
}
