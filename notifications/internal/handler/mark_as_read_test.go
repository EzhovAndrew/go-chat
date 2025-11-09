package handler

import (
	"context"
	"errors"
	"testing"

	"github.com/go-chat/notifications/internal/domain"
	notificationsv1 "github.com/go-chat/notifications/pkg/api/notifications/v1"
)

func TestMarkAsRead_ValidRequest_ReturnsSuccess(t *testing.T) {
	mockSvc := &mockNotificationService{
		markAsReadFunc: func(ctx context.Context, notificationID domain.NotificationID, userID domain.UserID) error {
			if notificationID.String() != "550e8400-e29b-41d4-a716-446655440000" {
				t.Errorf("Expected notification ID '550e8400-e29b-41d4-a716-446655440000', got '%s'", notificationID.String())
			}
			return nil
		},
	}

	server := NewServer(mockSvc)
	req := &notificationsv1.MarkAsReadRequest{
		NotificationId: "550e8400-e29b-41d4-a716-446655440000",
	}

	resp, err := server.MarkAsRead(context.Background(), req)

	if err != nil {
		t.Fatalf("MarkAsRead() returned error: %v", err)
	}

	if resp == nil {
		t.Fatal("Expected response, got nil")
	}
}

func TestMarkAsRead_NotificationNotFound_ReturnsError(t *testing.T) {
	mockSvc := &mockNotificationService{
		markAsReadFunc: func(ctx context.Context, notificationID domain.NotificationID, userID domain.UserID) error {
			return domain.ErrNotificationNotFound
		},
	}

	server := NewServer(mockSvc)
	req := &notificationsv1.MarkAsReadRequest{
		NotificationId: "550e8400-e29b-41d4-a716-446655440000",
	}

	_, err := server.MarkAsRead(context.Background(), req)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if !errors.Is(err, domain.ErrNotificationNotFound) {
		t.Errorf("Expected ErrNotificationNotFound, got: %v", err)
	}
}

func TestMarkAsRead_PermissionDenied_ReturnsError(t *testing.T) {
	mockSvc := &mockNotificationService{
		markAsReadFunc: func(ctx context.Context, notificationID domain.NotificationID, userID domain.UserID) error {
			return domain.ErrPermissionDenied
		},
	}

	server := NewServer(mockSvc)
	req := &notificationsv1.MarkAsReadRequest{
		NotificationId: "550e8400-e29b-41d4-a716-446655440000",
	}

	_, err := server.MarkAsRead(context.Background(), req)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if !errors.Is(err, domain.ErrPermissionDenied) {
		t.Errorf("Expected ErrPermissionDenied, got: %v", err)
	}
}

