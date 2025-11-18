package dto

import (
	"testing"
	"time"

	"github.com/go-chat/chat/internal/domain"
)

func TestToProtoChat_ValidChat_ConvertsCorrectly(t *testing.T) {
	chat := &domain.Chat{
		ChatID: domain.NewChatID("550e8400-e29b-41d4-a716-446655440000"),
		ParticipantIDs: []domain.UserID{
			domain.NewUserID("550e8400-e29b-41d4-a716-446655440001"),
			domain.NewUserID("550e8400-e29b-41d4-a716-446655440002"),
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	protoChat := ToProtoChat(chat)

	if protoChat.ChatId != "550e8400-e29b-41d4-a716-446655440000" {
		t.Errorf("Expected ChatId '550e8400-e29b-41d4-a716-446655440000', got '%s'", protoChat.ChatId)
	}

	if len(protoChat.ParticipantIds) != 2 {
		t.Fatalf("Expected 2 participants, got %d", len(protoChat.ParticipantIds))
	}

	if protoChat.ParticipantIds[0] != "550e8400-e29b-41d4-a716-446655440001" {
		t.Errorf("Expected first participant '550e8400-e29b-41d4-a716-446655440001', got '%s'", protoChat.ParticipantIds[0])
	}
}

func TestToProtoChats_MultipleChats_ConvertsAll(t *testing.T) {
	chats := []*domain.Chat{
		{
			ChatID: domain.NewChatID("550e8400-e29b-41d4-a716-446655440000"),
			ParticipantIDs: []domain.UserID{
				domain.NewUserID("user1"),
				domain.NewUserID("user2"),
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			ChatID: domain.NewChatID("550e8400-e29b-41d4-a716-446655440001"),
			ParticipantIDs: []domain.UserID{
				domain.NewUserID("user3"),
				domain.NewUserID("user4"),
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	protoChats := ToProtoChats(chats)

	if len(protoChats) != 2 {
		t.Fatalf("Expected 2 chats, got %d", len(protoChats))
	}

	if protoChats[0].ChatId != "550e8400-e29b-41d4-a716-446655440000" {
		t.Errorf("Expected first chat ID '550e8400-e29b-41d4-a716-446655440000', got '%s'", protoChats[0].ChatId)
	}

	if protoChats[1].ChatId != "550e8400-e29b-41d4-a716-446655440001" {
		t.Errorf("Expected second chat ID '550e8400-e29b-41d4-a716-446655440001', got '%s'", protoChats[1].ChatId)
	}
}

func TestToProtoChats_EmptySlice_ReturnsEmptySlice(t *testing.T) {
	chats := []*domain.Chat{}

	protoChats := ToProtoChats(chats)

	if len(protoChats) != 0 {
		t.Errorf("Expected empty slice, got %d chats", len(protoChats))
	}
}

