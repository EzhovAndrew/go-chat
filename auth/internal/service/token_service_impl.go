package service

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/go-chat/auth/internal/domain"
	"github.com/go-chat/auth/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// tokenService implements the TokenService interface
type tokenService struct {
	privateKey       *rsa.PrivateKey
	publicKey        *rsa.PublicKey
	keyID            string
	refreshTokenRepo repository.RefreshTokenRepository
}

// NewTokenService creates a new token service with RSA key loading and token repository
func NewTokenService(privateKeyPath string, refreshTokenRepo repository.RefreshTokenRepository) (TokenService, error) {
	if refreshTokenRepo == nil {
		panic("refreshTokenRepo cannot be nil")
	}
	// Load and parse RSA private key
	privateKey, err := loadRSAPrivateKey(privateKeyPath)
	if err != nil {
		return nil, err
	}

	// Generate a key ID for JWK
	keyID := uuid.New().String()

	return &tokenService{
		privateKey:       privateKey,
		publicKey:        &privateKey.PublicKey,
		keyID:            keyID,
		refreshTokenRepo: refreshTokenRepo,
	}, nil
}

// Interface methods

// GenerateTokenPair creates JWT access token and JWT refresh token
// Returns the token pair and refresh token metadata for efficient storage
func (s *tokenService) GenerateTokenPair(ctx context.Context, userID domain.UserID, email string) (*domain.TokenPair, *domain.RefreshToken, error) {
	now := time.Now()

	// Generate access token (15 minutes)
	accessTokenClaims := jwt.MapClaims{
		"sub":   userID.String(),
		"email": email,
		"iat":   now.Unix(),
		"exp":   now.Add(15 * time.Minute).Unix(),
		"type":  "access",
	}

	accessTokenString, err := s.signJWT(accessTokenClaims)
	if err != nil {
		return nil, nil, fmt.Errorf("sign access token: %w", err)
	}

	// Generate JWT refresh token (30 days)
	jti := uuid.New().String()
	expiresAt := now.Add(30 * 24 * time.Hour)

	refreshTokenClaims := jwt.MapClaims{
		"sub":  userID.String(),
		"jti":  jti, // Unique token ID for revocation and collision detection
		"iat":  now.Unix(),
		"exp":  expiresAt.Unix(),
		"type": "refresh",
	}

	refreshTokenString, err := s.signJWT(refreshTokenClaims)
	if err != nil {
		return nil, nil, fmt.Errorf("sign refresh token: %w", err)
	}

	tokenPair := &domain.TokenPair{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
	}

	// Prepare refresh token metadata for storage (without parsing)
	refreshTokenMetadata := &domain.RefreshToken{
		ID:        uuid.New().String(),
		UserID:    userID,
		Token:     jti, // Store JTI directly (UUID is already cryptographically random)
		ExpiresAt: expiresAt,
		Revoked:   false,
		CreatedAt: now,
	}

	return tokenPair, refreshTokenMetadata, nil
}

// StoreRefreshToken stores the refresh token metadata in repository
func (s *tokenService) StoreRefreshToken(ctx context.Context, refreshToken *domain.RefreshToken) error {
	if err := s.refreshTokenRepo.Create(ctx, refreshToken); err != nil {
		return fmt.Errorf("store refresh token: %w", err)
	}
	return nil
}

// ValidateAndRevokeRefreshToken validates refresh token, checks database, and revokes it
// Returns user ID if valid, error otherwise
func (s *tokenService) ValidateAndRevokeRefreshToken(ctx context.Context, refreshToken string) (domain.UserID, error) {
	// Parse and validate JWT signature first (fail fast)
	claims, err := s.parseRefreshToken(refreshToken)
	if err != nil {
		return "", domain.ErrInvalidToken
	}

	// Extract jti (JWT ID) for database lookup
	jti, ok := claims["jti"].(string)
	if !ok {
		return "", domain.ErrInvalidToken
	}

	// Extract user ID from claims for cross-validation
	claimsUserID, ok := claims["sub"].(string)
	if !ok {
		return "", domain.ErrInvalidToken
	}

	// Look up token in database using jti
	storedToken, err := s.refreshTokenRepo.GetByToken(ctx, jti)
	if err != nil {
		return "", domain.ErrInvalidToken
	}

	// CRITICAL: Cross-validate user ID from JWT claims against stored user ID
	// This detects JTI collisions and tampering
	if storedToken.UserID.String() != claimsUserID {
		return "", domain.ErrInvalidToken
	}

	if time.Now().After(storedToken.ExpiresAt) {
		return "", domain.ErrTokenExpired
	}

	if storedToken.Revoked {
		return "", domain.ErrTokenRevoked
	}

	if err := s.refreshTokenRepo.Revoke(ctx, jti); err != nil {
		return "", fmt.Errorf("revoke token: %w", err)
	}

	return storedToken.UserID, nil
}

// GetPublicKeys returns public keys in JWK format for JWT validation
func (s *tokenService) GetPublicKeys(ctx context.Context) ([]*domain.PublicKey, error) {
	// Convert RSA public key to JWK format
	n := s.publicKey.N
	e := s.publicKey.E

	// Encode N and E to base64url
	nBytes := n.Bytes()
	nBase64 := base64.RawURLEncoding.EncodeToString(nBytes)

	eBytes := big.NewInt(int64(e)).Bytes()
	eBase64 := base64.RawURLEncoding.EncodeToString(eBytes)

	publicKey := &domain.PublicKey{
		Kid: s.keyID,
		Alg: "RS256",
		Use: "sig",
		N:   nBase64,
		E:   eBase64,
	}

	return []*domain.PublicKey{publicKey}, nil
}

// Helper functions

// parseRefreshToken parses and validates a JWT refresh token
// Returns the claims if valid, otherwise returns an error
func (s *tokenService) parseRefreshToken(tokenString string) (map[string]interface{}, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Verify signing method
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.publicKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("parse token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid claims format")
	}

	// Verify token type
	tokenType, ok := claims["type"].(string)
	if !ok || tokenType != "refresh" {
		return nil, fmt.Errorf("invalid token type")
	}

	return claims, nil
}

// signJWT creates and signs a JWT token with the given claims
func (s *tokenService) signJWT(claims jwt.MapClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = s.keyID

	tokenString, err := token.SignedString(s.privateKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// loadRSAPrivateKey loads and parses an RSA private key from a PEM file
// Supports both PKCS1 and PKCS8 formats
func loadRSAPrivateKey(path string) (*rsa.PrivateKey, error) {
	// Read key file
	keyData, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read private key file: %w", err)
	}

	// Decode PEM block
	block, _ := pem.Decode(keyData)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block containing private key")
	}

	// Try PKCS1 format first
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err == nil {
		return privateKey, nil
	}

	// Try PKCS8 format
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("parse private key (tried PKCS1 and PKCS8): %w", err)
	}

	// Ensure it's an RSA key
	rsaKey, ok := key.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("key is not RSA private key")
	}

	return rsaKey, nil
}
