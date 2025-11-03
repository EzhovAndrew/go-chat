package handler

import (
	"context"
	"log"

	socialv1 "github.com/go-chat/social/pkg/api/social/v1"
)

// BlockUser blocks a user
// Returns empty response until database integration is added
func (s *Server) BlockUser(ctx context.Context, req *socialv1.BlockUserRequest) (*socialv1.BlockUserResponse, error) {
	log.Printf("BlockUser called with target_user_id: %s", req.TargetUserId)

	// TODO: Extract authenticated user_id (blocker_id) from JWT
	// TODO: Check if target user exists (call User Service or return NOT_FOUND)
	// TODO: Check if already blocked (idempotent operation)
	// TODO: Remove friendship if exists (blocking removes friend)
	// TODO: Create block record in database
	// TODO: Return empty response

	return &socialv1.BlockUserResponse{}, nil
}

// UnblockUser unblocks a user
// Returns empty response until database integration is added
func (s *Server) UnblockUser(ctx context.Context, req *socialv1.UnblockUserRequest) (*socialv1.UnblockUserResponse, error) {
	log.Printf("UnblockUser called with target_user_id: %s", req.TargetUserId)

	// TODO: Extract authenticated user_id from JWT
	// TODO: Query block record from database
	// TODO: Return NOT_FOUND if block doesn't exist
	// TODO: Delete block record
	// TODO: Return empty response

	return &socialv1.UnblockUserResponse{}, nil
}

