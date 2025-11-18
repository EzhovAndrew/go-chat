package domain

// UserID represents a unique user identifier
type UserID string

// String returns the string representation of UserID
func (id UserID) String() string {
	return string(id)
}

// NewUserID creates a new UserID from a string
func NewUserID(id string) UserID {
	return UserID(id)
}

// ChatID represents a unique chat identifier
type ChatID string

// String returns the string representation of ChatID
func (id ChatID) String() string {
	return string(id)
}

// NewChatID creates a new ChatID from a string
func NewChatID(id string) ChatID {
	return ChatID(id)
}

// MessageID represents a unique message identifier
type MessageID string

// String returns the string representation of MessageID
func (id MessageID) String() string {
	return string(id)
}

// NewMessageID creates a new MessageID from a string
func NewMessageID(id string) MessageID {
	return MessageID(id)
}

