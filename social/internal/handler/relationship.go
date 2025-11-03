package handler

import (
	"context"
	"log"

	socialv1 "github.com/go-chat/social/pkg/api/social/v1"
)

// CheckRelationship checks the relationship status between two users
// This is an internal gRPC endpoint used by Chat Service to verify permissions
// Returns dummy status until database integration is added
func (s *Server) CheckRelationship(ctx context.Context, req *socialv1.CheckRelationshipRequest) (*socialv1.CheckRelationshipResponse, error) {
	log.Printf("CheckRelationship called for user_id: %s, target_user_id: %s", req.UserId, req.TargetUserId)

	// TODO: Query database for relationship between users
	// TODO: Check if blocked (either direction) - return BLOCKED
	// TODO: Check if friends - return FRIEND
	// TODO: Check if pending friend request exists - return PENDING
	// TODO: If none of above - return NONE

	// For dummy implementation, always return FRIEND to allow chat creation
	return &socialv1.CheckRelationshipResponse{
		Status: socialv1.RelationshipStatus_RELATIONSHIP_STATUS_FRIEND,
	}, nil
}

