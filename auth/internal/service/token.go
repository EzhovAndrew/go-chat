package service

import (
	"context"

	"github.com/go-chat/auth/internal/domain"
)

// TokenService handles JWT token operations and public key distribution
type TokenService interface {
	// GetPublicKeys returns public keys in JWK format for JWT validation
	GetPublicKeys(ctx context.Context) ([]*domain.PublicKey, error)
}

