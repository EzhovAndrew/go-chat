package handler

import (
	"context"

	"github.com/go-chat/notifications/internal/domain"
	notificationsv1 "github.com/go-chat/notifications/pkg/api/notifications/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
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
	pbNotifications := make([]*notificationsv1.Notification, len(notifications))
	for i, notification := range notifications {
		// Convert domain NotificationType to proto enum
		var protoType notificationsv1.NotificationType
		switch notification.Type {
		case domain.NotificationTypeMessage:
			protoType = notificationsv1.NotificationType_NOTIFICATION_TYPE_MESSAGE
		case domain.NotificationTypeFriendRequest:
			protoType = notificationsv1.NotificationType_NOTIFICATION_TYPE_FRIEND_REQUEST
		case domain.NotificationTypeFriendAccepted:
			protoType = notificationsv1.NotificationType_NOTIFICATION_TYPE_FRIEND_ACCEPTED
		default:
			protoType = notificationsv1.NotificationType_NOTIFICATION_TYPE_UNSPECIFIED
		}

		pbNotifications[i] = &notificationsv1.Notification{
			NotificationId: notification.NotificationID.String(),
			UserId:         notification.UserID.String(),
			Type:           protoType,
			Title:          notification.Title,
			Body:           notification.Body,
			Data:           notification.Data,
			CreatedAt:      timestamppb.New(notification.CreatedAt),
			Read:           notification.Read,
		}
	}

	return &notificationsv1.GetNotificationsResponse{
		Notifications: pbNotifications,
		NextCursor:    nextCursor,
	}, nil
}
