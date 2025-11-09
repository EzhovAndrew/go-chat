package domain

import "time"

// UserProfile represents a user's profile information
type UserProfile struct {
	UserID    UserID
	Nickname  string
	Bio       string
	AvatarURL string
	CreatedAt time.Time
	UpdatedAt time.Time
}

