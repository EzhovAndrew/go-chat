package dto

import (
	"testing"
	"time"

	"github.com/go-chat/notifications/internal/domain"
)

func TestToProtoNotification_ValidNotification_ConvertsCorrectly(t *testing.T) {
	notification := &domain.Notification{
		NotificationID: domain.NewNotificationID("550e8400-e29b-41d4-a716-446655440000"),
		UserID:         domain.NewUserID("550e8400-e29b-41d4-a716-446655440001"),
		Type:           domain.NotificationTypeMessage,
		Title:          "New Message",
		Body:           "You have a new message from John",
		Data:           `{"chat_id":"123"}`,
		CreatedAt:      time.Now(),
		Read:           false,
	}

	protoNotification := ToProtoNotification(notification)

	if protoNotification.NotificationId != "550e8400-e29b-41d4-a716-446655440000" {
		t.Errorf("Expected NotificationId '550e8400-e29b-41d4-a716-446655440000', got '%s'", protoNotification.NotificationId)
	}

	if protoNotification.UserId != "550e8400-e29b-41d4-a716-446655440001" {
		t.Errorf("Expected UserId '550e8400-e29b-41d4-a716-446655440001', got '%s'", protoNotification.UserId)
	}

	if protoNotification.Title != "New Message" {
		t.Errorf("Expected Title 'New Message', got '%s'", protoNotification.Title)
	}

	if protoNotification.Body != "You have a new message from John" {
		t.Errorf("Expected Body 'You have a new message from John', got '%s'", protoNotification.Body)
	}

	if protoNotification.Read != false {
		t.Errorf("Expected Read false, got true")
	}
}

func TestToProtoNotification_AllTypes_MapsCorrectly(t *testing.T) {
	tests := []struct {
		name         string
		domainType   domain.NotificationType
		expectedEnum int32
	}{
		{"Message type", domain.NotificationTypeMessage, 1},
		{"Friend request type", domain.NotificationTypeFriendRequest, 2},
		{"Friend accepted type", domain.NotificationTypeFriendAccepted, 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			notification := &domain.Notification{
				NotificationID: domain.NewNotificationID("550e8400-e29b-41d4-a716-446655440000"),
				UserID:         domain.NewUserID("550e8400-e29b-41d4-a716-446655440001"),
				Type:           tt.domainType,
				Title:          "Test",
				Body:           "Test body",
				Data:           "",
				CreatedAt:      time.Now(),
				Read:           false,
			}

			protoNotification := ToProtoNotification(notification)

			if int32(protoNotification.Type) != tt.expectedEnum {
				t.Errorf("Expected type enum %d, got %d", tt.expectedEnum, protoNotification.Type)
			}
		})
	}
}

func TestToProtoNotifications_MultipleNotifications_ConvertsAll(t *testing.T) {
	notifications := []*domain.Notification{
		{
			NotificationID: domain.NewNotificationID("550e8400-e29b-41d4-a716-446655440001"),
			UserID:         domain.NewUserID("550e8400-e29b-41d4-a716-446655440000"),
			Type:           domain.NotificationTypeMessage,
			Title:          "Message 1",
			Body:           "Body 1",
			Data:           "",
			CreatedAt:      time.Now(),
			Read:           false,
		},
		{
			NotificationID: domain.NewNotificationID("550e8400-e29b-41d4-a716-446655440002"),
			UserID:         domain.NewUserID("550e8400-e29b-41d4-a716-446655440000"),
			Type:           domain.NotificationTypeFriendRequest,
			Title:          "Friend Request",
			Body:           "Body 2",
			Data:           "",
			CreatedAt:      time.Now(),
			Read:           true,
		},
	}

	protoNotifications := ToProtoNotifications(notifications)

	if len(protoNotifications) != 2 {
		t.Fatalf("Expected 2 notifications, got %d", len(protoNotifications))
	}

	if protoNotifications[0].Title != "Message 1" {
		t.Errorf("Expected first notification title 'Message 1', got '%s'", protoNotifications[0].Title)
	}

	if protoNotifications[1].Title != "Friend Request" {
		t.Errorf("Expected second notification title 'Friend Request', got '%s'", protoNotifications[1].Title)
	}

	if protoNotifications[1].Read != true {
		t.Errorf("Expected second notification to be read")
	}
}

func TestToProtoNotifications_EmptySlice_ReturnsEmptySlice(t *testing.T) {
	notifications := []*domain.Notification{}

	protoNotifications := ToProtoNotifications(notifications)

	if len(protoNotifications) != 0 {
		t.Errorf("Expected empty slice, got %d notifications", len(protoNotifications))
	}
}

