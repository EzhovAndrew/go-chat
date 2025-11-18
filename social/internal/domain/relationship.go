package domain

// RelationshipStatus represents the relationship between two users
type RelationshipStatus int

const (
	// RelationshipStatusNone indicates no relationship exists between users
	RelationshipStatusNone RelationshipStatus = iota

	// RelationshipStatusPending indicates a friend request has been sent but not yet accepted
	RelationshipStatusPending

	// RelationshipStatusFriend indicates users are friends
	RelationshipStatusFriend

	// RelationshipStatusBlocked indicates one user has blocked the other
	RelationshipStatusBlocked
)

// String returns the string representation of RelationshipStatus
func (s RelationshipStatus) String() string {
	switch s {
	case RelationshipStatusNone:
		return "none"
	case RelationshipStatusPending:
		return "pending"
	case RelationshipStatusFriend:
		return "friend"
	case RelationshipStatusBlocked:
		return "blocked"
	default:
		return "unknown"
	}
}

