package handler

import (
	"context"
	"testing"
	"time"

	"github.com/go-chat/notifications/internal/domain"
	notificationsv1 "github.com/go-chat/notifications/pkg/api/notifications/v1"
)

func TestGetNotifications_ValidRequest_ReturnsNotifications(t *testing.T) {
	expectedNotifications := []*domain.Notification{
		{
			NotificationID: domain.NewNotificationID("550e8400-e29b-41d4-a716-446655440000"),
			UserID:         domain.NewUserID("550e8400-e29b-41d4-a716-446655440001"),
			Type:           domain.NotificationTypeMessage,
			Title:          "New Message",
			Body:           "You have a new message",
			Data:           `{"chat_id":"123"}`,
			CreatedAt:      time.Now(),
			Read:           false,
		},
		{
			NotificationID: domain.NewNotificationID("550e8400-e29b-41d4-a716-446655440002"),
			UserID:         domain.NewUserID("550e8400-e29b-41d4-a716-446655440001"),
			Type:           domain.NotificationTypeFriendRequest,
			Title:          "Friend Request",
			Body:           "You have a new friend request",
			Data:           `{"request_id":"456"}`,
			CreatedAt:      time.Now(),
			Read:           false,
		},
	}

	mockSvc := &mockNotificationService{
		getNotificationsFunc: func(ctx context.Context, userID domain.UserID, cursor string, limit int32) ([]*domain.Notification, string, error) {
			return expectedNotifications, "next_cursor_123", nil
		},
	}

	server := NewServer(mockSvc)
	req := &notificationsv1.GetNotificationsRequest{
		UserId: "550e8400-e29b-41d4-a716-446655440001",
		Limit:  10,
	}

	resp, err := server.GetNotifications(context.Background(), req)

	if err != nil {
		t.Fatalf("GetNotifications() returned error: %v", err)
	}

	if len(resp.Notifications) != 2 {
		t.Errorf("Expected 2 notifications, got %d", len(resp.Notifications))
	}

	if resp.NextCursor != "next_cursor_123" {
		t.Errorf("Expected cursor 'next_cursor_123', got '%s'", resp.NextCursor)
	}

	// Verify first notification type conversion
	if resp.Notifications[0].Type != notificationsv1.NotificationType_NOTIFICATION_TYPE_MESSAGE {
		t.Errorf("Expected type MESSAGE, got %v", resp.Notifications[0].Type)
	}

	// Verify second notification type conversion
	if resp.Notifications[1].Type != notificationsv1.NotificationType_NOTIFICATION_TYPE_FRIEND_REQUEST {
		t.Errorf("Expected type FRIEND_REQUEST, got %v", resp.Notifications[1].Type)
	}
}

func TestGetNotifications_WithDefaultLimit_UsesDefaultValue(t *testing.T) {
	mockSvc := &mockNotificationService{
		getNotificationsFunc: func(ctx context.Context, userID domain.UserID, cursor string, limit int32) ([]*domain.Notification, string, error) {
			if limit != 20 {
				t.Errorf("Expected default limit 20, got %d", limit)
			}
			return []*domain.Notification{}, "", nil
		},
	}

	server := NewServer(mockSvc)
	req := &notificationsv1.GetNotificationsRequest{
		UserId: "550e8400-e29b-41d4-a716-446655440001",
		Limit:  0, // Not provided, should default to 20
	}

	_, err := server.GetNotifications(context.Background(), req)

	if err != nil {
		t.Fatalf("GetNotifications() returned error: %v", err)
	}
}

func TestGetNotifications_NoResults_ReturnsEmptyList(t *testing.T) {
	mockSvc := &mockNotificationService{
		getNotificationsFunc: func(ctx context.Context, userID domain.UserID, cursor string, limit int32) ([]*domain.Notification, string, error) {
			return []*domain.Notification{}, "", nil
		},
	}

	server := NewServer(mockSvc)
	req := &notificationsv1.GetNotificationsRequest{
		UserId: "550e8400-e29b-41d4-a716-446655440001",
		Limit:  10,
	}

	resp, err := server.GetNotifications(context.Background(), req)

	if err != nil {
		t.Fatalf("GetNotifications() returned error: %v", err)
	}

	if len(resp.Notifications) != 0 {
		t.Errorf("Expected 0 notifications, got %d", len(resp.Notifications))
	}

	if resp.NextCursor != "" {
		t.Errorf("Expected empty cursor, got '%s'", resp.NextCursor)
	}
}

