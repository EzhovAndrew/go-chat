package dto

import (
	"testing"
	"time"

	"github.com/go-chat/users/internal/domain"
)

func TestToProtoProfile_ValidProfile_ConvertsCorrectly(t *testing.T) {
	profile := &domain.UserProfile{
		UserID:    domain.NewUserID("550e8400-e29b-41d4-a716-446655440000"),
		Nickname:  "johndoe",
		Bio:       "Software Engineer",
		AvatarURL: "https://example.com/avatar.jpg",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	protoProfile := ToProtoProfile(profile)

	if protoProfile.UserId != "550e8400-e29b-41d4-a716-446655440000" {
		t.Errorf("Expected UserId '550e8400-e29b-41d4-a716-446655440000', got '%s'", protoProfile.UserId)
	}

	if protoProfile.Nickname != "johndoe" {
		t.Errorf("Expected Nickname 'johndoe', got '%s'", protoProfile.Nickname)
	}

	if protoProfile.Bio != "Software Engineer" {
		t.Errorf("Expected Bio 'Software Engineer', got '%s'", protoProfile.Bio)
	}

	if protoProfile.AvatarUrl != "https://example.com/avatar.jpg" {
		t.Errorf("Expected AvatarUrl 'https://example.com/avatar.jpg', got '%s'", protoProfile.AvatarUrl)
	}
}

func TestToProtoProfile_EmptyFields_HandlesCorrectly(t *testing.T) {
	profile := &domain.UserProfile{
		UserID:    domain.NewUserID("550e8400-e29b-41d4-a716-446655440000"),
		Nickname:  "johndoe",
		Bio:       "",
		AvatarURL: "",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	protoProfile := ToProtoProfile(profile)

	if protoProfile.Bio != "" {
		t.Errorf("Expected empty Bio, got '%s'", protoProfile.Bio)
	}

	if protoProfile.AvatarUrl != "" {
		t.Errorf("Expected empty AvatarUrl, got '%s'", protoProfile.AvatarUrl)
	}
}

func TestToProtoProfiles_MultipleProfiles_ConvertsAll(t *testing.T) {
	profiles := []*domain.UserProfile{
		{
			UserID:    domain.NewUserID("550e8400-e29b-41d4-a716-446655440001"),
			Nickname:  "alice",
			Bio:       "Designer",
			AvatarURL: "https://example.com/alice.jpg",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			UserID:    domain.NewUserID("550e8400-e29b-41d4-a716-446655440002"),
			Nickname:  "bob",
			Bio:       "Developer",
			AvatarURL: "https://example.com/bob.jpg",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	protoProfiles := ToProtoProfiles(profiles)

	if len(protoProfiles) != 2 {
		t.Fatalf("Expected 2 profiles, got %d", len(protoProfiles))
	}

	if protoProfiles[0].Nickname != "alice" {
		t.Errorf("Expected first profile nickname 'alice', got '%s'", protoProfiles[0].Nickname)
	}

	if protoProfiles[1].Nickname != "bob" {
		t.Errorf("Expected second profile nickname 'bob', got '%s'", protoProfiles[1].Nickname)
	}
}

func TestToProtoProfiles_EmptySlice_ReturnsEmptySlice(t *testing.T) {
	profiles := []*domain.UserProfile{}

	protoProfiles := ToProtoProfiles(profiles)

	if len(protoProfiles) != 0 {
		t.Errorf("Expected empty slice, got %d profiles", len(protoProfiles))
	}
}

