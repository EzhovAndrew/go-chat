package domain

import "time"

// User represents a user account in the authentication system
type User struct {
	ID           UserID
	Email        string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

