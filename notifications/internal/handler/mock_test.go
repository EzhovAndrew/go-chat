package handler

import (
	"context"
	"errors"

	"github.com/go-chat/notifications/internal/domain"
)

// mockNotificationService is a mock implementation of service.NotificationService
type mockNotificationService struct {
	getNotificationsFunc func(ctx context.Context, userID domain.UserID, cursor string, limit int32) ([]*domain.Notification, string, error)
	markAsReadFunc       func(ctx context.Context, notificationID domain.NotificationID, userID domain.UserID) error
}

func (m *mockNotificationService) GetNotifications(ctx context.Context, userID domain.UserID, cursor string, limit int32) ([]*domain.Notification, string, error) {
	if m.getNotificationsFunc != nil {
		return m.getNotificationsFunc(ctx, userID, cursor, limit)
	}
	return nil, "", errors.New("not implemented")
}

func (m *mockNotificationService) MarkAsRead(ctx context.Context, notificationID domain.NotificationID, userID domain.UserID) error {
	if m.markAsReadFunc != nil {
		return m.markAsReadFunc(ctx, notificationID, userID)
	}
	return errors.New("not implemented")
}

