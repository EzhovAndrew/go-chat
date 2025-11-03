package handler

import (
	chatv1 "github.com/go-chat/chat/pkg/api/chat/v1"
)

// Server implements the ChatService gRPC interface
type Server struct {
	chatv1.UnimplementedChatServiceServer
}

// NewServer creates a new Chat service handler
// Dependencies (repository, social service client, kafka producer) will be added as needed
func NewServer() *Server {
	return &Server{}
}

