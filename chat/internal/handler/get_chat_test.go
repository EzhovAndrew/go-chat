package handler

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/go-chat/chat/internal/domain"
	chatv1 "github.com/go-chat/chat/pkg/api/chat/v1"
)

func TestGetChat_ValidRequest_ReturnsChat(t *testing.T) {
	expectedChat := &domain.Chat{
		ChatID: domain.NewChatID("550e8400-e29b-41d4-a716-446655440000"),
		ParticipantIDs: []domain.UserID{
			domain.NewUserID("550e8400-e29b-41d4-a716-446655440001"),
			domain.NewUserID("550e8400-e29b-41d4-a716-446655440002"),
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockChatSvc := &mockChatService{
		getChatFunc: func(ctx context.Context, chatID domain.ChatID, requesterID domain.UserID) (*domain.Chat, error) {
			return expectedChat, nil
		},
	}

	server := NewServer(mockChatSvc, nil)
	req := &chatv1.GetChatRequest{
		ChatId: "550e8400-e29b-41d4-a716-446655440000",
	}

	resp, err := server.GetChat(context.Background(), req)

	if err != nil {
		t.Fatalf("GetChat() returned error: %v", err)
	}

	if resp.Chat.ChatId != expectedChat.ChatID.String() {
		t.Errorf("Expected chat ID '%s', got '%s'", expectedChat.ChatID.String(), resp.Chat.ChatId)
	}

	if len(resp.Chat.ParticipantIds) != len(expectedChat.ParticipantIDs) {
		t.Errorf("Expected %d participants, got %d", len(expectedChat.ParticipantIDs), len(resp.Chat.ParticipantIds))
	}
}

func TestGetChat_ChatNotFound_ReturnsError(t *testing.T) {
	mockChatSvc := &mockChatService{
		getChatFunc: func(ctx context.Context, chatID domain.ChatID, requesterID domain.UserID) (*domain.Chat, error) {
			return nil, domain.ErrChatNotFound
		},
	}

	server := NewServer(mockChatSvc, nil)
	req := &chatv1.GetChatRequest{
		ChatId: "550e8400-e29b-41d4-a716-446655440000",
	}

	_, err := server.GetChat(context.Background(), req)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if !errors.Is(err, domain.ErrChatNotFound) {
		t.Errorf("Expected ErrChatNotFound, got: %v", err)
	}
}

func TestGetChat_PermissionDenied_ReturnsError(t *testing.T) {
	mockChatSvc := &mockChatService{
		getChatFunc: func(ctx context.Context, chatID domain.ChatID, requesterID domain.UserID) (*domain.Chat, error) {
			return nil, domain.ErrPermissionDenied
		},
	}

	server := NewServer(mockChatSvc, nil)
	req := &chatv1.GetChatRequest{
		ChatId: "550e8400-e29b-41d4-a716-446655440000",
	}

	_, err := server.GetChat(context.Background(), req)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if !errors.Is(err, domain.ErrPermissionDenied) {
		t.Errorf("Expected ErrPermissionDenied, got: %v", err)
	}
}

func TestListUserChats_ValidRequest_ReturnsChats(t *testing.T) {
	expectedChats := []*domain.Chat{
		{
			ChatID: domain.NewChatID("550e8400-e29b-41d4-a716-446655440000"),
			ParticipantIDs: []domain.UserID{
				domain.NewUserID("550e8400-e29b-41d4-a716-446655440001"),
				domain.NewUserID("550e8400-e29b-41d4-a716-446655440002"),
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	mockChatSvc := &mockChatService{
		listUserChatsFunc: func(ctx context.Context, userID domain.UserID, cursor string, limit int32) ([]*domain.Chat, string, error) {
			return expectedChats, "next_cursor_123", nil
		},
	}

	server := NewServer(mockChatSvc, nil)
	req := &chatv1.ListUserChatsRequest{
		UserId: "550e8400-e29b-41d4-a716-446655440001",
		Limit:  10,
	}

	resp, err := server.ListUserChats(context.Background(), req)

	if err != nil {
		t.Fatalf("ListUserChats() returned error: %v", err)
	}

	if len(resp.Chats) != 1 {
		t.Errorf("Expected 1 chat, got %d", len(resp.Chats))
	}

	if resp.NextCursor != "next_cursor_123" {
		t.Errorf("Expected cursor 'next_cursor_123', got '%s'", resp.NextCursor)
	}
}

func TestListChatMembers_ValidRequest_ReturnsMembers(t *testing.T) {
	expectedMembers := []domain.UserID{
		domain.NewUserID("550e8400-e29b-41d4-a716-446655440001"),
		domain.NewUserID("550e8400-e29b-41d4-a716-446655440002"),
	}

	mockChatSvc := &mockChatService{
		listChatMembersFunc: func(ctx context.Context, chatID domain.ChatID, requesterID domain.UserID) ([]domain.UserID, error) {
			return expectedMembers, nil
		},
	}

	server := NewServer(mockChatSvc, nil)
	req := &chatv1.ListChatMembersRequest{
		ChatId: "550e8400-e29b-41d4-a716-446655440000",
	}

	resp, err := server.ListChatMembers(context.Background(), req)

	if err != nil {
		t.Fatalf("ListChatMembers() returned error: %v", err)
	}

	if len(resp.UserIds) != 2 {
		t.Errorf("Expected 2 members, got %d", len(resp.UserIds))
	}
}

