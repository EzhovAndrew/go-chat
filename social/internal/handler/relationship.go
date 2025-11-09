package handler

import (
	"context"

	"github.com/go-chat/social/internal/domain"
	socialv1 "github.com/go-chat/social/pkg/api/social/v1"
)

// CheckRelationship checks the relationship status between two users
func (s *Server) CheckRelationship(ctx context.Context, req *socialv1.CheckRelationshipRequest) (*socialv1.CheckRelationshipResponse, error) {
	// Delegate to service layer
	status, err := s.relationshipService.CheckRelationship(
		ctx,
		domain.NewUserID(req.UserId),
		domain.NewUserID(req.TargetUserId),
	)
	if err != nil {
		return nil, err // Middleware will map domain error to gRPC status
	}

	// Convert domain status to proto enum
	var protoStatus socialv1.RelationshipStatus
	switch status {
	case domain.RelationshipStatusNone:
		protoStatus = socialv1.RelationshipStatus_RELATIONSHIP_STATUS_NONE
	case domain.RelationshipStatusPending:
		protoStatus = socialv1.RelationshipStatus_RELATIONSHIP_STATUS_PENDING
	case domain.RelationshipStatusFriend:
		protoStatus = socialv1.RelationshipStatus_RELATIONSHIP_STATUS_FRIEND
	case domain.RelationshipStatusBlocked:
		protoStatus = socialv1.RelationshipStatus_RELATIONSHIP_STATUS_BLOCKED
	default:
		protoStatus = socialv1.RelationshipStatus_RELATIONSHIP_STATUS_UNSPECIFIED
	}

	return &socialv1.CheckRelationshipResponse{
		Status: protoStatus,
	}, nil
}
