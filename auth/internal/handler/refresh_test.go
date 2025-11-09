package handler

import (
	"context"
	"errors"
	"testing"

	"github.com/go-chat/auth/internal/domain"
	authv1 "github.com/go-chat/auth/pkg/api/auth/v1"
)

func TestRefresh_ValidToken_ReturnsNewTokensAndUserID(t *testing.T) {
	expectedTokens := &domain.TokenPair{
		AccessToken:  "new_access_token_jwt",
		RefreshToken: "new_refresh_token_jwt",
	}
	expectedUserID := domain.NewUserID("550e8400-e29b-41d4-a716-446655440000")

	mockAuth := &mockAuthService{
		refreshFunc: func(ctx context.Context, refreshToken string) (*domain.TokenPair, domain.UserID, error) {
			if refreshToken != "old_refresh_token" {
				t.Errorf("Expected refresh token 'old_refresh_token', got '%s'", refreshToken)
			}
			return expectedTokens, expectedUserID, nil
		},
	}

	server := NewServer(mockAuth, nil)
	req := &authv1.RefreshRequest{
		RefreshToken: "old_refresh_token",
	}

	resp, err := server.Refresh(context.Background(), req)

	if err != nil {
		t.Fatalf("Refresh() returned error: %v", err)
	}

	if resp.AccessToken != expectedTokens.AccessToken {
		t.Errorf("Expected access token '%s', got '%s'", expectedTokens.AccessToken, resp.AccessToken)
	}

	if resp.RefreshToken != expectedTokens.RefreshToken {
		t.Errorf("Expected refresh token '%s', got '%s'", expectedTokens.RefreshToken, resp.RefreshToken)
	}

	if resp.UserId != expectedUserID.String() {
		t.Errorf("Expected user ID '%s', got '%s'", expectedUserID.String(), resp.UserId)
	}
}

func TestRefresh_InvalidToken_ReturnsError(t *testing.T) {
	mockAuth := &mockAuthService{
		refreshFunc: func(ctx context.Context, refreshToken string) (*domain.TokenPair, domain.UserID, error) {
			return nil, "", domain.ErrInvalidToken
		},
	}

	server := NewServer(mockAuth, nil)
	req := &authv1.RefreshRequest{
		RefreshToken: "invalid_token",
	}

	_, err := server.Refresh(context.Background(), req)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if !errors.Is(err, domain.ErrInvalidToken) {
		t.Errorf("Expected ErrInvalidToken, got: %v", err)
	}
}

func TestRefresh_ExpiredToken_ReturnsError(t *testing.T) {
	mockAuth := &mockAuthService{
		refreshFunc: func(ctx context.Context, refreshToken string) (*domain.TokenPair, domain.UserID, error) {
			return nil, "", domain.ErrTokenExpired
		},
	}

	server := NewServer(mockAuth, nil)
	req := &authv1.RefreshRequest{
		RefreshToken: "expired_token",
	}

	_, err := server.Refresh(context.Background(), req)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if !errors.Is(err, domain.ErrTokenExpired) {
		t.Errorf("Expected ErrTokenExpired, got: %v", err)
	}
}

func TestRefresh_RevokedToken_ReturnsError(t *testing.T) {
	mockAuth := &mockAuthService{
		refreshFunc: func(ctx context.Context, refreshToken string) (*domain.TokenPair, domain.UserID, error) {
			return nil, "", domain.ErrTokenRevoked
		},
	}

	server := NewServer(mockAuth, nil)
	req := &authv1.RefreshRequest{
		RefreshToken: "revoked_token",
	}

	_, err := server.Refresh(context.Background(), req)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if !errors.Is(err, domain.ErrTokenRevoked) {
		t.Errorf("Expected ErrTokenRevoked, got: %v", err)
	}
}

func TestRefresh_ServiceError_ReturnsError(t *testing.T) {
	expectedErr := errors.New("database error")

	mockAuth := &mockAuthService{
		refreshFunc: func(ctx context.Context, refreshToken string) (*domain.TokenPair, domain.UserID, error) {
			return nil, "", expectedErr
		},
	}

	server := NewServer(mockAuth, nil)
	req := &authv1.RefreshRequest{
		RefreshToken: "some_token",
	}

	_, err := server.Refresh(context.Background(), req)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if !errors.Is(err, expectedErr) {
		t.Errorf("Expected error '%v', got: %v", expectedErr, err)
	}
}
