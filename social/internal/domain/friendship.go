package domain

import "time"

// Friendship represents a friendship between two users
type Friendship struct {
	FriendshipID FriendshipID
	UserID1      UserID
	UserID2      UserID
	CreatedAt    time.Time
}

