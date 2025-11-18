package grpc

import (
	"context"
	"errors"
	"testing"

	"github.com/go-chat/notifications/internal/domain"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestErrorMapperInterceptor_MapsNotificationNotFound(t *testing.T) {
	interceptor := ErrorMapperInterceptor()
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, domain.ErrNotificationNotFound
	}

	_, err := interceptor(context.Background(), nil, &grpc.UnaryServerInfo{}, handler)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	st, ok := status.FromError(err)
	if !ok {
		t.Fatalf("Expected gRPC status error, got: %v", err)
	}

	if st.Code() != codes.NotFound {
		t.Errorf("Expected code NotFound, got: %v", st.Code())
	}
}

func TestErrorMapperInterceptor_MapsPermissionDenied(t *testing.T) {
	interceptor := ErrorMapperInterceptor()
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, domain.ErrPermissionDenied
	}

	_, err := interceptor(context.Background(), nil, &grpc.UnaryServerInfo{}, handler)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	st, ok := status.FromError(err)
	if !ok {
		t.Fatalf("Expected gRPC status error, got: %v", err)
	}

	if st.Code() != codes.PermissionDenied {
		t.Errorf("Expected code PermissionDenied, got: %v", st.Code())
	}
}

func TestErrorMapperInterceptor_MapsUnknownErrorToInternal(t *testing.T) {
	interceptor := ErrorMapperInterceptor()
	unknownErr := errors.New("unknown error")
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, unknownErr
	}

	_, err := interceptor(context.Background(), nil, &grpc.UnaryServerInfo{}, handler)

	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	st, ok := status.FromError(err)
	if !ok {
		t.Fatalf("Expected gRPC status error, got: %v", err)
	}

	if st.Code() != codes.Internal {
		t.Errorf("Expected code Internal, got: %v", st.Code())
	}
}

func TestErrorMapperInterceptor_PassesThroughSuccessfulResponse(t *testing.T) {
	interceptor := ErrorMapperInterceptor()
	expectedResp := "success"
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return expectedResp, nil
	}

	resp, err := interceptor(context.Background(), nil, &grpc.UnaryServerInfo{}, handler)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if resp != expectedResp {
		t.Errorf("Expected response %v, got: %v", expectedResp, resp)
	}
}

