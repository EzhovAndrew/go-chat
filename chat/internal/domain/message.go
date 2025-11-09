package domain

import "time"

// Message represents a chat message
type Message struct {
	MessageID MessageID
	ChatID    ChatID
	SenderID  UserID
	Text      string
	CreatedAt time.Time
}

