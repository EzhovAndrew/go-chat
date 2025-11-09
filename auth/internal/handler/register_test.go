package handler

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/go-chat/auth/internal/domain"
	authv1 "github.com/go-chat/auth/pkg/api/auth/v1"
)

// mockAuthService is a mock implementation of service.AuthService
type mockAuthService struct {
	registerFunc func(ctx context.Context, email, password string) (*domain.User, error)
	loginFunc    func(ctx context.Context, email, password string) (*domain.TokenPair, domain.UserID, error)
	refreshFunc  func(ctx context.Context, refreshToken string) (*domain.TokenPair, domain.UserID, error)
}

func (m *mockAuthService) Register(ctx context.Context, email, password string) (*domain.User, error) {
	if m.registerFunc != nil {
		return m.registerFunc(ctx, email, password)
	}
	return nil, errors.New("not implemented")
}

func (m *mockAuthService) Login(ctx context.Context, email, password string) (*domain.TokenPair, domain.UserID, error) {
	if m.loginFunc != nil {
		return m.loginFunc(ctx, email, password)
	}
	return nil, "", errors.New("not implemented")
}

func (m *mockAuthService) Refresh(ctx context.Context, refreshToken string) (*domain.TokenPair, domain.UserID, error) {
	if m.refreshFunc != nil {
		return m.refreshFunc(ctx, refreshToken)
	}
	return nil, "", errors.New("not implemented")
}

func TestRegister_ValidRequest_ReturnsUserID(t *testing.T) {
	expectedUser := &domain.User{
		ID:           domain.NewUserID("550e8400-e29b-41d4-a716-446655440000"),
		Email:        "test@example.com",
		PasswordHash: "hashed_password",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	mockAuth := &mockAuthService{
		registerFunc: func(ctx context.Context, email, password string) (*domain.User, error) {
			if email != "test@example.com" {
				t.Errorf("Expected email 'test@example.com', got '%s'", email)
			}
			if password != "SecurePass123!" {
				t.Errorf("Expected password 'SecurePass123!', got '%s'", password)
			}
			return expectedUser, nil
		},
	}

	server := NewServer(mockAuth, nil)
	req := &authv1.RegisterRequest{
		Email:    "test@example.com",
		Password: "SecurePass123!",
	}

	resp, err := server.Register(context.Background(), req)

	if err != nil {
		t.Fatalf("Register() returned error: %v", err)
	}

	if resp.UserId != expectedUser.ID.String() {
		t.Errorf("Expected user ID '%s', got '%s'", expectedUser.ID.String(), resp.UserId)
	}
}

func TestRegister_EmailAlreadyExists_ReturnsError(t *testing.T) {
	mockAuth := &mockAuthService{
		registerFunc: func(ctx context.Context, email, password string) (*domain.User, error) {
			return nil, domain.ErrEmailAlreadyExists
		},
	}

	server := NewServer(mockAuth, nil)
	req := &authv1.RegisterRequest{
		Email:    "existing@example.com",
		Password: "SecurePass123!",
	}

	_, err := server.Register(context.Background(), req)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if !errors.Is(err, domain.ErrEmailAlreadyExists) {
		t.Errorf("Expected ErrEmailAlreadyExists, got: %v", err)
	}
}

func TestRegister_ServiceError_ReturnsError(t *testing.T) {
	expectedErr := errors.New("database connection failed")

	mockAuth := &mockAuthService{
		registerFunc: func(ctx context.Context, email, password string) (*domain.User, error) {
			return nil, expectedErr
		},
	}

	server := NewServer(mockAuth, nil)
	req := &authv1.RegisterRequest{
		Email:    "test@example.com",
		Password: "SecurePass123!",
	}

	_, err := server.Register(context.Background(), req)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if !errors.Is(err, expectedErr) {
		t.Errorf("Expected error '%v', got: %v", expectedErr, err)
	}
}

