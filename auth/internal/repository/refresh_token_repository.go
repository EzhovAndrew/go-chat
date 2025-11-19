package repository

import (
	"context"

	"github.com/go-chat/auth/internal/domain"
)

// RefreshTokenRepository defines the interface for refresh token data access operations
type RefreshTokenRepository interface {
	// Create stores a new refresh token
	// The token.Token field should contain the JWT's jti claim (UUID)
	Create(ctx context.Context, token *domain.RefreshToken) error

	// GetByToken retrieves a refresh token by its jti
	// The jti parameter is the UUID from the JWT's jti claim
	GetByToken(ctx context.Context, jti string) (*domain.RefreshToken, error)

	// Revoke marks a refresh token as revoked by its jti
	// The jti parameter is the UUID from the JWT's jti claim
	Revoke(ctx context.Context, jti string) error
}
