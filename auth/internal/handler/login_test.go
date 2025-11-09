package handler

import (
	"context"
	"errors"
	"testing"

	"github.com/go-chat/auth/internal/domain"
	authv1 "github.com/go-chat/auth/pkg/api/auth/v1"
)

func TestLogin_ValidCredentials_ReturnsTokensAndUserID(t *testing.T) {
	expectedTokens := &domain.TokenPair{
		AccessToken:  "access_token_jwt",
		RefreshToken: "refresh_token_jwt",
	}
	expectedUserID := domain.NewUserID("550e8400-e29b-41d4-a716-446655440000")

	mockAuth := &mockAuthService{
		loginFunc: func(ctx context.Context, email, password string) (*domain.TokenPair, domain.UserID, error) {
			if email != "test@example.com" {
				t.Errorf("Expected email 'test@example.com', got '%s'", email)
			}
			if password != "SecurePass123!" {
				t.Errorf("Expected password 'SecurePass123!', got '%s'", password)
			}
			return expectedTokens, expectedUserID, nil
		},
	}

	server := NewServer(mockAuth, nil)
	req := &authv1.LoginRequest{
		Email:    "test@example.com",
		Password: "SecurePass123!",
	}

	resp, err := server.Login(context.Background(), req)

	if err != nil {
		t.Fatalf("Login() returned error: %v", err)
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

func TestLogin_InvalidCredentials_ReturnsError(t *testing.T) {
	mockAuth := &mockAuthService{
		loginFunc: func(ctx context.Context, email, password string) (*domain.TokenPair, domain.UserID, error) {
			return nil, "", domain.ErrInvalidCredentials
		},
	}

	server := NewServer(mockAuth, nil)
	req := &authv1.LoginRequest{
		Email:    "test@example.com",
		Password: "WrongPassword",
	}

	_, err := server.Login(context.Background(), req)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if !errors.Is(err, domain.ErrInvalidCredentials) {
		t.Errorf("Expected ErrInvalidCredentials, got: %v", err)
	}
}

func TestLogin_ServiceError_ReturnsError(t *testing.T) {
	expectedErr := errors.New("token generation failed")

	mockAuth := &mockAuthService{
		loginFunc: func(ctx context.Context, email, password string) (*domain.TokenPair, domain.UserID, error) {
			return nil, "", expectedErr
		},
	}

	server := NewServer(mockAuth, nil)
	req := &authv1.LoginRequest{
		Email:    "test@example.com",
		Password: "SecurePass123!",
	}

	_, err := server.Login(context.Background(), req)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if !errors.Is(err, expectedErr) {
		t.Errorf("Expected error '%v', got: %v", expectedErr, err)
	}
}
