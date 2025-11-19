package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/go-chat/auth/internal/domain"
)

// Mock repositories for testing
type mockUserRepository struct {
	createFunc     func(ctx context.Context, user *domain.User) error
	getByIDFunc    func(ctx context.Context, userID domain.UserID) (*domain.User, error)
	getByEmailFunc func(ctx context.Context, email string) (*domain.User, error)
}

func (m *mockUserRepository) Create(ctx context.Context, user *domain.User) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, user)
	}
	return errors.New("not implemented")
}

func (m *mockUserRepository) GetByID(ctx context.Context, userID domain.UserID) (*domain.User, error) {
	if m.getByIDFunc != nil {
		return m.getByIDFunc(ctx, userID)
	}
	return nil, errors.New("not implemented")
}

func (m *mockUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	if m.getByEmailFunc != nil {
		return m.getByEmailFunc(ctx, email)
	}
	return nil, errors.New("not implemented")
}

type mockTokenService struct {
	generateTokenPairFunc             func(ctx context.Context, userID domain.UserID, email string) (*domain.TokenPair, *domain.RefreshToken, error)
	storeRefreshTokenFunc             func(ctx context.Context, refreshToken *domain.RefreshToken) error
	validateAndRevokeRefreshTokenFunc func(ctx context.Context, refreshToken string) (domain.UserID, error)
	getPublicKeysFunc                 func(ctx context.Context) ([]*domain.PublicKey, error)
}

func (m *mockTokenService) GenerateTokenPair(ctx context.Context, userID domain.UserID, email string) (*domain.TokenPair, *domain.RefreshToken, error) {
	if m.generateTokenPairFunc != nil {
		return m.generateTokenPairFunc(ctx, userID, email)
	}
	return nil, nil, errors.New("not implemented")
}

func (m *mockTokenService) StoreRefreshToken(ctx context.Context, refreshToken *domain.RefreshToken) error {
	if m.storeRefreshTokenFunc != nil {
		return m.storeRefreshTokenFunc(ctx, refreshToken)
	}
	return errors.New("not implemented")
}

func (m *mockTokenService) ValidateAndRevokeRefreshToken(ctx context.Context, refreshToken string) (domain.UserID, error) {
	if m.validateAndRevokeRefreshTokenFunc != nil {
		return m.validateAndRevokeRefreshTokenFunc(ctx, refreshToken)
	}
	return "", errors.New("not implemented")
}

func (m *mockTokenService) GetPublicKeys(ctx context.Context) ([]*domain.PublicKey, error) {
	if m.getPublicKeysFunc != nil {
		return m.getPublicKeysFunc(ctx)
	}
	return nil, errors.New("not implemented")
}

func TestNewAuthService_NilUserRepo_Panics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic with nil userRepo")
		}
	}()
	NewAuthService(nil, &mockTokenService{})
}

func TestNewAuthService_NilTokenService_Panics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected panic with nil tokenService")
		}
	}()
	NewAuthService(&mockUserRepository{}, nil)
}

func TestRegister_ValidInput_ReturnsUser(t *testing.T) {
	mockUserRepo := &mockUserRepository{
		createFunc: func(ctx context.Context, user *domain.User) error {
			return nil
		},
	}

	service := NewAuthService(mockUserRepo, &mockTokenService{})

	user, err := service.Register(context.Background(), "test@example.com", "password123")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if user == nil {
		t.Fatal("Expected user, got nil")
	}

	if user.Email != "test@example.com" {
		t.Errorf("Expected email 'test@example.com', got '%s'", user.Email)
	}
}

func TestRegister_EmailAlreadyExists_ReturnsError(t *testing.T) {
	mockUserRepo := &mockUserRepository{
		createFunc: func(ctx context.Context, user *domain.User) error {
			return domain.ErrEmailAlreadyExists
		},
	}

	service := NewAuthService(mockUserRepo, &mockTokenService{})

	_, err := service.Register(context.Background(), "test@example.com", "password123")
	if !errors.Is(err, domain.ErrEmailAlreadyExists) {
		t.Errorf("Expected ErrEmailAlreadyExists, got: %v", err)
	}
}

func TestRegister_CreateUserFails_ReturnsError(t *testing.T) {
	expectedErr := errors.New("database error")
	mockUserRepo := &mockUserRepository{
		createFunc: func(ctx context.Context, user *domain.User) error {
			return expectedErr
		},
	}

	service := NewAuthService(mockUserRepo, &mockTokenService{})

	_, err := service.Register(context.Background(), "test@example.com", "password123")
	if !errors.Is(err, expectedErr) {
		t.Errorf("Expected database error, got: %v", err)
	}
}

func TestLogin_UserNotFound_ReturnsInvalidCredentials(t *testing.T) {
	mockUserRepo := &mockUserRepository{
		getByEmailFunc: func(ctx context.Context, email string) (*domain.User, error) {
			return nil, errors.New("user not found")
		},
	}

	service := NewAuthService(mockUserRepo, &mockTokenService{})

	_, _, err := service.Login(context.Background(), "test@example.com", "password123")
	if !errors.Is(err, domain.ErrInvalidCredentials) {
		t.Errorf("Expected ErrInvalidCredentials, got: %v", err)
	}
}

