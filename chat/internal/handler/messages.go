package handler

import (
	"context"
	"time"

	"github.com/go-chat/chat/internal/domain"
	chatv1 "github.com/go-chat/chat/pkg/api/chat/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// SendMessage sends a message to a chat
func (s *Server) SendMessage(ctx context.Context, req *chatv1.SendMessageRequest) (*chatv1.SendMessageResponse, error) {
	// TODO: Extract authenticated user_id (sender_id) from JWT
	senderID := domain.NewUserID("authenticated-user-id") // Placeholder

	// Delegate to service layer
	message, err := s.messageService.SendMessage(
		ctx,
		domain.NewChatID(req.ChatId),
		senderID,
		req.Text,
		req.IdempotencyKey,
	)
	if err != nil {
		return nil, err // Middleware will map domain error to gRPC status
	}

	// Convert domain model to proto message
	return &chatv1.SendMessageResponse{
		Message: &chatv1.Message{
			MessageId: message.MessageID.String(),
			ChatId:    message.ChatID.String(),
			SenderId:  message.SenderID.String(),
			Text:      message.Text,
			CreatedAt: timestamppb.New(message.CreatedAt),
		},
	}, nil
}

// ListMessages retrieves message history for a chat with cursor-based pagination
func (s *Server) ListMessages(ctx context.Context, req *chatv1.ListMessagesRequest) (*chatv1.ListMessagesResponse, error) {
	// TODO: Extract authenticated user_id from JWT
	requesterID := domain.NewUserID("authenticated-user-id") // Placeholder

	// Default limit if not provided
	limit := req.Limit
	if limit == 0 {
		limit = 50
	}

	// Delegate to service layer
	messages, nextCursor, err := s.messageService.ListMessages(
		ctx,
		domain.NewChatID(req.ChatId),
		requesterID,
		req.Cursor,
		limit,
	)
	if err != nil {
		return nil, err // Middleware will map domain error to gRPC status
	}

	// Convert domain models to proto messages
	pbMessages := make([]*chatv1.Message, len(messages))
	for i, message := range messages {
		pbMessages[i] = &chatv1.Message{
			MessageId: message.MessageID.String(),
			ChatId:    message.ChatID.String(),
			SenderId:  message.SenderID.String(),
			Text:      message.Text,
			CreatedAt: timestamppb.New(message.CreatedAt),
		}
	}

	return &chatv1.ListMessagesResponse{
		Messages:   pbMessages,
		NextCursor: nextCursor,
	}, nil
}

// StreamMessages streams new messages in real-time (server-side streaming)
func (s *Server) StreamMessages(req *chatv1.StreamMessagesRequest, stream chatv1.ChatService_StreamMessagesServer) error {
	// TODO: Extract authenticated user_id from JWT in metadata
	requesterID := domain.NewUserID("authenticated-user-id") // Placeholder

	// Convert proto timestamp to time.Time
	since := time.Time{}
	if req.Since != nil {
		since = req.Since.AsTime()
	}

	// Delegate to service layer (placeholder implementation)
	err := s.messageService.StreamMessages(
		stream.Context(),
		domain.NewChatID(req.ChatId),
		requesterID,
		since,
	)
	if err != nil {
		return err // Middleware will map domain error to gRPC status
	}

	// Note: Actual streaming implementation will be done in the service layer
	// The service will send messages through a channel or similar mechanism
	return nil
}
