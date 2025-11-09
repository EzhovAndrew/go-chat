package handler

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/go-chat/chat/internal/domain"
	chatv1 "github.com/go-chat/chat/pkg/api/chat/v1"
)

func TestSendMessage_ValidRequest_ReturnsMessage(t *testing.T) {
	expectedMessage := &domain.Message{
		MessageID: domain.NewMessageID("550e8400-e29b-41d4-a716-446655440000"),
		ChatID:    domain.NewChatID("550e8400-e29b-41d4-a716-446655440001"),
		SenderID:  domain.NewUserID("550e8400-e29b-41d4-a716-446655440002"),
		Text:      "Hello, world!",
		CreatedAt: time.Now(),
	}

	mockMsgSvc := &mockMessageService{
		sendMessageFunc: func(ctx context.Context, chatID domain.ChatID, senderID domain.UserID, text, idempotencyKey string) (*domain.Message, error) {
			if text != "Hello, world!" {
				t.Errorf("Expected text 'Hello, world!', got '%s'", text)
			}
			return expectedMessage, nil
		},
	}

	server := NewServer(nil, mockMsgSvc)
	req := &chatv1.SendMessageRequest{
		ChatId:         "550e8400-e29b-41d4-a716-446655440001",
		Text:           "Hello, world!",
		IdempotencyKey: "550e8400-e29b-41d4-a716-446655440003",
	}

	resp, err := server.SendMessage(context.Background(), req)

	if err != nil {
		t.Fatalf("SendMessage() returned error: %v", err)
	}

	if resp.Message.MessageId != expectedMessage.MessageID.String() {
		t.Errorf("Expected message ID '%s', got '%s'", expectedMessage.MessageID.String(), resp.Message.MessageId)
	}

	if resp.Message.Text != expectedMessage.Text {
		t.Errorf("Expected text '%s', got '%s'", expectedMessage.Text, resp.Message.Text)
	}
}

func TestSendMessage_InvalidMessage_ReturnsError(t *testing.T) {
	mockMsgSvc := &mockMessageService{
		sendMessageFunc: func(ctx context.Context, chatID domain.ChatID, senderID domain.UserID, text, idempotencyKey string) (*domain.Message, error) {
			return nil, domain.ErrInvalidMessage
		},
	}

	server := NewServer(nil, mockMsgSvc)
	req := &chatv1.SendMessageRequest{
		ChatId:         "550e8400-e29b-41d4-a716-446655440001",
		Text:           "",
		IdempotencyKey: "550e8400-e29b-41d4-a716-446655440003",
	}

	_, err := server.SendMessage(context.Background(), req)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if !errors.Is(err, domain.ErrInvalidMessage) {
		t.Errorf("Expected ErrInvalidMessage, got: %v", err)
	}
}

func TestSendMessage_PermissionDenied_ReturnsError(t *testing.T) {
	mockMsgSvc := &mockMessageService{
		sendMessageFunc: func(ctx context.Context, chatID domain.ChatID, senderID domain.UserID, text, idempotencyKey string) (*domain.Message, error) {
			return nil, domain.ErrPermissionDenied
		},
	}

	server := NewServer(nil, mockMsgSvc)
	req := &chatv1.SendMessageRequest{
		ChatId:         "550e8400-e29b-41d4-a716-446655440001",
		Text:           "Hello",
		IdempotencyKey: "550e8400-e29b-41d4-a716-446655440003",
	}

	_, err := server.SendMessage(context.Background(), req)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if !errors.Is(err, domain.ErrPermissionDenied) {
		t.Errorf("Expected ErrPermissionDenied, got: %v", err)
	}
}

func TestListMessages_ValidRequest_ReturnsMessages(t *testing.T) {
	expectedMessages := []*domain.Message{
		{
			MessageID: domain.NewMessageID("550e8400-e29b-41d4-a716-446655440000"),
			ChatID:    domain.NewChatID("550e8400-e29b-41d4-a716-446655440001"),
			SenderID:  domain.NewUserID("550e8400-e29b-41d4-a716-446655440002"),
			Text:      "Message 1",
			CreatedAt: time.Now(),
		},
		{
			MessageID: domain.NewMessageID("550e8400-e29b-41d4-a716-446655440003"),
			ChatID:    domain.NewChatID("550e8400-e29b-41d4-a716-446655440001"),
			SenderID:  domain.NewUserID("550e8400-e29b-41d4-a716-446655440004"),
			Text:      "Message 2",
			CreatedAt: time.Now(),
		},
	}

	mockMsgSvc := &mockMessageService{
		listMessagesFunc: func(ctx context.Context, chatID domain.ChatID, requesterID domain.UserID, cursor string, limit int32) ([]*domain.Message, string, error) {
			return expectedMessages, "next_cursor_456", nil
		},
	}

	server := NewServer(nil, mockMsgSvc)
	req := &chatv1.ListMessagesRequest{
		ChatId: "550e8400-e29b-41d4-a716-446655440001",
		Limit:  10,
	}

	resp, err := server.ListMessages(context.Background(), req)

	if err != nil {
		t.Fatalf("ListMessages() returned error: %v", err)
	}

	if len(resp.Messages) != 2 {
		t.Errorf("Expected 2 messages, got %d", len(resp.Messages))
	}

	if resp.NextCursor != "next_cursor_456" {
		t.Errorf("Expected cursor 'next_cursor_456', got '%s'", resp.NextCursor)
	}
}

func TestListMessages_WithDefaultLimit_UsesDefaultValue(t *testing.T) {
	mockMsgSvc := &mockMessageService{
		listMessagesFunc: func(ctx context.Context, chatID domain.ChatID, requesterID domain.UserID, cursor string, limit int32) ([]*domain.Message, string, error) {
			if limit != 50 {
				t.Errorf("Expected default limit 50, got %d", limit)
			}
			return []*domain.Message{}, "", nil
		},
	}

	server := NewServer(nil, mockMsgSvc)
	req := &chatv1.ListMessagesRequest{
		ChatId: "550e8400-e29b-41d4-a716-446655440001",
		Limit:  0, // Not provided, should default to 50
	}

	_, err := server.ListMessages(context.Background(), req)

	if err != nil {
		t.Fatalf("ListMessages() returned error: %v", err)
	}
}

