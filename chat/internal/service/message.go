package service

import (
	"context"
	"time"

	"github.com/go-chat/chat/internal/domain"
)

// MessageService handles message operations
type MessageService interface {
	// SendMessage sends a message to a chat
	SendMessage(ctx context.Context, chatID domain.ChatID, senderID domain.UserID, text, idempotencyKey string) (*domain.Message, error)

	// ListMessages retrieves message history for a chat with pagination
	ListMessages(ctx context.Context, chatID domain.ChatID, requesterID domain.UserID, cursor string, limit int32) ([]*domain.Message, string, error)

	// StreamMessages streams new messages in real-time (placeholder for now)
	StreamMessages(ctx context.Context, chatID domain.ChatID, requesterID domain.UserID, since time.Time) error
}

