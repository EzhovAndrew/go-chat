package service

import (
	"context"

	"github.com/go-chat/notifications/internal/domain"
)

// NotificationService handles notification operations
type NotificationService interface {
	// GetNotifications retrieves notification history for a user with pagination
	GetNotifications(ctx context.Context, userID domain.UserID, cursor string, limit int32) ([]*domain.Notification, string, error)

	// MarkAsRead marks a notification as read
	MarkAsRead(ctx context.Context, notificationID domain.NotificationID, userID domain.UserID) error
}

