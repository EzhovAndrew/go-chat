package handler

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/go-chat/users/internal/domain"
	usersv1 "github.com/go-chat/users/pkg/api/users/v1"
)

func TestUpdateProfile_ValidRequest_ReturnsUpdatedProfile(t *testing.T) {
	expectedProfile := &domain.UserProfile{
		UserID:    domain.NewUserID("550e8400-e29b-41d4-a716-446655440000"),
		Nickname:  "new_nickname",
		Bio:       "Updated bio",
		AvatarURL: "https://example.com/new_avatar.jpg",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockService := &mockUserService{
		updateProfileFunc: func(ctx context.Context, userID domain.UserID, nickname, bio, avatarURL string) (*domain.UserProfile, error) {
			if userID.String() != "550e8400-e29b-41d4-a716-446655440000" {
				t.Errorf("Expected userID '550e8400-e29b-41d4-a716-446655440000', got '%s'", userID.String())
			}
			if nickname != "new_nickname" {
				t.Errorf("Expected nickname 'new_nickname', got '%s'", nickname)
			}
			return expectedProfile, nil
		},
	}

	server := NewServer(mockService)
	req := &usersv1.UpdateProfileRequest{
		UserId:    "550e8400-e29b-41d4-a716-446655440000",
		Nickname:  "new_nickname",
		Bio:       "Updated bio",
		AvatarUrl: "https://example.com/new_avatar.jpg",
	}

	resp, err := server.UpdateProfile(context.Background(), req)

	if err != nil {
		t.Fatalf("UpdateProfile() returned error: %v", err)
	}

	if resp.Profile.Nickname != expectedProfile.Nickname {
		t.Errorf("Expected nickname '%s', got '%s'", expectedProfile.Nickname, resp.Profile.Nickname)
	}
}

func TestUpdateProfile_ProfileNotFound_ReturnsError(t *testing.T) {
	mockService := &mockUserService{
		updateProfileFunc: func(ctx context.Context, userID domain.UserID, nickname, bio, avatarURL string) (*domain.UserProfile, error) {
			return nil, domain.ErrProfileNotFound
		},
	}

	server := NewServer(mockService)
	req := &usersv1.UpdateProfileRequest{
		UserId:   "550e8400-e29b-41d4-a716-446655440000",
		Nickname: "new_nickname",
	}

	_, err := server.UpdateProfile(context.Background(), req)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if !errors.Is(err, domain.ErrProfileNotFound) {
		t.Errorf("Expected ErrProfileNotFound, got: %v", err)
	}
}

func TestUpdateProfile_NicknameAlreadyExists_ReturnsError(t *testing.T) {
	mockService := &mockUserService{
		updateProfileFunc: func(ctx context.Context, userID domain.UserID, nickname, bio, avatarURL string) (*domain.UserProfile, error) {
			return nil, domain.ErrNicknameAlreadyExists
		},
	}

	server := NewServer(mockService)
	req := &usersv1.UpdateProfileRequest{
		UserId:   "550e8400-e29b-41d4-a716-446655440000",
		Nickname: "taken_nickname",
	}

	_, err := server.UpdateProfile(context.Background(), req)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if !errors.Is(err, domain.ErrNicknameAlreadyExists) {
		t.Errorf("Expected ErrNicknameAlreadyExists, got: %v", err)
	}
}

