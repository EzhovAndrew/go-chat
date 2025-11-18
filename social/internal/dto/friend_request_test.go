package dto

import (
	"testing"
	"time"

	"github.com/go-chat/social/internal/domain"
)

func TestToProtoFriendRequest_ValidRequest_ConvertsCorrectly(t *testing.T) {
	fr := &domain.FriendRequest{
		RequestID:   domain.NewRequestID("550e8400-e29b-41d4-a716-446655440000"),
		RequesterID: domain.NewUserID("550e8400-e29b-41d4-a716-446655440001"),
		TargetID:    domain.NewUserID("550e8400-e29b-41d4-a716-446655440002"),
		Status:      domain.FriendRequestStatusPending,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	protoFR := ToProtoFriendRequest(fr)

	if protoFR.RequestId != "550e8400-e29b-41d4-a716-446655440000" {
		t.Errorf("Expected RequestId '550e8400-e29b-41d4-a716-446655440000', got '%s'", protoFR.RequestId)
	}

	if protoFR.RequesterId != "550e8400-e29b-41d4-a716-446655440001" {
		t.Errorf("Expected RequesterId '550e8400-e29b-41d4-a716-446655440001', got '%s'", protoFR.RequesterId)
	}

	if protoFR.Status != "pending" {
		t.Errorf("Expected Status 'pending', got '%s'", protoFR.Status)
	}
}

func TestToProtoFriendRequests_MultipleRequests_ConvertsAll(t *testing.T) {
	requests := []*domain.FriendRequest{
		{
			RequestID:   domain.NewRequestID("req1"),
			RequesterID: domain.NewUserID("user1"),
			TargetID:    domain.NewUserID("user2"),
			Status:      domain.FriendRequestStatusPending,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			RequestID:   domain.NewRequestID("req2"),
			RequesterID: domain.NewUserID("user3"),
			TargetID:    domain.NewUserID("user4"),
			Status:      domain.FriendRequestStatusAccepted,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	protoRequests := ToProtoFriendRequests(requests)

	if len(protoRequests) != 2 {
		t.Fatalf("Expected 2 requests, got %d", len(protoRequests))
	}

	if protoRequests[0].Status != "pending" {
		t.Errorf("Expected first request status 'pending', got '%s'", protoRequests[0].Status)
	}

	if protoRequests[1].Status != "accepted" {
		t.Errorf("Expected second request status 'accepted', got '%s'", protoRequests[1].Status)
	}
}

func TestToProtoFriendRequests_EmptySlice_ReturnsEmptySlice(t *testing.T) {
	requests := []*domain.FriendRequest{}

	protoRequests := ToProtoFriendRequests(requests)

	if len(protoRequests) != 0 {
		t.Errorf("Expected empty slice, got %d requests", len(protoRequests))
	}
}

