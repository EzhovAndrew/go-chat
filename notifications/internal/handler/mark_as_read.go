package handler

import (
	"context"
	"log"

	notificationsv1 "github.com/go-chat/notifications/pkg/api/notifications/v1"
)

// MarkAsRead marks a notification as read
// Returns empty response until database integration is added
func (s *Server) MarkAsRead(ctx context.Context, req *notificationsv1.MarkAsReadRequest) (*notificationsv1.MarkAsReadResponse, error) {
	log.Printf("MarkAsRead called for notification_id: %s", req.NotificationId)

	// TODO: Extract authenticated user_id from JWT
	// TODO: Query notification from database
	// TODO: Return NOT_FOUND if notification doesn't exist
	// TODO: Verify notification belongs to authenticated user (PERMISSION_DENIED if not)
	// TODO: Update is_read to true
	// TODO: Return empty response

	return &notificationsv1.MarkAsReadResponse{}, nil
}

