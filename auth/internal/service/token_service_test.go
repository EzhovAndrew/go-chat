package service

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"math/big"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-chat/auth/internal/domain"
	"github.com/golang-jwt/jwt/v5"
)

// mockRefreshTokenRepository for token service tests
type mockRefreshTokenRepository struct {
	createFunc     func(ctx context.Context, token *domain.RefreshToken) error
	getByTokenFunc func(ctx context.Context, jtiHash string) (*domain.RefreshToken, error)
	revokeFunc     func(ctx context.Context, jtiHash string) error
}

func (m *mockRefreshTokenRepository) Create(ctx context.Context, token *domain.RefreshToken) error {
	if m.createFunc != nil {
		return m.createFunc(ctx, token)
	}
	return errors.New("not implemented")
}

func (m *mockRefreshTokenRepository) GetByToken(ctx context.Context, jtiHash string) (*domain.RefreshToken, error) {
	if m.getByTokenFunc != nil {
		return m.getByTokenFunc(ctx, jtiHash)
	}
	return nil, errors.New("not implemented")
}

func (m *mockRefreshTokenRepository) Revoke(ctx context.Context, jtiHash string) error {
	if m.revokeFunc != nil {
		return m.revokeFunc(ctx, jtiHash)
	}
	return errors.New("not implemented")
}

func (m *mockRefreshTokenRepository) DeleteExpired(ctx context.Context, userID domain.UserID) error {
	return errors.New("not implemented")
}

func TestNewTokenService_InvalidKeyPath_ReturnsError(t *testing.T) {
	_, err := NewTokenService("/nonexistent/path/to/key.pem", &mockRefreshTokenRepository{})
	if err == nil {
		t.Error("Expected error for invalid key path")
	}
}

func TestNewTokenService_InvalidKeyFormat_ReturnsError(t *testing.T) {
	// Create temporary file with invalid content
	tmpDir := t.TempDir()
	keyPath := filepath.Join(tmpDir, "invalid_key.pem")
	if err := os.WriteFile(keyPath, []byte("invalid key content"), 0600); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	_, err := NewTokenService(keyPath, &mockRefreshTokenRepository{})
	if err == nil {
		t.Error("Expected error for invalid key format")
	}
}

func TestNewTokenService_ValidKey_ReturnsService(t *testing.T) {
	// Generate a test RSA key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("Failed to generate RSA key: %v", err)
	}

	// Create temporary file with valid PEM key
	tmpDir := t.TempDir()
	keyPath := filepath.Join(tmpDir, "test_key.pem")

	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	pemBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	}

	if err := os.WriteFile(keyPath, pem.EncodeToMemory(pemBlock), 0600); err != nil {
		t.Fatalf("Failed to write test key file: %v", err)
	}

	service, err := NewTokenService(keyPath, &mockRefreshTokenRepository{})
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if service == nil {
		t.Fatal("Expected service, got nil")
	}
}

func TestGenerateTokenPair_ValidInput_ReturnsTokens(t *testing.T) {
	// Generate a test RSA key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("Failed to generate RSA key: %v", err)
	}

	// Create temporary file with valid PEM key
	tmpDir := t.TempDir()
	keyPath := filepath.Join(tmpDir, "test_key.pem")

	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	pemBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	}

	if err := os.WriteFile(keyPath, pem.EncodeToMemory(pemBlock), 0600); err != nil {
		t.Fatalf("Failed to write test key file: %v", err)
	}

	service, err := NewTokenService(keyPath, &mockRefreshTokenRepository{})
	if err != nil {
		t.Fatalf("Failed to create service: %v", err)
	}

	userID := domain.NewUserID("user-123")
	email := "test@example.com"

	tokenPair, refreshTokenMetadata, err := service.GenerateTokenPair(context.Background(), userID, email)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if tokenPair == nil {
		t.Fatal("Expected token pair, got nil")
	}

	if refreshTokenMetadata == nil {
		t.Fatal("Expected refresh token metadata, got nil")
	}

	if refreshTokenMetadata.UserID != userID {
		t.Errorf("Expected user ID %s, got %s", userID, refreshTokenMetadata.UserID)
	}

	if tokenPair.AccessToken == "" {
		t.Error("Expected access token to be set")
	}

	if tokenPair.RefreshToken == "" {
		t.Error("Expected refresh token to be set")
	}

	// Verify the access token can be parsed
	token, err := jwt.Parse(tokenPair.AccessToken, func(token *jwt.Token) (interface{}, error) {
		return &privateKey.PublicKey, nil
	})

	if err != nil {
		t.Errorf("Failed to parse access token: %v", err)
	}

	if !token.Valid {
		t.Error("Expected token to be valid")
	}

	// Verify claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		t.Fatal("Failed to parse claims")
	}

	if claims["sub"] != userID.String() {
		t.Errorf("Expected sub '%s', got '%v'", userID, claims["sub"])
	}

	if claims["email"] != email {
		t.Errorf("Expected email '%s', got '%v'", email, claims["email"])
	}

	if claims["type"] != "access" {
		t.Errorf("Expected type 'access', got '%v'", claims["type"])
	}
}

