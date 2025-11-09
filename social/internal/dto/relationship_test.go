package dto

import (
	"testing"

	"github.com/go-chat/social/internal/domain"
)

func TestToProtoRelationshipStatus_AllStatuses_MapCorrectly(t *testing.T) {
	tests := []struct {
		name         string
		domainStatus domain.RelationshipStatus
		expectedEnum int32
	}{
		{"None status", domain.RelationshipStatusNone, 1},
		{"Pending status", domain.RelationshipStatusPending, 2},
		{"Friend status", domain.RelationshipStatusFriend, 3},
		{"Blocked status", domain.RelationshipStatusBlocked, 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			protoStatus := ToProtoRelationshipStatus(tt.domainStatus)

			if int32(protoStatus) != tt.expectedEnum {
				t.Errorf("Expected enum %d, got %d", tt.expectedEnum, protoStatus)
			}
		})
	}
}

func TestToProtoRelationshipStatus_InvalidStatus_ReturnsUnspecified(t *testing.T) {
	// Test with an invalid status value
	invalidStatus := domain.RelationshipStatus(999)

	protoStatus := ToProtoRelationshipStatus(invalidStatus)

	if int32(protoStatus) != 0 {
		t.Errorf("Expected unspecified status (0), got %d", protoStatus)
	}
}

