package handler

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/go-chat/users/internal/domain"
	usersv1 "github.com/go-chat/users/pkg/api/users/v1"
)

func TestGetProfileByID_ValidRequest_ReturnsProfile(t *testing.T) {
	expectedProfile := &domain.UserProfile{
		UserID:    domain.NewUserID("550e8400-e29b-41d4-a716-446655440000"),
		Nickname:  "john_doe",
		Bio:       "Software engineer",
		AvatarURL: "https://example.com/avatar.jpg",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockService := &mockUserService{
		getProfileByIDFunc: func(ctx context.Context, userID domain.UserID) (*domain.UserProfile, error) {
			return expectedProfile, nil
		},
	}

	server := NewServer(mockService)
	req := &usersv1.GetProfileByIDRequest{
		UserId: "550e8400-e29b-41d4-a716-446655440000",
	}

	resp, err := server.GetProfileByID(context.Background(), req)

	if err != nil {
		t.Fatalf("GetProfileByID() returned error: %v", err)
	}

	if resp.Profile.UserId != expectedProfile.UserID.String() {
		t.Errorf("Expected user ID '%s', got '%s'", expectedProfile.UserID.String(), resp.Profile.UserId)
	}
}

func TestGetProfileByID_ProfileNotFound_ReturnsError(t *testing.T) {
	mockService := &mockUserService{
		getProfileByIDFunc: func(ctx context.Context, userID domain.UserID) (*domain.UserProfile, error) {
			return nil, domain.ErrProfileNotFound
		},
	}

	server := NewServer(mockService)
	req := &usersv1.GetProfileByIDRequest{
		UserId: "550e8400-e29b-41d4-a716-446655440000",
	}

	_, err := server.GetProfileByID(context.Background(), req)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if !errors.Is(err, domain.ErrProfileNotFound) {
		t.Errorf("Expected ErrProfileNotFound, got: %v", err)
	}
}

func TestGetProfilesByIDs_ValidRequest_ReturnsProfiles(t *testing.T) {
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
			Nickname:  "jane_smith",
			Bio:       "Designer",
			AvatarURL: "https://example.com/avatar2.jpg",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	mockService := &mockUserService{
		getProfilesByIDsFunc: func(ctx context.Context, userIDs []domain.UserID) ([]*domain.UserProfile, error) {
			if len(userIDs) != 2 {
				t.Errorf("Expected 2 user IDs, got %d", len(userIDs))
			}
			return expectedProfiles, nil
		},
	}

	server := NewServer(mockService)
	req := &usersv1.GetProfilesByIDsRequest{
		UserIds: []string{
			"550e8400-e29b-41d4-a716-446655440000",
			"550e8400-e29b-41d4-a716-446655440001",
		},
	}

	resp, err := server.GetProfilesByIDs(context.Background(), req)

	if err != nil {
		t.Fatalf("GetProfilesByIDs() returned error: %v", err)
	}

	if len(resp.Profiles) != 2 {
		t.Errorf("Expected 2 profiles, got %d", len(resp.Profiles))
	}
}

func TestGetProfileByNickname_ValidRequest_ReturnsProfile(t *testing.T) {
	expectedProfile := &domain.UserProfile{
		UserID:    domain.NewUserID("550e8400-e29b-41d4-a716-446655440000"),
		Nickname:  "john_doe",
		Bio:       "Software engineer",
		AvatarURL: "https://example.com/avatar.jpg",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockService := &mockUserService{
		getProfileByNicknameFunc: func(ctx context.Context, nickname string) (*domain.UserProfile, error) {
			if nickname != "john_doe" {
				t.Errorf("Expected nickname 'john_doe', got '%s'", nickname)
			}
			return expectedProfile, nil
		},
	}

	server := NewServer(mockService)
	req := &usersv1.GetProfileByNicknameRequest{
		Nickname: "john_doe",
	}

	resp, err := server.GetProfileByNickname(context.Background(), req)

	if err != nil {
		t.Fatalf("GetProfileByNickname() returned error: %v", err)
	}

	if resp.Profile.Nickname != expectedProfile.Nickname {
		t.Errorf("Expected nickname '%s', got '%s'", expectedProfile.Nickname, resp.Profile.Nickname)
	}
}

func TestGetProfileByNickname_ProfileNotFound_ReturnsError(t *testing.T) {
	mockService := &mockUserService{
		getProfileByNicknameFunc: func(ctx context.Context, nickname string) (*domain.UserProfile, error) {
			return nil, domain.ErrProfileNotFound
		},
	}

	server := NewServer(mockService)
	req := &usersv1.GetProfileByNicknameRequest{
		Nickname: "nonexistent",
	}

	_, err := server.GetProfileByNickname(context.Background(), req)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if !errors.Is(err, domain.ErrProfileNotFound) {
		t.Errorf("Expected ErrProfileNotFound, got: %v", err)
	}
}

