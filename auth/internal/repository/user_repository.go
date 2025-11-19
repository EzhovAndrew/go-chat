package repository

import (
	"context"

	"github.com/go-chat/auth/internal/domain"
)

// UserRepository defines the interface for user data access operations
type UserRepository interface {
	// Create stores a new user
	// Returns domain.ErrEmailAlreadyExists if email already exists (unique constraint violation)
	Create(ctx context.Context, user *domain.User) error

	// GetByID retrieves a user by ID
	GetByID(ctx context.Context, userID domain.UserID) (*domain.User, error)

	// GetByEmail retrieves a user by email
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
}
