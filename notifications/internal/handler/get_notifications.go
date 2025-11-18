package handler

import (
	"context"

	"github.com/go-chat/notifications/internal/domain"
	"github.com/go-chat/notifications/internal/dto"
	notificationsv1 "github.com/go-chat/notifications/pkg/api/notifications/v1"
)

// GetNotifications retrieves notification history for a user
func (s *Server) GetNotifications(ctx context.Context, req *notificationsv1.GetNotificationsRequest) (*notificationsv1.GetNotificationsResponse, error) {
	// Default limit if not provided
	limit := req.Limit
	if limit == 0 {
		limit = 20
	}

	// Delegate to service layer
	notifications, nextCursor, err := s.notificationService.GetNotifications(
		ctx,
		domain.NewUserID(req.UserId),
		req.Cursor,
		limit,
	)
	if err != nil {
		return nil, err // Middleware will map domain error to gRPC status
	}

	// Convert domain models to proto messages
	return &notificationsv1.GetNotificationsResponse{
		Notifications: dto.ToProtoNotifications(notifications),
		NextCursor:    nextCursor,
	}, nil
}
