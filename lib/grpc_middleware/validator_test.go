package grpc_middleware

import (
	"context"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestNewValidationMiddleware(t *testing.T) {
	middleware, err := NewValidationMiddleware()
	if err != nil {
		t.Fatalf("NewValidationMiddleware() failed: %v", err)
	}

	if middleware == nil {
		t.Fatal("NewValidationMiddleware() returned nil")
	}

	if middleware.validator == nil {
		t.Error("validator should not be nil")
	}
}

func TestValidationMiddleware_UnaryServerInterceptor_ValidMessage(t *testing.T) {
	middleware, err := NewValidationMiddleware()
	if err != nil {
		t.Fatalf("NewValidationMiddleware() failed: %v", err)
	}

	interceptor := middleware.UnaryServerInterceptor()

	// Use a valid empty message (no validation constraints)
	req := &emptypb.Empty{}

	handlerCalled := false
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		handlerCalled = true
		return &emptypb.Empty{}, nil
	}

	info := &grpc.UnaryServerInfo{
		FullMethod: "/test.Service/Method",
	}

	resp, err := interceptor(context.Background(), req, info, handler)
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	if !handlerCalled {
		t.Error("Handler should have been called for valid message")
	}

	if resp == nil {
		t.Error("Response should not be nil")
	}
}

func TestValidationMiddleware_UnaryServerInterceptor_NilMessage(t *testing.T) {
	middleware, err := NewValidationMiddleware()
	if err != nil {
		t.Fatalf("NewValidationMiddleware() failed: %v", err)
	}

	interceptor := middleware.UnaryServerInterceptor()

	handlerCalled := false
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		handlerCalled = true
		return &emptypb.Empty{}, nil
	}

	info := &grpc.UnaryServerInfo{
		FullMethod: "/test.Service/Method",
	}

	// Nil message should cause validation error
	_, err = interceptor(context.Background(), nil, info, handler)
	if err == nil {
		t.Error("Expected error for nil message")
	}

	if handlerCalled {
		t.Error("Handler should not be called when validation fails")
	}

	// Check error code - nil message is an internal error, not invalid argument
	st, ok := status.FromError(err)
	if !ok {
		t.Error("Error should be a gRPC status error")
	} else if st.Code() != codes.Internal {
		t.Errorf("Expected Internal error code for non-proto message, got: %v", st.Code())
	}
}

func TestValidationMiddleware_UnaryServerInterceptor_HandlerError(t *testing.T) {
	middleware, err := NewValidationMiddleware()
	if err != nil {
		t.Fatalf("NewValidationMiddleware() failed: %v", err)
	}

	interceptor := middleware.UnaryServerInterceptor()

	req := &emptypb.Empty{}

	expectedErr := status.Error(codes.Internal, "handler error")
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, expectedErr
	}

	info := &grpc.UnaryServerInfo{
		FullMethod: "/test.Service/Method",
	}

	_, err = interceptor(context.Background(), req, info, handler)
	if err != expectedErr {
		t.Errorf("Expected handler error to be returned, got: %v", err)
	}
}

// mockServerStream is a mock implementation of grpc.ServerStream for testing
type mockServerStream struct {
	grpc.ServerStream
	recvMsgFunc func(m interface{}) error
	sendMsgFunc func(m interface{}) error
}

func (m *mockServerStream) RecvMsg(msg interface{}) error {
	if m.recvMsgFunc != nil {
		return m.recvMsgFunc(msg)
	}
	return nil
}

func (m *mockServerStream) SendMsg(msg interface{}) error {
	if m.sendMsgFunc != nil {
		return m.sendMsgFunc(msg)
	}
	return nil
}

func TestValidatingServerStream_RecvMsg_ValidMessage(t *testing.T) {
	middleware, err := NewValidationMiddleware()
	if err != nil {
		t.Fatalf("NewValidationMiddleware() failed: %v", err)
	}

	msg := &emptypb.Empty{}
	mockStream := &mockServerStream{
		recvMsgFunc: func(m interface{}) error {
			// Copy the empty message
			if empty, ok := m.(*emptypb.Empty); ok {
				proto.Merge(empty, msg)
			}
			return nil
		},
	}

	wrapped := &validatingServerStream{
		ServerStream: mockStream,
		validator:    middleware.validator,
	}

	err = wrapped.RecvMsg(&emptypb.Empty{})
	if err != nil {
		t.Errorf("Expected no error for valid message, got: %v", err)
	}
}

func TestValidatingServerStream_RecvMsg_NilMessage(t *testing.T) {
	middleware, err := NewValidationMiddleware()
	if err != nil {
		t.Fatalf("NewValidationMiddleware() failed: %v", err)
	}

	mockStream := &mockServerStream{
		recvMsgFunc: func(m interface{}) error {
			// Simulate receiving a message
			return nil
		},
	}

	wrapped := &validatingServerStream{
		ServerStream: mockStream,
		validator:    middleware.validator,
	}

	err = wrapped.RecvMsg(nil)
	if err == nil {
		t.Error("Expected error for nil message")
	}

	// Check error code - nil message is an internal error, not invalid argument
	st, ok := status.FromError(err)
	if !ok {
		t.Error("Error should be a gRPC status error")
	} else if st.Code() != codes.Internal {
		t.Errorf("Expected Internal error code for non-proto message, got: %v", st.Code())
	}
}

