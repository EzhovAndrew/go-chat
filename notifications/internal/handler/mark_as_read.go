package handler

import (
	"context"

	"github.com/go-chat/notifications/internal/domain"
	notificationsv1 "github.com/go-chat/notifications/pkg/api/notifications/v1"
)

// MarkAsRead marks a notification as read
func (s *Server) MarkAsRead(ctx context.Context, req *notificationsv1.MarkAsReadRequest) (*notificationsv1.MarkAsReadResponse, error) {
	// TODO: Extract authenticated user_id from JWT
	userID := domain.NewUserID("authenticated-user-id") // Placeholder

	// Delegate to service layer
	err := s.notificationService.MarkAsRead(
		ctx,
		domain.NewNotificationID(req.NotificationId),
		userID,
	)
	if err != nil {
		return nil, err // Middleware will map domain error to gRPC status
	}

	// Return empty response on success (as per proto)
	return &notificationsv1.MarkAsReadResponse{}, nil
}
