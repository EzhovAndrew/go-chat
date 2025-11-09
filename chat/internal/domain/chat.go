package domain

import "time"

// Chat represents a chat conversation
type Chat struct {
	ChatID         ChatID
	ParticipantIDs []UserID
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

