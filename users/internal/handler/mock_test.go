package handler

import (
	"context"
	"errors"

	"github.com/go-chat/users/internal/domain"
)

// mockUserService is a mock implementation of service.UserService
type mockUserService struct {
	createProfileFunc        func(ctx context.Context, userID domain.UserID, nickname, bio, avatarURL string) (*domain.UserProfile, error)
	updateProfileFunc        func(ctx context.Context, userID domain.UserID, nickname, bio, avatarURL string) (*domain.UserProfile, error)
	getProfileByIDFunc       func(ctx context.Context, userID domain.UserID) (*domain.UserProfile, error)
	getProfilesByIDsFunc     func(ctx context.Context, userIDs []domain.UserID) ([]*domain.UserProfile, error)
	getProfileByNicknameFunc func(ctx context.Context, nickname string) (*domain.UserProfile, error)
	searchByNicknameFunc     func(ctx context.Context, query, cursor string, limit int32) ([]*domain.UserProfile, string, error)
}

func (m *mockUserService) CreateProfile(ctx context.Context, userID domain.UserID, nickname, bio, avatarURL string) (*domain.UserProfile, error) {
	if m.createProfileFunc != nil {
		return m.createProfileFunc(ctx, userID, nickname, bio, avatarURL)
	}
	return nil, errors.New("not implemented")
}

func (m *mockUserService) UpdateProfile(ctx context.Context, userID domain.UserID, nickname, bio, avatarURL string) (*domain.UserProfile, error) {
	if m.updateProfileFunc != nil {
		return m.updateProfileFunc(ctx, userID, nickname, bio, avatarURL)
	}
	return nil, errors.New("not implemented")
}

func (m *mockUserService) GetProfileByID(ctx context.Context, userID domain.UserID) (*domain.UserProfile, error) {
	if m.getProfileByIDFunc != nil {
		return m.getProfileByIDFunc(ctx, userID)
	}
	return nil, errors.New("not implemented")
}

func (m *mockUserService) GetProfilesByIDs(ctx context.Context, userIDs []domain.UserID) ([]*domain.UserProfile, error) {
	if m.getProfilesByIDsFunc != nil {
		return m.getProfilesByIDsFunc(ctx, userIDs)
	}
	return nil, errors.New("not implemented")
}

func (m *mockUserService) GetProfileByNickname(ctx context.Context, nickname string) (*domain.UserProfile, error) {
	if m.getProfileByNicknameFunc != nil {
		return m.getProfileByNicknameFunc(ctx, nickname)
	}
	return nil, errors.New("not implemented")
}

func (m *mockUserService) SearchByNickname(ctx context.Context, query, cursor string, limit int32) ([]*domain.UserProfile, string, error) {
	if m.searchByNicknameFunc != nil {
		return m.searchByNicknameFunc(ctx, query, cursor, limit)
	}
	return nil, "", errors.New("not implemented")
}
