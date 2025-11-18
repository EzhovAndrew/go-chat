package grpc

import (
	"context"
	"errors"

	"github.com/go-chat/social/internal/domain"
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
	case errors.Is(err, domain.ErrRequestNotFound):
		return status.Error(codes.NotFound, "friend request not found")
	case errors.Is(err, domain.ErrRequestAlreadyExists):
		return status.Error(codes.AlreadyExists, "friend request already exists")
	case errors.Is(err, domain.ErrAlreadyFriends):
		return status.Error(codes.AlreadyExists, "users are already friends")
	case errors.Is(err, domain.ErrNotFriends):
		return status.Error(codes.NotFound, "users are not friends")
	case errors.Is(err, domain.ErrUserBlocked):
		return status.Error(codes.PermissionDenied, "user is blocked")
	case errors.Is(err, domain.ErrPermissionDenied):
		return status.Error(codes.PermissionDenied, "permission denied")
	case errors.Is(err, domain.ErrSelfAction):
		return status.Error(codes.InvalidArgument, "cannot perform this action on yourself")
	default:
		// Log internal error details here if needed
		// For now, return a generic internal error
		return status.Error(codes.Internal, "internal server error")
	}
}

