package handler

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/go-chat/social/internal/domain"
	socialv1 "github.com/go-chat/social/pkg/api/social/v1"
)

func TestSendFriendRequest_ValidRequest_ReturnsRequest(t *testing.T) {
	expectedRequest := &domain.FriendRequest{
		RequestID:   domain.NewRequestID("550e8400-e29b-41d4-a716-446655440000"),
		RequesterID: domain.NewUserID("550e8400-e29b-41d4-a716-446655440001"),
		TargetID:    domain.NewUserID("550e8400-e29b-41d4-a716-446655440002"),
		Status:      domain.FriendRequestStatusPending,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	mockFRSvc := &mockFriendRequestService{
		sendFriendRequestFunc: func(ctx context.Context, requesterID, targetID domain.UserID) (*domain.FriendRequest, error) {
			return expectedRequest, nil
		},
	}

	server := NewServer(mockFRSvc, nil, nil, nil)
	req := &socialv1.SendFriendRequestRequest{
		TargetUserId: "550e8400-e29b-41d4-a716-446655440002",
	}

	resp, err := server.SendFriendRequest(context.Background(), req)

	if err != nil {
		t.Fatalf("SendFriendRequest() returned error: %v", err)
	}

	if resp.Request.RequestId != expectedRequest.RequestID.String() {
		t.Errorf("Expected request ID '%s', got '%s'", expectedRequest.RequestID.String(), resp.Request.RequestId)
	}
}

func TestSendFriendRequest_AlreadyFriends_ReturnsError(t *testing.T) {
	mockFRSvc := &mockFriendRequestService{
		sendFriendRequestFunc: func(ctx context.Context, requesterID, targetID domain.UserID) (*domain.FriendRequest, error) {
			return nil, domain.ErrAlreadyFriends
		},
	}

	server := NewServer(mockFRSvc, nil, nil, nil)
	req := &socialv1.SendFriendRequestRequest{
		TargetUserId: "550e8400-e29b-41d4-a716-446655440002",
	}

	_, err := server.SendFriendRequest(context.Background(), req)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if !errors.Is(err, domain.ErrAlreadyFriends) {
		t.Errorf("Expected ErrAlreadyFriends, got: %v", err)
	}
}

func TestAcceptFriendRequest_ValidRequest_ReturnsAcceptedRequest(t *testing.T) {
	expectedRequest := &domain.FriendRequest{
		RequestID:   domain.NewRequestID("550e8400-e29b-41d4-a716-446655440000"),
		RequesterID: domain.NewUserID("550e8400-e29b-41d4-a716-446655440001"),
		TargetID:    domain.NewUserID("550e8400-e29b-41d4-a716-446655440002"),
		Status:      domain.FriendRequestStatusAccepted,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	mockFRSvc := &mockFriendRequestService{
		acceptFriendRequestFunc: func(ctx context.Context, requestID domain.RequestID, userID domain.UserID) (*domain.FriendRequest, error) {
			return expectedRequest, nil
		},
	}

	server := NewServer(mockFRSvc, nil, nil, nil)
	req := &socialv1.AcceptFriendRequestRequest{
		RequestId: "550e8400-e29b-41d4-a716-446655440000",
	}

	resp, err := server.AcceptFriendRequest(context.Background(), req)

	if err != nil {
		t.Fatalf("AcceptFriendRequest() returned error: %v", err)
	}

	if resp.Request.Status != string(domain.FriendRequestStatusAccepted) {
		t.Errorf("Expected status 'accepted', got '%s'", resp.Request.Status)
	}
}

func TestListRequests_ValidRequest_ReturnsRequests(t *testing.T) {
	expectedRequests := []*domain.FriendRequest{
		{
			RequestID:   domain.NewRequestID("550e8400-e29b-41d4-a716-446655440000"),
			RequesterID: domain.NewUserID("550e8400-e29b-41d4-a716-446655440001"),
			TargetID:    domain.NewUserID("550e8400-e29b-41d4-a716-446655440002"),
			Status:      domain.FriendRequestStatusPending,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	mockFRSvc := &mockFriendRequestService{
		listRequestsFunc: func(ctx context.Context, userID domain.UserID, cursor string, limit int32) ([]*domain.FriendRequest, string, error) {
			return expectedRequests, "next_cursor_123", nil
		},
	}

	server := NewServer(mockFRSvc, nil, nil, nil)
	req := &socialv1.ListRequestsRequest{
		UserId: "550e8400-e29b-41d4-a716-446655440002",
		Limit:  10,
	}

	resp, err := server.ListRequests(context.Background(), req)

	if err != nil {
		t.Fatalf("ListRequests() returned error: %v", err)
	}

	if len(resp.Requests) != 1 {
		t.Errorf("Expected 1 request, got %d", len(resp.Requests))
	}

	if resp.NextCursor != "next_cursor_123" {
		t.Errorf("Expected cursor 'next_cursor_123', got '%s'", resp.NextCursor)
	}
}

