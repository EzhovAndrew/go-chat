package domain

// NotificationID represents a unique notification identifier
type NotificationID string

// String returns the string representation of NotificationID
func (id NotificationID) String() string {
	return string(id)
}

// NewNotificationID creates a new NotificationID from a string
func NewNotificationID(id string) NotificationID {
	return NotificationID(id)
}

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

// NotificationType represents the type of notification
type NotificationType int

const (
	// NotificationTypeMessage indicates a new chat message notification
	NotificationTypeMessage NotificationType = iota + 1

	// NotificationTypeFriendRequest indicates an incoming friend request notification
	NotificationTypeFriendRequest

	// NotificationTypeFriendAccepted indicates a friend request acceptance notification
	NotificationTypeFriendAccepted
)

// String returns the string representation of NotificationType
func (t NotificationType) String() string {
	switch t {
	case NotificationTypeMessage:
		return "message"
	case NotificationTypeFriendRequest:
		return "friend_request"
	case NotificationTypeFriendAccepted:
		return "friend_accepted"
	default:
		return "unknown"
	}
}

