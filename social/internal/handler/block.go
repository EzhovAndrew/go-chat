package handler

import (
	"context"

	"github.com/go-chat/social/internal/domain"
	socialv1 "github.com/go-chat/social/pkg/api/social/v1"
)

// BlockUser blocks a user
func (s *Server) BlockUser(ctx context.Context, req *socialv1.BlockUserRequest) (*socialv1.BlockUserResponse, error) {
	// TODO: Extract authenticated user_id from JWT
	blockerID := domain.NewUserID("authenticated-user-id") // Placeholder

	// Delegate to service layer
	err := s.blockService.BlockUser(
		ctx,
		blockerID,
		domain.NewUserID(req.TargetUserId),
	)
	if err != nil {
		return nil, err // Middleware will map domain error to gRPC status
	}

	// Return empty response on success (as per proto)
	return &socialv1.BlockUserResponse{}, nil
}

// UnblockUser unblocks a previously blocked user
func (s *Server) UnblockUser(ctx context.Context, req *socialv1.UnblockUserRequest) (*socialv1.UnblockUserResponse, error) {
	// TODO: Extract authenticated user_id from JWT
	blockerID := domain.NewUserID("authenticated-user-id") // Placeholder

	// Delegate to service layer
	err := s.blockService.UnblockUser(
		ctx,
		blockerID,
		domain.NewUserID(req.TargetUserId),
	)
	if err != nil {
		return nil, err // Middleware will map domain error to gRPC status
	}

	// Return empty response on success (as per proto)
	return &socialv1.UnblockUserResponse{}, nil
}
