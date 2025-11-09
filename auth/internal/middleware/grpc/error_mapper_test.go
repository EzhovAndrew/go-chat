package grpc

import (
	"context"
	"errors"
	"testing"

	"github.com/go-chat/auth/internal/domain"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestErrorMapperInterceptor_NoDomainError_ReturnsOriginalResponse(t *testing.T) {
	interceptor := ErrorMapperInterceptor()
	expectedResp := "test response"

	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return expectedResp, nil
	}

	resp, err := interceptor(context.Background(), nil, &grpc.UnaryServerInfo{}, handler)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if resp != expectedResp {
		t.Errorf("Expected response %v, got %v", expectedResp, resp)
	}
}

func TestMapDomainError_EmailAlreadyExists_ReturnsAlreadyExists(t *testing.T) {
	err := mapDomainError(domain.ErrEmailAlreadyExists)

	st, ok := status.FromError(err)
	if !ok {
		t.Fatal("Expected gRPC status error")
	}

	if st.Code() != codes.AlreadyExists {
		t.Errorf("Expected code AlreadyExists, got %v", st.Code())
	}

	if st.Message() != "email already registered" {
		t.Errorf("Expected message 'email already registered', got '%s'", st.Message())
	}
}

func TestMapDomainError_InvalidCredentials_ReturnsUnauthenticated(t *testing.T) {
	err := mapDomainError(domain.ErrInvalidCredentials)

	st, ok := status.FromError(err)
	if !ok {
		t.Fatal("Expected gRPC status error")
	}

	if st.Code() != codes.Unauthenticated {
		t.Errorf("Expected code Unauthenticated, got %v", st.Code())
	}

	if st.Message() != "invalid email or password" {
		t.Errorf("Expected message 'invalid email or password', got '%s'", st.Message())
	}
}

func TestMapDomainError_TokenErrors_ReturnsUnauthenticated(t *testing.T) {
	tests := []struct {
		name string
		err  error
	}{
		{"invalid token", domain.ErrInvalidToken},
		{"expired token", domain.ErrTokenExpired},
		{"revoked token", domain.ErrTokenRevoked},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := mapDomainError(tt.err)

			st, ok := status.FromError(err)
			if !ok {
				t.Fatal("Expected gRPC status error")
			}

			if st.Code() != codes.Unauthenticated {
				t.Errorf("Expected code Unauthenticated, got %v", st.Code())
			}

			if st.Message() != "invalid or expired token" {
				t.Errorf("Expected message 'invalid or expired token', got '%s'", st.Message())
			}
		})
	}
}

func TestMapDomainError_WrappedTokenError_ReturnsUnauthenticated(t *testing.T) {
	wrappedErr := errors.New("some context: " + domain.ErrInvalidToken.Error())
	// This won't match with errors.Is, so let's test proper wrapping
	properlyWrappedErr := errors.Join(errors.New("some context"), domain.ErrInvalidToken)

	err := mapDomainError(properlyWrappedErr)

	st, ok := status.FromError(err)
	if !ok {
		t.Fatal("Expected gRPC status error")
	}

	if st.Code() != codes.Unauthenticated {
		t.Errorf("Expected code Unauthenticated for wrapped error, got %v", st.Code())
	}

	// For non-wrapped string concatenation, it should return Internal
	err2 := mapDomainError(wrappedErr)
	st2, _ := status.FromError(err2)
	if st2.Code() != codes.Internal {
		t.Errorf("Expected code Internal for improperly wrapped error, got %v", st2.Code())
	}
}

func TestMapDomainError_UnknownError_ReturnsInternal(t *testing.T) {
	unknownErr := errors.New("some unknown error")

	err := mapDomainError(unknownErr)

	st, ok := status.FromError(err)
	if !ok {
		t.Fatal("Expected gRPC status error")
	}

	if st.Code() != codes.Internal {
		t.Errorf("Expected code Internal, got %v", st.Code())
	}

	if st.Message() != "internal server error" {
		t.Errorf("Expected message 'internal server error', got '%s'", st.Message())
	}
}

func TestErrorMapperInterceptor_DomainError_ReturnsMappedError(t *testing.T) {
	interceptor := ErrorMapperInterceptor()

	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, domain.ErrEmailAlreadyExists
	}

	_, err := interceptor(context.Background(), nil, &grpc.UnaryServerInfo{}, handler)

	if err == nil {
		t.Fatal("Expected error to be returned")
	}

	st, ok := status.FromError(err)
	if !ok {
		t.Fatal("Expected gRPC status error")
	}

	if st.Code() != codes.AlreadyExists {
		t.Errorf("Expected code AlreadyExists, got %v", st.Code())
	}
}

