package handler

import (
	notificationsv1 "github.com/go-chat/notifications/pkg/api/notifications/v1"
)

// Server implements the NotificationService gRPC interface
type Server struct {
	notificationsv1.UnimplementedNotificationServiceServer
}

// NewServer creates a new Notification service handler
// Dependencies (repository, kafka consumer) will be added as needed
func NewServer() *Server {
	return &Server{}
}

