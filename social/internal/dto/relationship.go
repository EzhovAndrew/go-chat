package dto

import (
	"github.com/go-chat/social/internal/domain"
	socialv1 "github.com/go-chat/social/pkg/api/social/v1"
)

// ToProtoRelationshipStatus converts domain.RelationshipStatus to proto RelationshipStatus
func ToProtoRelationshipStatus(status domain.RelationshipStatus) socialv1.RelationshipStatus {
	switch status {
	case domain.RelationshipStatusNone:
		return socialv1.RelationshipStatus_RELATIONSHIP_STATUS_NONE
	case domain.RelationshipStatusPending:
		return socialv1.RelationshipStatus_RELATIONSHIP_STATUS_PENDING
	case domain.RelationshipStatusFriend:
		return socialv1.RelationshipStatus_RELATIONSHIP_STATUS_FRIEND
	case domain.RelationshipStatusBlocked:
		return socialv1.RelationshipStatus_RELATIONSHIP_STATUS_BLOCKED
	default:
		return socialv1.RelationshipStatus_RELATIONSHIP_STATUS_UNSPECIFIED
	}
}

