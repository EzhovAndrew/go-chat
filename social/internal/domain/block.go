package domain

import "time"

// Block represents a block relationship where one user blocks another
type Block struct {
	BlockID   BlockID
	BlockerID UserID  // User who initiated the block
	BlockedID UserID  // User who is blocked
	CreatedAt time.Time
}

