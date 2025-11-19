package handler

import (
	"context"
	"errors"
	"testing"

	"github.com/go-chat/auth/internal/domain"
	authv1 "github.com/go-chat/auth/pkg/api/auth/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// mockTokenService is a mock implementation of service.TokenService
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

func TestGetPublicKeys_ValidRequest_ReturnsKeys(t *testing.T) {
	expectedKeys := []*domain.PublicKey{
		{
			Kid: "key-2024-01",
			Alg: "RS256",
			Use: "sig",
			N:   "modulus_base64url_encoded",
			E:   "AQAB",
		},
		{
			Kid: "key-2024-02",
			Alg: "RS256",
			Use: "sig",
			N:   "another_modulus_base64url_encoded",
			E:   "AQAB",
		},
	}

	mockToken := &mockTokenService{
		getPublicKeysFunc: func(ctx context.Context) ([]*domain.PublicKey, error) {
			return expectedKeys, nil
		},
	}

	server := NewServer(nil, mockToken)
	req := &authv1.GetPublicKeysRequest{}

	resp, err := server.GetPublicKeys(context.Background(), req)

	if err != nil {
		t.Fatalf("GetPublicKeys() returned error: %v", err)
	}

	if len(resp.Keys) != len(expectedKeys) {
		t.Fatalf("Expected %d keys, got %d", len(expectedKeys), len(resp.Keys))
	}

	for i, key := range resp.Keys {
		expected := expectedKeys[i]
		if key.Kid != expected.Kid {
			t.Errorf("Key[%d]: Expected Kid '%s', got '%s'", i, expected.Kid, key.Kid)
		}
		if key.Alg != expected.Alg {
			t.Errorf("Key[%d]: Expected Alg '%s', got '%s'", i, expected.Alg, key.Alg)
		}
		if key.Use != expected.Use {
			t.Errorf("Key[%d]: Expected Use '%s', got '%s'", i, expected.Use, key.Use)
		}
		if key.N != expected.N {
			t.Errorf("Key[%d]: Expected N '%s', got '%s'", i, expected.N, key.N)
		}
		if key.E != expected.E {
			t.Errorf("Key[%d]: Expected E '%s', got '%s'", i, expected.E, key.E)
		}
	}
}

func TestGetPublicKeys_EmptyKeys_ReturnsEmptyArray(t *testing.T) {
	mockToken := &mockTokenService{
		getPublicKeysFunc: func(ctx context.Context) ([]*domain.PublicKey, error) {
			return []*domain.PublicKey{}, nil
		},
	}

	server := NewServer(nil, mockToken)
	req := &authv1.GetPublicKeysRequest{}

	resp, err := server.GetPublicKeys(context.Background(), req)

	if err != nil {
		t.Fatalf("GetPublicKeys() returned error: %v", err)
	}

	if len(resp.Keys) != 0 {
		t.Errorf("Expected empty keys array, got %d keys", len(resp.Keys))
	}
}

func TestGetPublicKeys_ServiceError_ReturnsInternalError(t *testing.T) {
	mockToken := &mockTokenService{
		getPublicKeysFunc: func(ctx context.Context) ([]*domain.PublicKey, error) {
			return nil, errors.New("failed to load keys")
		},
	}

	server := NewServer(nil, mockToken)
	req := &authv1.GetPublicKeysRequest{}

	_, err := server.GetPublicKeys(context.Background(), req)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	st, ok := status.FromError(err)
	if !ok {
		t.Fatal("Expected gRPC status error")
	}

	if st.Code() != codes.Internal {
		t.Errorf("Expected code Internal, got %v", st.Code())
	}

	if st.Message() != "failed to retrieve public keys" {
		t.Errorf("Expected message 'failed to retrieve public keys', got '%s'", st.Message())
	}
}

func TestGetPublicKeys_SingleKey_ReturnsOneKey(t *testing.T) {
	expectedKey := &domain.PublicKey{
		Kid: "key-single",
		Alg: "RS256",
		Use: "sig",
		N:   "single_key_modulus",
		E:   "AQAB",
	}

	mockToken := &mockTokenService{
		getPublicKeysFunc: func(ctx context.Context) ([]*domain.PublicKey, error) {
			return []*domain.PublicKey{expectedKey}, nil
		},
	}

	server := NewServer(nil, mockToken)
	req := &authv1.GetPublicKeysRequest{}

	resp, err := server.GetPublicKeys(context.Background(), req)

	if err != nil {
		t.Fatalf("GetPublicKeys() returned error: %v", err)
	}

	if len(resp.Keys) != 1 {
		t.Fatalf("Expected 1 key, got %d", len(resp.Keys))
	}

	if resp.Keys[0].Kid != expectedKey.Kid {
		t.Errorf("Expected Kid '%s', got '%s'", expectedKey.Kid, resp.Keys[0].Kid)
	}
}
