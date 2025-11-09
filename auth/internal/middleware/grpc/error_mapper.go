package grpc

import (
	"context"
	"errors"

	"github.com/go-chat/auth/internal/domain"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ErrorMapperInterceptor returns a unary server interceptor that maps domain errors to gRPC status codes
func ErrorMapperInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		resp, err := handler(ctx, req)
		if err != nil {
			return resp, mapDomainError(err)
		}
		return resp, nil
	}
}

// mapDomainError converts domain errors to appropriate gRPC status codes
func mapDomainError(err error) error {
	switch {
	case errors.Is(err, domain.ErrEmailAlreadyExists):
		return status.Error(codes.AlreadyExists, "email already registered")
	case errors.Is(err, domain.ErrInvalidCredentials):
		return status.Error(codes.Unauthenticated, "invalid email or password")
	case errors.Is(err, domain.ErrInvalidToken),
		errors.Is(err, domain.ErrTokenExpired),
		errors.Is(err, domain.ErrTokenRevoked):
		return status.Error(codes.Unauthenticated, "invalid or expired token")
	default:
		// Log internal error details here if needed
		// For now, return a generic internal error
		return status.Error(codes.Internal, "internal server error")
	}
}

