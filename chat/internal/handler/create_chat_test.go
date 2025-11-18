package handler

import (
	"context"
	"errors"
	"testing"

	"github.com/go-chat/chat/internal/domain"
	chatv1 "github.com/go-chat/chat/pkg/api/chat/v1"
)

func TestCreateDirectChat_ValidRequest_ReturnsChatID(t *testing.T) {
	expectedChatID := domain.NewChatID("550e8400-e29b-41d4-a716-446655440000")

	mockChatSvc := &mockChatService{
		createDirectChatFunc: func(ctx context.Context, requesterID, participantID domain.UserID) (domain.ChatID, error) {
			return expectedChatID, nil
		},
	}

	server := NewServer(mockChatSvc, nil)
	req := &chatv1.CreateDirectChatRequest{
		ParticipantId: "550e8400-e29b-41d4-a716-446655440001",
	}

	resp, err := server.CreateDirectChat(context.Background(), req)

	if err != nil {
		t.Fatalf("CreateDirectChat() returned error: %v", err)
	}

	if resp.ChatId != expectedChatID.String() {
		t.Errorf("Expected chat ID '%s', got '%s'", expectedChatID.String(), resp.ChatId)
	}
}

func TestCreateDirectChat_UsersNotFriends_ReturnsError(t *testing.T) {
	mockChatSvc := &mockChatService{
		createDirectChatFunc: func(ctx context.Context, requesterID, participantID domain.UserID) (domain.ChatID, error) {
			return "", domain.ErrUsersNotFriends
		},
	}

	server := NewServer(mockChatSvc, nil)
	req := &chatv1.CreateDirectChatRequest{
		ParticipantId: "550e8400-e29b-41d4-a716-446655440001",
	}

	_, err := server.CreateDirectChat(context.Background(), req)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if !errors.Is(err, domain.ErrUsersNotFriends) {
		t.Errorf("Expected ErrUsersNotFriends, got: %v", err)
	}
}

func TestCreateDirectChat_ChatAlreadyExists_ReturnsError(t *testing.T) {
	mockChatSvc := &mockChatService{
		createDirectChatFunc: func(ctx context.Context, requesterID, participantID domain.UserID) (domain.ChatID, error) {
			return "", domain.ErrChatAlreadyExists
		},
	}

	server := NewServer(mockChatSvc, nil)
	req := &chatv1.CreateDirectChatRequest{
		ParticipantId: "550e8400-e29b-41d4-a716-446655440001",
	}

	_, err := server.CreateDirectChat(context.Background(), req)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if !errors.Is(err, domain.ErrChatAlreadyExists) {
		t.Errorf("Expected ErrChatAlreadyExists, got: %v", err)
	}
}

