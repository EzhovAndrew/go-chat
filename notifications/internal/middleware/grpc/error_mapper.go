package grpc

import (
	"context"
	"errors"

	"github.com/go-chat/notifications/internal/domain"
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
	case errors.Is(err, domain.ErrNotificationNotFound):
		return status.Error(codes.NotFound, "notification not found")
	case errors.Is(err, domain.ErrPermissionDenied):
		return status.Error(codes.PermissionDenied, "permission denied")
	default:
		// Log internal error details here if needed
		// For now, return a generic internal error
		return status.Error(codes.Internal, "internal server error")
	}
}

