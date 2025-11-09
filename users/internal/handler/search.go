package handler

import (
	"context"

	"github.com/go-chat/users/internal/dto"
	usersv1 "github.com/go-chat/users/pkg/api/users/v1"
)

// SearchByNickname searches for profiles matching a query with pagination
func (s *Server) SearchByNickname(ctx context.Context, req *usersv1.SearchByNicknameRequest) (*usersv1.SearchByNicknameResponse, error) {
	// Default limit if not provided
	limit := req.Limit
	if limit == 0 {
		limit = 20
	}

	// Delegate to service layer
	profiles, nextCursor, err := s.userService.SearchByNickname(ctx, req.Query, req.Cursor, limit)
	if err != nil {
		return nil, err // Middleware will map domain error to gRPC status
	}

	// Convert domain models to proto messages
	return &usersv1.SearchByNicknameResponse{
		Profiles:   dto.ToProtoProfiles(profiles),
		NextCursor: nextCursor,
	}, nil
}
