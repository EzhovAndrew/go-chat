package handler

import (
	"context"

	"github.com/go-chat/chat/internal/domain"
	"github.com/go-chat/chat/internal/dto"
	chatv1 "github.com/go-chat/chat/pkg/api/chat/v1"
)

// GetChat retrieves chat information
func (s *Server) GetChat(ctx context.Context, req *chatv1.GetChatRequest) (*chatv1.GetChatResponse, error) {
	// TODO: Extract authenticated user_id from JWT
	requesterID := domain.NewUserID("authenticated-user-id") // Placeholder

	// Delegate to service layer
	chat, err := s.chatService.GetChat(ctx, domain.NewChatID(req.ChatId), requesterID)
	if err != nil {
		return nil, err // Middleware will map domain error to gRPC status
	}

	// Convert domain model to proto message
	return &chatv1.GetChatResponse{
		Chat: dto.ToProtoChat(chat),
	}, nil
}

// ListUserChats lists all chats for a user with cursor-based pagination
func (s *Server) ListUserChats(ctx context.Context, req *chatv1.ListUserChatsRequest) (*chatv1.ListUserChatsResponse, error) {
	// Default limit if not provided
	limit := req.Limit
	if limit == 0 {
		limit = 20
	}

	// Delegate to service layer
	chats, nextCursor, err := s.chatService.ListUserChats(
		ctx,
		domain.NewUserID(req.UserId),
		req.Cursor,
		limit,
	)
	if err != nil {
		return nil, err // Middleware will map domain error to gRPC status
	}

	// Convert domain models to proto messages
	return &chatv1.ListUserChatsResponse{
		Chats:      dto.ToProtoChats(chats),
		NextCursor: nextCursor,
	}, nil
}

// ListChatMembers lists all participants in a chat
func (s *Server) ListChatMembers(ctx context.Context, req *chatv1.ListChatMembersRequest) (*chatv1.ListChatMembersResponse, error) {
	// TODO: Extract authenticated user_id from JWT
	requesterID := domain.NewUserID("authenticated-user-id") // Placeholder

	// Delegate to service layer
	members, err := s.chatService.ListChatMembers(ctx, domain.NewChatID(req.ChatId), requesterID)
	if err != nil {
		return nil, err // Middleware will map domain error to gRPC status
	}

	// Convert domain types to strings
	userIDs := make([]string, len(members))
	for i, id := range members {
		userIDs[i] = id.String()
	}

	return &chatv1.ListChatMembersResponse{
		UserIds: userIDs,
	}, nil
}