func TestLogin_TokenGenerationFails_ReturnsError(t *testing.T) {
	mockUserRepo := &mockUserRepository{
		getByEmailFunc: func(ctx context.Context, email string) (*domain.User, error) {
			passwordHash := "$argon2id$v=19$m=65536,t=3,p=2$somesalt$somehash"
			return &domain.User{
				ID:           domain.NewUserID("user-123"),
				Email:        email,
				PasswordHash: passwordHash,
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			}, nil
		},
	}

	expectedErr := errors.New("token generation failed")
	mockTokenService := &mockTokenService{
		generateTokenPairFunc: func(ctx context.Context, userID domain.UserID, email string) (*domain.TokenPair, *domain.RefreshToken, error) {
			return nil, nil, expectedErr
		},
	}

	service := NewAuthService(mockUserRepo, mockTokenService)

	_, _, err := service.Login(context.Background(), "test@example.com", "password")
	if err == nil {
		t.Error("Expected error")
	}
}

func TestRefresh_ValidToken_ReturnsNewTokenPair(t *testing.T) {
	userID := domain.NewUserID("user-123")

	mockUserRepo := &mockUserRepository{
		getByIDFunc: func(ctx context.Context, uid domain.UserID) (*domain.User, error) {
			return &domain.User{
				ID:           userID,
				Email:        "test@example.com",
				PasswordHash: "hash",
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			}, nil
		},
	}

	mockTokenService := &mockTokenService{
		validateAndRevokeRefreshTokenFunc: func(ctx context.Context, refreshToken string) (domain.UserID, error) {
			return userID, nil
		},
		generateTokenPairFunc: func(ctx context.Context, uid domain.UserID, email string) (*domain.TokenPair, *domain.RefreshToken, error) {
			return &domain.TokenPair{
					AccessToken:  "new-access-token",
					RefreshToken: "new-refresh-token-jwt",
				}, &domain.RefreshToken{
					ID:        "new-token-id",
					UserID:    uid,
					Token:     "new-jti-hash",
					ExpiresAt: time.Now().Add(30 * 24 * time.Hour),
					Revoked:   false,
					CreatedAt: time.Now(),
				}, nil
		},
		storeRefreshTokenFunc: func(ctx context.Context, refreshToken *domain.RefreshToken) error {
			return nil
		},
	}

	service := NewAuthService(mockUserRepo, mockTokenService)

	tokenPair, returnedUserID, err := service.Refresh(context.Background(), "valid-refresh-token-jwt")
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if tokenPair == nil {
		t.Fatal("Expected token pair, got nil")
	}

	if returnedUserID != userID {
		t.Errorf("Expected user ID '%s', got '%s'", userID, returnedUserID)
	}

	if tokenPair.AccessToken != "new-access-token" {
		t.Errorf("Expected access token 'new-access-token', got '%s'", tokenPair.AccessToken)
	}
}

func TestRefresh_ExpiredToken_ReturnsError(t *testing.T) {
	mockTokenService := &mockTokenService{
		validateAndRevokeRefreshTokenFunc: func(ctx context.Context, refreshToken string) (domain.UserID, error) {
			return "", domain.ErrTokenExpired
		},
	}

	service := NewAuthService(&mockUserRepository{}, mockTokenService)

	_, _, err := service.Refresh(context.Background(), "expired-token-jwt")
	if !errors.Is(err, domain.ErrTokenExpired) {
		t.Errorf("Expected ErrTokenExpired, got: %v", err)
	}
}

func TestRefresh_RevokedToken_ReturnsError(t *testing.T) {
	mockTokenService := &mockTokenService{
		validateAndRevokeRefreshTokenFunc: func(ctx context.Context, refreshToken string) (domain.UserID, error) {
			return "", domain.ErrTokenRevoked
		},
	}

	service := NewAuthService(&mockUserRepository{}, mockTokenService)

	_, _, err := service.Refresh(context.Background(), "revoked-token-jwt")
	if !errors.Is(err, domain.ErrTokenRevoked) {
		t.Errorf("Expected ErrTokenRevoked, got: %v", err)
	}
}

func TestRefresh_InvalidToken_ReturnsError(t *testing.T) {
	mockTokenService := &mockTokenService{
		validateAndRevokeRefreshTokenFunc: func(ctx context.Context, refreshToken string) (domain.UserID, error) {
			return "", domain.ErrInvalidToken
		},
	}

	service := NewAuthService(&mockUserRepository{}, mockTokenService)

	_, _, err := service.Refresh(context.Background(), "invalid-token-jwt")
	if !errors.Is(err, domain.ErrInvalidToken) {
		t.Errorf("Expected ErrInvalidToken, got: %v", err)
	}
}

func TestRefresh_RevokeFails_ReturnsError(t *testing.T) {
	expectedErr := errors.New("revoke failed")
	mockTokenService := &mockTokenService{
		validateAndRevokeRefreshTokenFunc: func(ctx context.Context, refreshToken string) (domain.UserID, error) {
			return "", expectedErr
		},
	}

	service := NewAuthService(&mockUserRepository{}, mockTokenService)

	_, _, err := service.Refresh(context.Background(), "valid-token-jwt")
	if err == nil {
		t.Error("Expected error when revoke fails")
	}
}
