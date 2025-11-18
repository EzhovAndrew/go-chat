package handler

import (
	"github.com/go-chat/notifications/internal/service"
	notificationsv1 "github.com/go-chat/notifications/pkg/api/notifications/v1"
)

// Server implements the NotificationService gRPC interface
type Server struct {
	notificationsv1.UnimplementedNotificationServiceServer
	notificationService service.NotificationService
}

// NewServer creates a new Notification service handler with injected dependencies
func NewServer(notificationService service.NotificationService) *Server {
	return &Server{
		notificationService: notificationService,
	}
}
