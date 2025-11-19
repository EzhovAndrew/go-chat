package service

import (
	"context"
	"fmt"
	"time"

	"github.com/go-chat/auth/internal/domain"
	"github.com/go-chat/auth/internal/repository"
	"github.com/go-chat/auth/internal/utils"
	"github.com/google/uuid"
)

// authService implements the AuthService interface
type authService struct {
	userRepo     repository.UserRepository
	tokenService TokenService
}

// NewAuthService creates a new auth service with injected dependencies
func NewAuthService(
	userRepo repository.UserRepository,
	tokenService TokenService,
) AuthService {
	if userRepo == nil {
		panic("userRepo cannot be nil")
	}
	if tokenService == nil {
		panic("tokenService cannot be nil")
	}

	return &authService{
		userRepo:     userRepo,
		tokenService: tokenService,
	}
}

// Register creates a new user account
func (s *authService) Register(ctx context.Context, email, password string) (*domain.User, error) {
	// Hash password
	passwordHash, err := utils.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	// Create user entity
	now := time.Now()
	user := &domain.User{
		ID:           domain.NewUserID(uuid.New().String()),
		Email:        email,
		PasswordHash: passwordHash,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	// Store user - repository will return ErrEmailAlreadyExists on unique constraint violation
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// Login authenticates user and returns tokens with user ID
func (s *authService) Login(ctx context.Context, email, password string) (*domain.TokenPair, domain.UserID, error) {
	// Fetch user by email
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, "", domain.ErrInvalidCredentials
	}

	// Compare password
	if err := utils.ComparePassword(user.PasswordHash, password); err != nil {
		return nil, "", domain.ErrInvalidCredentials
	}

	// Generate token pair with refresh token metadata
	tokenPair, refreshTokenMetadata, err := s.tokenService.GenerateTokenPair(ctx, user.ID, user.Email)
	if err != nil {
		return nil, "", fmt.Errorf("generate tokens: %w", err)
	}

	// Store refresh token metadata in database
	if err := s.tokenService.StoreRefreshToken(ctx, refreshTokenMetadata); err != nil {
		return nil, "", fmt.Errorf("store refresh token: %w", err)
	}

	return tokenPair, user.ID, nil
}

// Refresh validates refresh token and returns new token pair with user ID
func (s *authService) Refresh(ctx context.Context, refreshToken string) (*domain.TokenPair, domain.UserID, error) {
	// Validate and revoke old refresh token
	userID, err := s.tokenService.ValidateAndRevokeRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, "", err
	}

	// Get user to generate new tokens
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, "", fmt.Errorf("get user: %w", err)
	}

	// Generate new token pair with refresh token metadata
	newTokenPair, newRefreshTokenMetadata, err := s.tokenService.GenerateTokenPair(ctx, user.ID, user.Email)
	if err != nil {
		return nil, "", fmt.Errorf("generate new tokens: %w", err)
	}

	// Store new refresh token metadata
	if err := s.tokenService.StoreRefreshToken(ctx, newRefreshTokenMetadata); err != nil {
		return nil, "", fmt.Errorf("store new refresh token: %w", err)
	}

	return newTokenPair, user.ID, nil
}
