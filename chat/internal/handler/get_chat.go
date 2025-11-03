package handler

import (
	"context"
	"log"

	chatv1 "github.com/go-chat/chat/pkg/api/chat/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// GetChat retrieves chat information
// Returns dummy chat until database integration is added
func (s *Server) GetChat(ctx context.Context, req *chatv1.GetChatRequest) (*chatv1.GetChatResponse, error) {
	log.Printf("GetChat called for chat_id: %s", req.ChatId)

	// TODO: Extract authenticated user_id from JWT
	// TODO: Query chat from database by chat_id
	// TODO: Return NOT_FOUND if chat doesn't exist
	// TODO: Verify user is a participant (PERMISSION_DENIED if not)

	return &chatv1.GetChatResponse{
		Chat: &chatv1.Chat{
			ChatId: req.ChatId,
			ParticipantIds: []string{
				"550e8400-e29b-41d4-a716-446655440001",
				"550e8400-e29b-41d4-a716-446655440002",
			},
			CreatedAt: timestamppb.Now(),
		},
	}, nil
}

// ListUserChats lists all chats for a user with cursor-based pagination
// Returns dummy chats until database integration is added
func (s *Server) ListUserChats(ctx context.Context, req *chatv1.ListUserChatsRequest) (*chatv1.ListUserChatsResponse, error) {
	log.Printf("ListUserChats called for user_id: %s, limit: %d", req.UserId, req.Limit)

	// TODO: Extract authenticated user_id from JWT
	// TODO: Verify req.UserId matches authenticated user (or allow service-to-service calls)
	// TODO: Decode cursor to get pagination position
	// TODO: Query chats from database where user is participant
	// TODO: Order by last_message_at DESC for most recent chats first
	// TODO: Apply limit + 1 to check for more results
	// TODO: Encode next_cursor if more results exist

	return &chatv1.ListUserChatsResponse{
		Chats: []*chatv1.Chat{
			{
				ChatId: "550e8400-e29b-41d4-a716-446655440001",
				ParticipantIds: []string{
					req.UserId,
					"550e8400-e29b-41d4-a716-446655440002",
				},
				CreatedAt: timestamppb.Now(),
			},
		},
		NextCursor: "", // Empty means no more results
	}, nil
}

// ListChatMembers lists all participants in a chat
// Returns dummy participants until database integration is added
func (s *Server) ListChatMembers(ctx context.Context, req *chatv1.ListChatMembersRequest) (*chatv1.ListChatMembersResponse, error) {
	log.Printf("ListChatMembers called for chat_id: %s", req.ChatId)

	// TODO: Extract authenticated user_id from JWT
	// TODO: Query chat from database
	// TODO: Return NOT_FOUND if chat doesn't exist
	// TODO: Verify user is a participant (PERMISSION_DENIED if not)
	// TODO: Return participant user IDs

	return &chatv1.ListChatMembersResponse{
		UserIds: []string{
			"550e8400-e29b-41d4-a716-446655440001",
			"550e8400-e29b-41d4-a716-446655440002",
		},
	}, nil
}

