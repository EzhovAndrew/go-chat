package service

import (
	"context"

	"github.com/go-chat/social/internal/domain"
)

// BlockService handles user blocking operations
type BlockService interface {
	// BlockUser blocks a user
	BlockUser(ctx context.Context, blockerID, blockedID domain.UserID) error

	// UnblockUser unblocks a previously blocked user
	UnblockUser(ctx context.Context, blockerID, blockedID domain.UserID) error
}

