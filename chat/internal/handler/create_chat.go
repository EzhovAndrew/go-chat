package handler

import (
	"context"

	"github.com/go-chat/chat/internal/domain"
	chatv1 "github.com/go-chat/chat/pkg/api/chat/v1"
)

// CreateDirectChat creates a 1-on-1 chat with another user
func (s *Server) CreateDirectChat(ctx context.Context, req *chatv1.CreateDirectChatRequest) (*chatv1.CreateDirectChatResponse, error) {
	// TODO: Extract authenticated user_id from JWT in gRPC metadata
	// For now, assume requester ID will be extracted from context by service layer
	requesterID := domain.NewUserID("authenticated-user-id") // Placeholder

	// Delegate to service layer
	chatID, err := s.chatService.CreateDirectChat(
		ctx,
		requesterID,
		domain.NewUserID(req.ParticipantId),
	)
	if err != nil {
		return nil, err // Middleware will map domain error to gRPC status
	}

	return &chatv1.CreateDirectChatResponse{
		ChatId: chatID.String(),
	}, nil
}
