package service

import (
	"context"

	"github.com/go-chat/auth/internal/domain"
)

// AuthService handles user authentication and token management
type AuthService interface {
	// Register creates a new user account
	Register(ctx context.Context, email, password string) (*domain.User, error)

	// Login authenticates user and returns tokens with user ID
	Login(ctx context.Context, email, password string) (*domain.TokenPair, domain.UserID, error)

	// Refresh validates refresh token and returns new token pair with user ID
	Refresh(ctx context.Context, refreshToken string) (*domain.TokenPair, domain.UserID, error)
}
