package service

import (
	"context"

	"github.com/go-chat/users/internal/domain"
)

// UserService handles user profile management
type UserService interface {
	// CreateProfile creates a new user profile
	CreateProfile(ctx context.Context, userID domain.UserID, nickname, bio, avatarURL string) (*domain.UserProfile, error)

	// UpdateProfile updates an existing user profile
	UpdateProfile(ctx context.Context, userID domain.UserID, nickname, bio, avatarURL string) (*domain.UserProfile, error)

	// GetProfileByID retrieves a profile by user ID
	GetProfileByID(ctx context.Context, userID domain.UserID) (*domain.UserProfile, error)

	// GetProfilesByIDs retrieves multiple profiles by user IDs (batch operation)
	GetProfilesByIDs(ctx context.Context, userIDs []domain.UserID) ([]*domain.UserProfile, error)

	// GetProfileByNickname retrieves a profile by nickname
	GetProfileByNickname(ctx context.Context, nickname string) (*domain.UserProfile, error)

	// SearchByNickname searches for profiles matching a query with pagination
	SearchByNickname(ctx context.Context, query, cursor string, limit int32) ([]*domain.UserProfile, string, error)
}
