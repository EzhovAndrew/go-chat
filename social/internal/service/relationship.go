package service

import (
	"context"

	"github.com/go-chat/social/internal/domain"
)

// RelationshipService handles relationship queries
type RelationshipService interface {
	// CheckRelationship checks the relationship status between two users
	CheckRelationship(ctx context.Context, userID, targetID domain.UserID) (domain.RelationshipStatus, error)
}

