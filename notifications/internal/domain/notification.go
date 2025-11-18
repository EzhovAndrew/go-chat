package domain

import "time"

// Notification represents a notification for a user
type Notification struct {
	NotificationID NotificationID
	UserID         UserID
	Type           NotificationType
	Title          string
	Body           string
	Data           string // JSON-encoded additional data
	CreatedAt      time.Time
	Read           bool
}

