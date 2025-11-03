package handler

import (
	"context"
	"log"

	notificationsv1 "github.com/go-chat/notifications/pkg/api/notifications/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// GetNotifications retrieves notification history for a user with cursor-based pagination
// Returns dummy notifications until database integration is added
func (s *Server) GetNotifications(ctx context.Context, req *notificationsv1.GetNotificationsRequest) (*notificationsv1.GetNotificationsResponse, error) {
	log.Printf("GetNotifications called for user_id: %s, limit: %d", req.UserId, req.Limit)

	// TODO: Extract authenticated user_id from JWT
	// TODO: Verify req.UserId matches authenticated user
	// TODO: Decode cursor for pagination
	// TODO: Query notifications from database where user_id = req.UserId
	// TODO: Order by created_at DESC (newest first)
	// TODO: Apply limit + 1 to check for more results
	// TODO: Encode next_cursor if more results exist

	return &notificationsv1.GetNotificationsResponse{
		Notifications: []*notificationsv1.Notification{
			{
				NotificationId: "550e8400-e29b-41d4-a716-446655440001",
				UserId:         req.UserId,
				Type:           notificationsv1.NotificationType_NOTIFICATION_TYPE_MESSAGE,
				Title:          "New message",
				Body:           "You have a new message from user_123",
				Data:           `{"chat_id": "550e8400-e29b-41d4-a716-446655440000"}`,
				Read:           false,
				CreatedAt:      timestamppb.Now(),
			},
			{
				NotificationId: "550e8400-e29b-41d4-a716-446655440002",
				UserId:         req.UserId,
				Type:           notificationsv1.NotificationType_NOTIFICATION_TYPE_FRIEND_REQUEST,
				Title:          "New friend request",
				Body:           "user_456 sent you a friend request",
				Data:           `{"requester_id": "550e8400-e29b-41d4-a716-446655440003"}`,
				Read:           false,
				CreatedAt:      timestamppb.Now(),
			},
		},
		NextCursor: "", // Empty means no more results
	}, nil
}

