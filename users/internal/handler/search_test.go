package handler

import (
	"context"
	"testing"
	"time"

	"github.com/go-chat/users/internal/domain"
	usersv1 "github.com/go-chat/users/pkg/api/users/v1"
)

func TestSearchByNickname_ValidRequest_ReturnsResults(t *testing.T) {
	expectedProfiles := []*domain.UserProfile{
		{
			UserID:    domain.NewUserID("550e8400-e29b-41d4-a716-446655440000"),
			Nickname:  "john_doe",
			Bio:       "Software engineer",
			AvatarURL: "https://example.com/avatar1.jpg",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			UserID:    domain.NewUserID("550e8400-e29b-41d4-a716-446655440001"),
			Nickname:  "john_smith",
			Bio:       "Designer",
			AvatarURL: "https://example.com/avatar2.jpg",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	mockService := &mockUserService{
		searchByNicknameFunc: func(ctx context.Context, query, cursor string, limit int32) ([]*domain.UserProfile, string, error) {
			if query != "john" {
				t.Errorf("Expected query 'john', got '%s'", query)
			}
			return expectedProfiles, "next_cursor_123", nil
		},
	}

	server := NewServer(mockService)
	req := &usersv1.SearchByNicknameRequest{
		Query: "john",
		Limit: 10,
	}

	resp, err := server.SearchByNickname(context.Background(), req)

	if err != nil {
		t.Fatalf("SearchByNickname() returned error: %v", err)
	}

	if len(resp.Profiles) != 2 {
		t.Errorf("Expected 2 profiles, got %d", len(resp.Profiles))
	}

	if resp.NextCursor != "next_cursor_123" {
		t.Errorf("Expected cursor 'next_cursor_123', got '%s'", resp.NextCursor)
	}
}

func TestSearchByNickname_WithDefaultLimit_UsesDefaultValue(t *testing.T) {
	mockService := &mockUserService{
		searchByNicknameFunc: func(ctx context.Context, query, cursor string, limit int32) ([]*domain.UserProfile, string, error) {
			if limit != 20 {
				t.Errorf("Expected default limit 20, got %d", limit)
			}
			return []*domain.UserProfile{}, "", nil
		},
	}

	server := NewServer(mockService)
	req := &usersv1.SearchByNicknameRequest{
		Query: "test",
		Limit: 0, // Not provided, should default to 20
	}

	_, err := server.SearchByNickname(context.Background(), req)

	if err != nil {
		t.Fatalf("SearchByNickname() returned error: %v", err)
	}
}

func TestSearchByNickname_NoResults_ReturnsEmptyList(t *testing.T) {
	mockService := &mockUserService{
		searchByNicknameFunc: func(ctx context.Context, query, cursor string, limit int32) ([]*domain.UserProfile, string, error) {
			return []*domain.UserProfile{}, "", nil
		},
	}

	server := NewServer(mockService)
	req := &usersv1.SearchByNicknameRequest{
		Query: "nonexistent",
		Limit: 10,
	}

	resp, err := server.SearchByNickname(context.Background(), req)

	if err != nil {
		t.Fatalf("SearchByNickname() returned error: %v", err)
	}

	if len(resp.Profiles) != 0 {
		t.Errorf("Expected 0 profiles, got %d", len(resp.Profiles))
	}

	if resp.NextCursor != "" {
		t.Errorf("Expected empty cursor, got '%s'", resp.NextCursor)
	}
}
