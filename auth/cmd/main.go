package main

import (
	"log"
	"net"

	"github.com/go-chat/auth/internal/handler"
	authv1 "github.com/go-chat/auth/pkg/api/auth/v1"
	"github.com/go-chat/lib/grpc_middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":8080"
)

func main() {
	log.Println("Auth Service starting...")

	// Create middleware manager with validation enabled by default
	mgr, err := grpc_middleware.NewManager()
	if err != nil {
		log.Fatalf("Failed to create middleware manager: %v", err)
	}

	// Get interceptor chains
	unaryInterceptors, err := mgr.UnaryInterceptors()
	if err != nil {
		log.Fatalf("Failed to get unary interceptors: %v", err)
	}

	streamInterceptors, err := mgr.StreamInterceptors()
	if err != nil {
		log.Fatalf("Failed to get stream interceptors: %v", err)
	}

	// Create gRPC server with middleware
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(unaryInterceptors...),
		grpc.ChainStreamInterceptor(streamInterceptors...),
	)
	authHandler := handler.NewServer()
	authv1.RegisterAuthServiceServer(grpcServer, authHandler)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen on %s: %v", port, err)
	}

	log.Printf("Auth Service listening on %s", port)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
