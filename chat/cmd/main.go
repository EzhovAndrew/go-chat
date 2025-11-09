package main

import (
	"log"
	"net"

	"github.com/go-chat/chat/internal/handler"
	grpcmw "github.com/go-chat/chat/internal/middleware/grpc"
	"github.com/go-chat/chat/internal/service"
	chatv1 "github.com/go-chat/chat/pkg/api/chat/v1"
	"github.com/go-chat/lib/grpc_middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":8080"
)

func main() {
	log.Println("Chat Service starting...")

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

	// Add chat service error mapper middleware
	chatErrorMapper := grpcmw.ErrorMapperInterceptor()
	unaryInterceptors = append(unaryInterceptors, chatErrorMapper)

	streamInterceptors, err := mgr.StreamInterceptors()
	if err != nil {
		log.Fatalf("Failed to get stream interceptors: %v", err)
	}

	// Create gRPC server with middleware
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(unaryInterceptors...),
		grpc.ChainStreamInterceptor(streamInterceptors...),
	)

	// TODO: Replace with actual service implementations in next iteration
	// For now, use nil services - handlers will panic if called
	var chatService service.ChatService = nil
	var messageService service.MessageService = nil

	chatHandler := handler.NewServer(chatService, messageService)
	chatv1.RegisterChatServiceServer(grpcServer, chatHandler)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen on %s: %v", port, err)
	}

	log.Printf("Chat Service listening on %s", port)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
