package handler

import (
	"context"
	"log"

	chatv1 "github.com/go-chat/chat/pkg/api/chat/v1"
)

// CreateDirectChat creates a 1-on-1 chat with another user
// Returns dummy chat_id until database and Social Service integration is added
func (s *Server) CreateDirectChat(ctx context.Context, req *chatv1.CreateDirectChatRequest) (*chatv1.CreateDirectChatResponse, error) {
	log.Printf("CreateDirectChat called with participant_id: %s", req.ParticipantId)

	// TODO: Extract authenticated user_id from JWT in gRPC metadata
	// TODO: Call SocialService.CheckRelationship to verify users are friends
	// TODO: Return PERMISSION_DENIED if not friends or if blocked
	// TODO: Check if chat already exists between these users
	// TODO: Return ALREADY_EXISTS if chat exists
	// TODO: Create chat in database with both user_ids as participants
	// TODO: Return new chat_id

	return &chatv1.CreateDirectChatResponse{
		ChatId: "550e8400-e29b-41d4-a716-446655440000",
	}, nil
}
