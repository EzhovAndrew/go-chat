package handler

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/go-chat/users/internal/domain"
	usersv1 "github.com/go-chat/users/pkg/api/users/v1"
)

func TestCreateProfile_ValidRequest_ReturnsProfile(t *testing.T) {
	expectedProfile := &domain.UserProfile{
		UserID:    domain.NewUserID("550e8400-e29b-41d4-a716-446655440000"),
		Nickname:  "john_doe",
		Bio:       "Software engineer",
		AvatarURL: "https://example.com/avatar.jpg",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockService := &mockUserService{
		createProfileFunc: func(ctx context.Context, userID domain.UserID, nickname, bio, avatarURL string) (*domain.UserProfile, error) {
			if userID.String() != "550e8400-e29b-41d4-a716-446655440000" {
				t.Errorf("Expected userID '550e8400-e29b-41d4-a716-446655440000', got '%s'", userID.String())
			}
			if nickname != "john_doe" {
				t.Errorf("Expected nickname 'john_doe', got '%s'", nickname)
			}
			return expectedProfile, nil
		},
	}

	server := NewServer(mockService)
	avatarURL := "https://example.com/avatar.jpg"
	req := &usersv1.CreateProfileRequest{
		UserId:    "550e8400-e29b-41d4-a716-446655440000",
		Nickname:  "john_doe",
		Bio:       "Software engineer",
		AvatarUrl: &avatarURL,
	}

	resp, err := server.CreateProfile(context.Background(), req)

	if err != nil {
		t.Fatalf("CreateProfile() returned error: %v", err)
	}

	if resp.Profile.UserId != expectedProfile.UserID.String() {
		t.Errorf("Expected user ID '%s', got '%s'", expectedProfile.UserID.String(), resp.Profile.UserId)
	}
	if resp.Profile.Nickname != expectedProfile.Nickname {
		t.Errorf("Expected nickname '%s', got '%s'", expectedProfile.Nickname, resp.Profile.Nickname)
	}
}

func TestCreateProfile_ProfileAlreadyExists_ReturnsError(t *testing.T) {
	mockService := &mockUserService{
		createProfileFunc: func(ctx context.Context, userID domain.UserID, nickname, bio, avatarURL string) (*domain.UserProfile, error) {
			return nil, domain.ErrProfileAlreadyExists
		},
	}

	server := NewServer(mockService)
	req := &usersv1.CreateProfileRequest{
		UserId:   "550e8400-e29b-41d4-a716-446655440000",
		Nickname: "john_doe",
	}

	_, err := server.CreateProfile(context.Background(), req)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if !errors.Is(err, domain.ErrProfileAlreadyExists) {
		t.Errorf("Expected ErrProfileAlreadyExists, got: %v", err)
	}
}

func TestCreateProfile_NicknameAlreadyExists_ReturnsError(t *testing.T) {
	mockService := &mockUserService{
		createProfileFunc: func(ctx context.Context, userID domain.UserID, nickname, bio, avatarURL string) (*domain.UserProfile, error) {
			return nil, domain.ErrNicknameAlreadyExists
		},
	}

	server := NewServer(mockService)
	req := &usersv1.CreateProfileRequest{
		UserId:   "550e8400-e29b-41d4-a716-446655440000",
		Nickname: "taken_nickname",
	}

	_, err := server.CreateProfile(context.Background(), req)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if !errors.Is(err, domain.ErrNicknameAlreadyExists) {
		t.Errorf("Expected ErrNicknameAlreadyExists, got: %v", err)
	}
}

func TestCreateProfile_WithoutOptionalAvatarURL_HandlesCorrectly(t *testing.T) {
	mockService := &mockUserService{
		createProfileFunc: func(ctx context.Context, userID domain.UserID, nickname, bio, avatarURL string) (*domain.UserProfile, error) {
			if avatarURL != "" {
				t.Errorf("Expected empty avatarURL, got '%s'", avatarURL)
			}
			return &domain.UserProfile{
				UserID:    userID,
				Nickname:  nickname,
				Bio:       bio,
				AvatarURL: "",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}, nil
		},
	}

	server := NewServer(mockService)
	req := &usersv1.CreateProfileRequest{
		UserId:   "550e8400-e29b-41d4-a716-446655440000",
		Nickname: "john_doe",
		Bio:      "Software engineer",
		// AvatarUrl is not provided (nil)
	}

	resp, err := server.CreateProfile(context.Background(), req)

	if err != nil {
		t.Fatalf("CreateProfile() returned error: %v", err)
	}

	if resp.Profile.AvatarUrl != "" {
		t.Errorf("Expected empty avatar URL, got '%s'", resp.Profile.AvatarUrl)
	}
}

