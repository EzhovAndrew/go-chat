package handler

import (
	"context"
	"errors"
	"testing"

	"github.com/go-chat/social/internal/domain"
	socialv1 "github.com/go-chat/social/pkg/api/social/v1"
)

func TestBlockUser_ValidRequest_ReturnsSuccess(t *testing.T) {
	mockBlockSvc := &mockBlockService{
		blockUserFunc: func(ctx context.Context, blockerID, blockedID domain.UserID) error {
			if blockedID.String() != "550e8400-e29b-41d4-a716-446655440002" {
				t.Errorf("Expected blocked ID '550e8400-e29b-41d4-a716-446655440002', got '%s'", blockedID.String())
			}
			return nil
		},
	}

	server := NewServer(nil, nil, mockBlockSvc, nil)
	req := &socialv1.BlockUserRequest{
		TargetUserId: "550e8400-e29b-41d4-a716-446655440002",
	}

	resp, err := server.BlockUser(context.Background(), req)

	if err != nil {
		t.Fatalf("BlockUser() returned error: %v", err)
	}

	if resp == nil {
		t.Fatal("Expected response, got nil")
	}
}

func TestBlockUser_SelfAction_ReturnsError(t *testing.T) {
	mockBlockSvc := &mockBlockService{
		blockUserFunc: func(ctx context.Context, blockerID, blockedID domain.UserID) error {
			return domain.ErrSelfAction
		},
	}

	server := NewServer(nil, nil, mockBlockSvc, nil)
	req := &socialv1.BlockUserRequest{
		TargetUserId: "550e8400-e29b-41d4-a716-446655440001",
	}

	_, err := server.BlockUser(context.Background(), req)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if !errors.Is(err, domain.ErrSelfAction) {
		t.Errorf("Expected ErrSelfAction, got: %v", err)
	}
}

func TestUnblockUser_ValidRequest_ReturnsSuccess(t *testing.T) {
	mockBlockSvc := &mockBlockService{
		unblockUserFunc: func(ctx context.Context, blockerID, blockedID domain.UserID) error {
			if blockedID.String() != "550e8400-e29b-41d4-a716-446655440002" {
				t.Errorf("Expected blocked ID '550e8400-e29b-41d4-a716-446655440002', got '%s'", blockedID.String())
			}
			return nil
		},
	}

	server := NewServer(nil, nil, mockBlockSvc, nil)
	req := &socialv1.UnblockUserRequest{
		TargetUserId: "550e8400-e29b-41d4-a716-446655440002",
	}

	resp, err := server.UnblockUser(context.Background(), req)

	if err != nil {
		t.Fatalf("UnblockUser() returned error: %v", err)
	}

	if resp == nil {
		t.Fatal("Expected response, got nil")
	}
}

func TestUnblockUser_UserNotBlocked_ReturnsError(t *testing.T) {
	mockBlockSvc := &mockBlockService{
		unblockUserFunc: func(ctx context.Context, blockerID, blockedID domain.UserID) error {
			return domain.ErrUserBlocked
		},
	}

	server := NewServer(nil, nil, mockBlockSvc, nil)
	req := &socialv1.UnblockUserRequest{
		TargetUserId: "550e8400-e29b-41d4-a716-446655440002",
	}

	_, err := server.UnblockUser(context.Background(), req)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if !errors.Is(err, domain.ErrUserBlocked) {
		t.Errorf("Expected ErrUserBlocked, got: %v", err)
	}
}