func TestGetPublicKeys_ReturnsJWK(t *testing.T) {
	// Generate a test RSA key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("Failed to generate RSA key: %v", err)
	}

	// Create temporary file with valid PEM key
	tmpDir := t.TempDir()
	keyPath := filepath.Join(tmpDir, "test_key.pem")

	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	pemBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	}

	if err := os.WriteFile(keyPath, pem.EncodeToMemory(pemBlock), 0600); err != nil {
		t.Fatalf("Failed to write test key file: %v", err)
	}

	service, err := NewTokenService(keyPath, &mockRefreshTokenRepository{})
	if err != nil {
		t.Fatalf("Failed to create service: %v", err)
	}

	publicKeys, err := service.GetPublicKeys(context.Background())
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if len(publicKeys) == 0 {
		t.Fatal("Expected at least one public key")
	}

	key := publicKeys[0]

	if key.Kid == "" {
		t.Error("Expected key ID to be set")
	}

	if key.Alg != "RS256" {
		t.Errorf("Expected algorithm 'RS256', got '%s'", key.Alg)
	}

	if key.Use != "sig" {
		t.Errorf("Expected use 'sig', got '%s'", key.Use)
	}

	if key.N == "" {
		t.Error("Expected N (modulus) to be set")
	}

	if key.E == "" {
		t.Error("Expected E (exponent) to be set")
	}

	// Verify N and E can be decoded
	if _, err := base64.RawURLEncoding.DecodeString(key.N); err != nil {
		t.Errorf("Failed to decode N: %v", err)
	}

	if _, err := base64.RawURLEncoding.DecodeString(key.E); err != nil {
		t.Errorf("Failed to decode E: %v", err)
	}
}

func TestGetPublicKeys_MatchesPrivateKey(t *testing.T) {
	// Generate a test RSA key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("Failed to generate RSA key: %v", err)
	}

	// Create temporary file with valid PEM key
	tmpDir := t.TempDir()
	keyPath := filepath.Join(tmpDir, "test_key.pem")

	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	pemBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	}

	if err := os.WriteFile(keyPath, pem.EncodeToMemory(pemBlock), 0600); err != nil {
		t.Fatalf("Failed to write test key file: %v", err)
	}

	service, err := NewTokenService(keyPath, &mockRefreshTokenRepository{})
	if err != nil {
		t.Fatalf("Failed to create service: %v", err)
	}

	publicKeys, err := service.GetPublicKeys(context.Background())
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	key := publicKeys[0]

	// Decode N and E
	nBytes, err := base64.RawURLEncoding.DecodeString(key.N)
	if err != nil {
		t.Fatalf("Failed to decode N: %v", err)
	}

	eBytes, err := base64.RawURLEncoding.DecodeString(key.E)
	if err != nil {
		t.Fatalf("Failed to decode E: %v", err)
	}

	// Verify N matches
	n := new(big.Int).SetBytes(nBytes)
	if n.Cmp(privateKey.PublicKey.N) != 0 {
		t.Error("Public key N does not match")
	}

	// Verify E matches
	e := new(big.Int).SetBytes(eBytes)
	if e.Int64() != int64(privateKey.PublicKey.E) {
		t.Errorf("Public key E does not match: expected %d, got %d", privateKey.PublicKey.E, e.Int64())
	}
}
