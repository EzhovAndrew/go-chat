package handler

import (
	"context"
	"log"
	"time"

	chatv1 "github.com/go-chat/chat/pkg/api/chat/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// SendMessage sends a message to a chat
// Returns dummy message until database and Kafka integration is added
func (s *Server) SendMessage(ctx context.Context, req *chatv1.SendMessageRequest) (*chatv1.SendMessageResponse, error) {
	log.Printf("SendMessage called for chat_id: %s, text length: %d", req.ChatId, len(req.Text))

	// TODO: Extract authenticated user_id (sender_id) from JWT
	// TODO: Query chat from database
	// TODO: Return NOT_FOUND if chat doesn't exist
	// TODO: Verify sender is a participant (PERMISSION_DENIED if not)
	// TODO: Check idempotency_key if provided (return existing message if duplicate)
	// TODO: Validate text length (INVALID_ARGUMENT if too long)
	// TODO: Store message in database
	// TODO: Publish message.sent event to Kafka for notifications
	// TODO: Return created message

	return &chatv1.SendMessageResponse{
		Message: &chatv1.Message{
			MessageId: "550e8400-e29b-41d4-a716-446655440000",
			ChatId:    req.ChatId,
			SenderId:  "550e8400-e29b-41d4-a716-446655440001",
			Text:      req.Text,
			CreatedAt: timestamppb.Now(),
		},
	}, nil
}

// ListMessages retrieves message history for a chat with cursor-based pagination
// Returns dummy messages until database integration is added
func (s *Server) ListMessages(ctx context.Context, req *chatv1.ListMessagesRequest) (*chatv1.ListMessagesResponse, error) {
	log.Printf("ListMessages called for chat_id: %s, limit: %d", req.ChatId, req.Limit)

	// TODO: Extract authenticated user_id from JWT
	// TODO: Query chat from database
	// TODO: Return NOT_FOUND if chat doesn't exist
	// TODO: Verify user is a participant (PERMISSION_DENIED if not)
	// TODO: Decode cursor to get pagination position
	// TODO: Query messages from database ordered by created_at DESC
	// TODO: Apply limit + 1 to check for more results
	// TODO: Encode next_cursor if more results exist

	return &chatv1.ListMessagesResponse{
		Messages: []*chatv1.Message{
			{
				MessageId: "550e8400-e29b-41d4-a716-446655440001",
				ChatId:    req.ChatId,
				SenderId:  "550e8400-e29b-41d4-a716-446655440001",
				Text:      "Hello, this is a dummy message!",
				CreatedAt: timestamppb.Now(),
			},
			{
				MessageId: "550e8400-e29b-41d4-a716-446655440002",
				ChatId:    req.ChatId,
				SenderId:  "550e8400-e29b-41d4-a716-446655440002",
				Text:      "This is another dummy message.",
				CreatedAt: timestamppb.Now(),
			},
		},
		NextCursor: "", // Empty means no more results
	}, nil
}

// StreamMessages streams new messages in real-time (server-side streaming)
// Sends dummy messages until database and real-time implementation is added
func (s *Server) StreamMessages(req *chatv1.StreamMessagesRequest, stream chatv1.ChatService_StreamMessagesServer) error {
	log.Printf("StreamMessages called for chat_id: %s", req.ChatId)

	// TODO: Extract authenticated user_id from JWT in metadata
	// TODO: Query chat from database
	// TODO: Verify user is a participant (PERMISSION_DENIED if not)
	// TODO: Subscribe to message events for this chat (via channel or pub/sub)
	// TODO: If since_unix_ms provided, send messages created after that timestamp
	// TODO: Stream new messages as they arrive
	// TODO: Handle context cancellation for cleanup

	// For now, send a single dummy message and close the stream
	dummyMessage := &chatv1.StreamMessagesResponse{
		Message: &chatv1.Message{
			MessageId: "550e8400-e29b-41d4-a716-446655440000",
			ChatId:    req.ChatId,
			SenderId:  "550e8400-e29b-41d4-a716-446655440001",
			Text:      "This is a dummy streamed message",
			CreatedAt: timestamppb.Now(),
		},
	}

	if err := stream.Send(dummyMessage); err != nil {
		log.Printf("Error sending message: %v", err)
		return err
	}

	// Keep stream open for 10 seconds to demonstrate streaming
	time.Sleep(10 * time.Second)

	return nil
}

