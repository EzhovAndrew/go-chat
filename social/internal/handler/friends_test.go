package handler

import (
	"context"
	"errors"
	"testing"

	"github.com/go-chat/social/internal/domain"
	socialv1 "github.com/go-chat/social/pkg/api/social/v1"
)

func TestListFriends_ValidRequest_ReturnsFriends(t *testing.T) {
	expectedFriends := []domain.UserID{
		domain.NewUserID("550e8400-e29b-41d4-a716-446655440001"),
		domain.NewUserID("550e8400-e29b-41d4-a716-446655440002"),
	}

	mockFriendshipSvc := &mockFriendshipService{
		listFriendsFunc: func(ctx context.Context, userID domain.UserID, cursor string, limit int32) ([]domain.UserID, string, error) {
			return expectedFriends, "next_cursor_789", nil
		},
	}

	server := NewServer(nil, mockFriendshipSvc, nil, nil)
	req := &socialv1.ListFriendsRequest{
		UserId: "550e8400-e29b-41d4-a716-446655440000",
		Limit:  10,
	}

	resp, err := server.ListFriends(context.Background(), req)

	if err != nil {
		t.Fatalf("ListFriends() returned error: %v", err)
	}

	if len(resp.UserIds) != 2 {
		t.Errorf("Expected 2 friends, got %d", len(resp.UserIds))
	}

	if resp.NextCursor != "next_cursor_789" {
		t.Errorf("Expected cursor 'next_cursor_789', got '%s'", resp.NextCursor)
	}
}

func TestListFriends_WithDefaultLimit_UsesDefaultValue(t *testing.T) {
	mockFriendshipSvc := &mockFriendshipService{
		listFriendsFunc: func(ctx context.Context, userID domain.UserID, cursor string, limit int32) ([]domain.UserID, string, error) {
			if limit != 20 {
				t.Errorf("Expected default limit 20, got %d", limit)
			}
			return []domain.UserID{}, "", nil
		},
	}

	server := NewServer(nil, mockFriendshipSvc, nil, nil)
	req := &socialv1.ListFriendsRequest{
		UserId: "550e8400-e29b-41d4-a716-446655440000",
		Limit:  0, // Not provided, should default to 20
	}

	_, err := server.ListFriends(context.Background(), req)

	if err != nil {
		t.Fatalf("ListFriends() returned error: %v", err)
	}
}

func TestRemoveFriend_ValidRequest_ReturnsSuccess(t *testing.T) {
	mockFriendshipSvc := &mockFriendshipService{
		removeFriendFunc: func(ctx context.Context, userID, friendID domain.UserID) error {
			if friendID.String() != "550e8400-e29b-41d4-a716-446655440002" {
				t.Errorf("Expected friend ID '550e8400-e29b-41d4-a716-446655440002', got '%s'", friendID.String())
			}
			return nil
		},
	}

	server := NewServer(nil, mockFriendshipSvc, nil, nil)
	req := &socialv1.RemoveFriendRequest{
		FriendUserId: "550e8400-e29b-41d4-a716-446655440002",
	}

	resp, err := server.RemoveFriend(context.Background(), req)

	if err != nil {
		t.Fatalf("RemoveFriend() returned error: %v", err)
	}

	if resp == nil {
		t.Fatal("Expected response, got nil")
	}
}

func TestRemoveFriend_NotFriends_ReturnsError(t *testing.T) {
	mockFriendshipSvc := &mockFriendshipService{
		removeFriendFunc: func(ctx context.Context, userID, friendID domain.UserID) error {
			return domain.ErrNotFriends
		},
	}

	server := NewServer(nil, mockFriendshipSvc, nil, nil)
	req := &socialv1.RemoveFriendRequest{
		FriendUserId: "550e8400-e29b-41d4-a716-446655440002",
	}

	_, err := server.RemoveFriend(context.Background(), req)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if !errors.Is(err, domain.ErrNotFriends) {
		t.Errorf("Expected ErrNotFriends, got: %v", err)
	}
}

