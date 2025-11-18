package grpc

import (
	"context"
	"errors"

	"github.com/go-chat/chat/internal/domain"
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
	case errors.Is(err, domain.ErrChatNotFound):
		return status.Error(codes.NotFound, "chat not found")
	case errors.Is(err, domain.ErrChatAlreadyExists):
		return status.Error(codes.AlreadyExists, "chat already exists")
	case errors.Is(err, domain.ErrPermissionDenied):
		return status.Error(codes.PermissionDenied, "permission denied")
	case errors.Is(err, domain.ErrInvalidMessage):
		return status.Error(codes.InvalidArgument, "invalid message")
	case errors.Is(err, domain.ErrInvalidChatID):
		return status.Error(codes.InvalidArgument, "invalid chat ID")
	case errors.Is(err, domain.ErrUsersNotFriends):
		return status.Error(codes.PermissionDenied, "users are not friends")
	case errors.Is(err, domain.ErrUserBlocked):
		return status.Error(codes.PermissionDenied, "user is blocked")
	default:
		// Log internal error details here if needed
		// For now, return a generic internal error
		return status.Error(codes.Internal, "internal server error")
	}
}

