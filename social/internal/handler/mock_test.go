package handler

import (
	"context"
	"errors"

	"github.com/go-chat/social/internal/domain"
)

// mockFriendRequestService is a mock implementation of service.FriendRequestService
type mockFriendRequestService struct {
	sendFriendRequestFunc    func(ctx context.Context, requesterID, targetID domain.UserID) (*domain.FriendRequest, error)
	listRequestsFunc         func(ctx context.Context, userID domain.UserID, cursor string, limit int32) ([]*domain.FriendRequest, string, error)
	acceptFriendRequestFunc  func(ctx context.Context, requestID domain.RequestID, userID domain.UserID) (*domain.FriendRequest, error)
	declineFriendRequestFunc func(ctx context.Context, requestID domain.RequestID, userID domain.UserID) error
}

func (m *mockFriendRequestService) SendFriendRequest(ctx context.Context, requesterID, targetID domain.UserID) (*domain.FriendRequest, error) {
	if m.sendFriendRequestFunc != nil {
		return m.sendFriendRequestFunc(ctx, requesterID, targetID)
	}
	return nil, errors.New("not implemented")
}

func (m *mockFriendRequestService) ListRequests(ctx context.Context, userID domain.UserID, cursor string, limit int32) ([]*domain.FriendRequest, string, error) {
	if m.listRequestsFunc != nil {
		return m.listRequestsFunc(ctx, userID, cursor, limit)
	}
	return nil, "", errors.New("not implemented")
}

func (m *mockFriendRequestService) AcceptFriendRequest(ctx context.Context, requestID domain.RequestID, userID domain.UserID) (*domain.FriendRequest, error) {
	if m.acceptFriendRequestFunc != nil {
		return m.acceptFriendRequestFunc(ctx, requestID, userID)
	}
	return nil, errors.New("not implemented")
}

func (m *mockFriendRequestService) DeclineFriendRequest(ctx context.Context, requestID domain.RequestID, userID domain.UserID) error {
	if m.declineFriendRequestFunc != nil {
		return m.declineFriendRequestFunc(ctx, requestID, userID)
	}
	return errors.New("not implemented")
}

// mockFriendshipService is a mock implementation of service.FriendshipService
type mockFriendshipService struct {
	listFriendsFunc  func(ctx context.Context, userID domain.UserID, cursor string, limit int32) ([]domain.UserID, string, error)
	removeFriendFunc func(ctx context.Context, userID, friendID domain.UserID) error
}

func (m *mockFriendshipService) ListFriends(ctx context.Context, userID domain.UserID, cursor string, limit int32) ([]domain.UserID, string, error) {
	if m.listFriendsFunc != nil {
		return m.listFriendsFunc(ctx, userID, cursor, limit)
	}
	return nil, "", errors.New("not implemented")
}

func (m *mockFriendshipService) RemoveFriend(ctx context.Context, userID, friendID domain.UserID) error {
	if m.removeFriendFunc != nil {
		return m.removeFriendFunc(ctx, userID, friendID)
	}
	return errors.New("not implemented")
}

// mockBlockService is a mock implementation of service.BlockService
type mockBlockService struct {
	blockUserFunc   func(ctx context.Context, blockerID, blockedID domain.UserID) error
	unblockUserFunc func(ctx context.Context, blockerID, blockedID domain.UserID) error
}

func (m *mockBlockService) BlockUser(ctx context.Context, blockerID, blockedID domain.UserID) error {
	if m.blockUserFunc != nil {
		return m.blockUserFunc(ctx, blockerID, blockedID)
	}
	return errors.New("not implemented")
}

func (m *mockBlockService) UnblockUser(ctx context.Context, blockerID, blockedID domain.UserID) error {
	if m.unblockUserFunc != nil {
		return m.unblockUserFunc(ctx, blockerID, blockedID)
	}
	return errors.New("not implemented")
}

// mockRelationshipService is a mock implementation of service.RelationshipService
type mockRelationshipService struct {
	checkRelationshipFunc func(ctx context.Context, userID, targetID domain.UserID) (domain.RelationshipStatus, error)
}

func (m *mockRelationshipService) CheckRelationship(ctx context.Context, userID, targetID domain.UserID) (domain.RelationshipStatus, error) {
	if m.checkRelationshipFunc != nil {
		return m.checkRelationshipFunc(ctx, userID, targetID)
	}
	return domain.RelationshipStatusNone, errors.New("not implemented")
}

