package handler

import (
	"context"
	"errors"
	"time"

	"github.com/go-chat/chat/internal/domain"
)

// mockChatService is a mock implementation of service.ChatService
type mockChatService struct {
	createDirectChatFunc func(ctx context.Context, requesterID, participantID domain.UserID) (domain.ChatID, error)
	getChatFunc          func(ctx context.Context, chatID domain.ChatID, requesterID domain.UserID) (*domain.Chat, error)
	listUserChatsFunc    func(ctx context.Context, userID domain.UserID, cursor string, limit int32) ([]*domain.Chat, string, error)
	listChatMembersFunc  func(ctx context.Context, chatID domain.ChatID, requesterID domain.UserID) ([]domain.UserID, error)
}

func (m *mockChatService) CreateDirectChat(ctx context.Context, requesterID, participantID domain.UserID) (domain.ChatID, error) {
	if m.createDirectChatFunc != nil {
		return m.createDirectChatFunc(ctx, requesterID, participantID)
	}
	return "", errors.New("not implemented")
}

func (m *mockChatService) GetChat(ctx context.Context, chatID domain.ChatID, requesterID domain.UserID) (*domain.Chat, error) {
	if m.getChatFunc != nil {
		return m.getChatFunc(ctx, chatID, requesterID)
	}
	return nil, errors.New("not implemented")
}

func (m *mockChatService) ListUserChats(ctx context.Context, userID domain.UserID, cursor string, limit int32) ([]*domain.Chat, string, error) {
	if m.listUserChatsFunc != nil {
		return m.listUserChatsFunc(ctx, userID, cursor, limit)
	}
	return nil, "", errors.New("not implemented")
}

func (m *mockChatService) ListChatMembers(ctx context.Context, chatID domain.ChatID, requesterID domain.UserID) ([]domain.UserID, error) {
	if m.listChatMembersFunc != nil {
		return m.listChatMembersFunc(ctx, chatID, requesterID)
	}
	return nil, errors.New("not implemented")
}

// mockMessageService is a mock implementation of service.MessageService
type mockMessageService struct {
	sendMessageFunc    func(ctx context.Context, chatID domain.ChatID, senderID domain.UserID, text, idempotencyKey string) (*domain.Message, error)
	listMessagesFunc   func(ctx context.Context, chatID domain.ChatID, requesterID domain.UserID, cursor string, limit int32) ([]*domain.Message, string, error)
	streamMessagesFunc func(ctx context.Context, chatID domain.ChatID, requesterID domain.UserID, since time.Time) error
}

func (m *mockMessageService) SendMessage(ctx context.Context, chatID domain.ChatID, senderID domain.UserID, text, idempotencyKey string) (*domain.Message, error) {
	if m.sendMessageFunc != nil {
		return m.sendMessageFunc(ctx, chatID, senderID, text, idempotencyKey)
	}
	return nil, errors.New("not implemented")
}

func (m *mockMessageService) ListMessages(ctx context.Context, chatID domain.ChatID, requesterID domain.UserID, cursor string, limit int32) ([]*domain.Message, string, error) {
	if m.listMessagesFunc != nil {
		return m.listMessagesFunc(ctx, chatID, requesterID, cursor, limit)
	}
	return nil, "", errors.New("not implemented")
}

func (m *mockMessageService) StreamMessages(ctx context.Context, chatID domain.ChatID, requesterID domain.UserID, since time.Time) error {
	if m.streamMessagesFunc != nil {
		return m.streamMessagesFunc(ctx, chatID, requesterID, since)
	}
	return errors.New("not implemented")
}

