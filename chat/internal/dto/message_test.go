package dto

import (
	"testing"
	"time"

	"github.com/go-chat/chat/internal/domain"
)

func TestToProtoMessage_ValidMessage_ConvertsCorrectly(t *testing.T) {
	message := &domain.Message{
		MessageID: domain.NewMessageID("550e8400-e29b-41d4-a716-446655440000"),
		ChatID:    domain.NewChatID("550e8400-e29b-41d4-a716-446655440001"),
		SenderID:  domain.NewUserID("550e8400-e29b-41d4-a716-446655440002"),
		Text:      "Hello, World!",
		CreatedAt: time.Now(),
	}

	protoMessage := ToProtoMessage(message)

	if protoMessage.MessageId != "550e8400-e29b-41d4-a716-446655440000" {
		t.Errorf("Expected MessageId '550e8400-e29b-41d4-a716-446655440000', got '%s'", protoMessage.MessageId)
	}

	if protoMessage.ChatId != "550e8400-e29b-41d4-a716-446655440001" {
		t.Errorf("Expected ChatId '550e8400-e29b-41d4-a716-446655440001', got '%s'", protoMessage.ChatId)
	}

	if protoMessage.SenderId != "550e8400-e29b-41d4-a716-446655440002" {
		t.Errorf("Expected SenderId '550e8400-e29b-41d4-a716-446655440002', got '%s'", protoMessage.SenderId)
	}

	if protoMessage.Text != "Hello, World!" {
		t.Errorf("Expected Text 'Hello, World!', got '%s'", protoMessage.Text)
	}
}

func TestToProtoMessages_MultipleMessages_ConvertsAll(t *testing.T) {
	messages := []*domain.Message{
		{
			MessageID: domain.NewMessageID("msg1"),
			ChatID:    domain.NewChatID("chat1"),
			SenderID:  domain.NewUserID("user1"),
			Text:      "First message",
			CreatedAt: time.Now(),
		},
		{
			MessageID: domain.NewMessageID("msg2"),
			ChatID:    domain.NewChatID("chat1"),
			SenderID:  domain.NewUserID("user2"),
			Text:      "Second message",
			CreatedAt: time.Now(),
		},
	}

	protoMessages := ToProtoMessages(messages)

	if len(protoMessages) != 2 {
		t.Fatalf("Expected 2 messages, got %d", len(protoMessages))
	}

	if protoMessages[0].Text != "First message" {
		t.Errorf("Expected first message text 'First message', got '%s'", protoMessages[0].Text)
	}

	if protoMessages[1].Text != "Second message" {
		t.Errorf("Expected second message text 'Second message', got '%s'", protoMessages[1].Text)
	}
}

func TestToProtoMessages_EmptySlice_ReturnsEmptySlice(t *testing.T) {
	messages := []*domain.Message{}

	protoMessages := ToProtoMessages(messages)

	if len(protoMessages) != 0 {
		t.Errorf("Expected empty slice, got %d messages", len(protoMessages))
	}
}

