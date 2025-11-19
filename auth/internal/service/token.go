package service

import (
	"context"

	"github.com/go-chat/auth/internal/domain"
)

// TokenService handles JWT token operations, storage, and public key distribution
type TokenService interface {
	// GenerateTokenPair creates JWT access token and JWT refresh token
	// Returns the token pair and refresh token metadata for efficient storage
	GenerateTokenPair(ctx context.Context, userID domain.UserID, email string) (*domain.TokenPair, *domain.RefreshToken, error)

	// StoreRefreshToken stores the refresh token metadata in repository
	StoreRefreshToken(ctx context.Context, refreshToken *domain.RefreshToken) error

	// ValidateAndRevokeRefreshToken validates refresh token, checks database, and revokes it
	// Returns user ID if valid, error otherwise
	ValidateAndRevokeRefreshToken(ctx context.Context, refreshToken string) (domain.UserID, error)

	// GetPublicKeys returns public keys in JWK format for JWT validation
	GetPublicKeys(ctx context.Context) ([]*domain.PublicKey, error)
}
