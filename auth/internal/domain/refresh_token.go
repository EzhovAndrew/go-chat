package domain

import "time"

// RefreshToken represents a refresh token for JWT token rotation
type RefreshToken struct {
	ID        string    // UUID
	UserID    UserID    // User identifier
	Token     string    // Hashed refresh token
	ExpiresAt time.Time
	Revoked   bool
	CreatedAt time.Time
}

