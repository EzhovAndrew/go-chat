package handler

import (
	"context"
	"testing"

	"github.com/go-chat/social/internal/domain"
	socialv1 "github.com/go-chat/social/pkg/api/social/v1"
)

func TestCheckRelationship_ReturnsFriendStatus(t *testing.T) {
	mockRelSvc := &mockRelationshipService{
		checkRelationshipFunc: func(ctx context.Context, userID, targetID domain.UserID) (domain.RelationshipStatus, error) {
			return domain.RelationshipStatusFriend, nil
		},
	}

	server := NewServer(nil, nil, nil, mockRelSvc)
	req := &socialv1.CheckRelationshipRequest{
		UserId:       "550e8400-e29b-41d4-a716-446655440001",
		TargetUserId: "550e8400-e29b-41d4-a716-446655440002",
	}

	resp, err := server.CheckRelationship(context.Background(), req)

	if err != nil {
		t.Fatalf("CheckRelationship() returned error: %v", err)
	}

	if resp.Status != socialv1.RelationshipStatus_RELATIONSHIP_STATUS_FRIEND {
		t.Errorf("Expected status FRIEND, got %v", resp.Status)
	}
}

func TestCheckRelationship_ReturnsBlockedStatus(t *testing.T) {
	mockRelSvc := &mockRelationshipService{
		checkRelationshipFunc: func(ctx context.Context, userID, targetID domain.UserID) (domain.RelationshipStatus, error) {
			return domain.RelationshipStatusBlocked, nil
		},
	}

	server := NewServer(nil, nil, nil, mockRelSvc)
	req := &socialv1.CheckRelationshipRequest{
		UserId:       "550e8400-e29b-41d4-a716-446655440001",
		TargetUserId: "550e8400-e29b-41d4-a716-446655440002",
	}

	resp, err := server.CheckRelationship(context.Background(), req)

	if err != nil {
		t.Fatalf("CheckRelationship() returned error: %v", err)
	}

	if resp.Status != socialv1.RelationshipStatus_RELATIONSHIP_STATUS_BLOCKED {
		t.Errorf("Expected status BLOCKED, got %v", resp.Status)
	}
}

