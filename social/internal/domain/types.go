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

// RequestID represents a unique friend request identifier
type RequestID string

// String returns the string representation of RequestID
func (id RequestID) String() string {
	return string(id)
}

// NewRequestID creates a new RequestID from a string
func NewRequestID(id string) RequestID {
	return RequestID(id)
}

// FriendshipID represents a unique friendship identifier
type FriendshipID string

// String returns the string representation of FriendshipID
func (id FriendshipID) String() string {
	return string(id)
}

// NewFriendshipID creates a new FriendshipID from a string
func NewFriendshipID(id string) FriendshipID {
	return FriendshipID(id)
}

// BlockID represents a unique block identifier
type BlockID string

// String returns the string representation of BlockID
func (id BlockID) String() string {
	return string(id)
}

// NewBlockID creates a new BlockID from a string
func NewBlockID(id string) BlockID {
	return BlockID(id)
}

