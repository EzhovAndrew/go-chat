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
