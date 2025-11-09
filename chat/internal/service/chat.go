package service

import (
	"context"

	"github.com/go-chat/chat/internal/domain"
)

// ChatService handles chat management operations
type ChatService interface {
	// CreateDirectChat creates a 1-on-1 chat between two users
	CreateDirectChat(ctx context.Context, requesterID, participantID domain.UserID) (domain.ChatID, error)

	// GetChat retrieves chat information
	GetChat(ctx context.Context, chatID domain.ChatID, requesterID domain.UserID) (*domain.Chat, error)

	// ListUserChats lists all chats for a user with pagination
	ListUserChats(ctx context.Context, userID domain.UserID, cursor string, limit int32) ([]*domain.Chat, string, error)

	// ListChatMembers lists all participants in a chat
	ListChatMembers(ctx context.Context, chatID domain.ChatID, requesterID domain.UserID) ([]domain.UserID, error)
}

