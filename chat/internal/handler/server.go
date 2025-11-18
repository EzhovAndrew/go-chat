package handler

import (
	"github.com/go-chat/chat/internal/service"
	chatv1 "github.com/go-chat/chat/pkg/api/chat/v1"
)

// Server implements the ChatService gRPC interface
type Server struct {
	chatv1.UnimplementedChatServiceServer
	chatService    service.ChatService
	messageService service.MessageService
}

// NewServer creates a new Chat service handler with injected dependencies
func NewServer(chatService service.ChatService, messageService service.MessageService) *Server {
	return &Server{
		chatService:    chatService,
		messageService: messageService,
	}
}
