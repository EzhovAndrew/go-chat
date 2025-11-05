package grpc_middleware

import (
	"context"
	"fmt"

	"buf.build/go/protovalidate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

// ValidationMiddleware validates incoming gRPC requests against proto validation rules.
// It uses protovalidate-go to enforce constraints defined in proto files using buf.validate annotations.
type ValidationMiddleware struct {
	validator protovalidate.Validator
}

// NewValidationMiddleware creates a new validation middleware instance.
// It initializes the protovalidate validator that will be reused for all requests.
func NewValidationMiddleware() (*ValidationMiddleware, error) {
	validator, err := protovalidate.New()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize protovalidate: %w", err)
	}
	return &ValidationMiddleware{validator: validator}, nil
}

// UnaryServerInterceptor returns a unary server interceptor for request validation.
// It validates the request message before passing it to the actual handler.
// If validation fails, it returns an InvalidArgument error without calling the handler.
func (v *ValidationMiddleware) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		// Type assert to proto.Message for validation
		msg, ok := req.(proto.Message)
		if !ok {
			return nil, status.Error(codes.Internal, "request is not a proto message")
		}

		// Validate request message against proto validation rules
		if err := v.validator.Validate(msg); err != nil {
			return nil, formatValidationError(err)
		}

		// Validation passed, call actual handler
		return handler(ctx, req)
	}
}

// StreamServerInterceptor returns a stream server interceptor for message validation.
// It wraps the server stream to validate each message received from the client.
func (v *ValidationMiddleware) StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		// Wrap stream to validate each message
		wrapped := &validatingServerStream{
			ServerStream: ss,
			validator:    v.validator,
		}
		return handler(srv, wrapped)
	}
}

// validatingServerStream wraps grpc.ServerStream to validate received messages.
type validatingServerStream struct {
	grpc.ServerStream
	validator protovalidate.Validator
}

// RecvMsg validates messages received from the client.
// It first receives the message, then validates it before returning to the handler.
func (s *validatingServerStream) RecvMsg(m interface{}) error {
	// Receive message from client
	if err := s.ServerStream.RecvMsg(m); err != nil {
		return err
	}

	// Type assert to proto.Message for validation
	msg, ok := m.(proto.Message)
	if !ok {
		return status.Error(codes.Internal, "message is not a proto message")
	}

	// Validate the received message
	if err := s.validator.Validate(msg); err != nil {
		return formatValidationError(err)
	}

	return nil
}
