package main

import (
	"log"
	"net"

	"github.com/go-chat/lib/grpc_middleware"
	"github.com/go-chat/social/internal/handler"
	grpcmw "github.com/go-chat/social/internal/middleware/grpc"
	"github.com/go-chat/social/internal/service"
	socialv1 "github.com/go-chat/social/pkg/api/social/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":8080"
)

func main() {
	log.Println("Social Service starting...")

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

	// Add social service error mapper middleware
	socialErrorMapper := grpcmw.ErrorMapperInterceptor()
	unaryInterceptors = append(unaryInterceptors, socialErrorMapper)

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
	var friendRequestService service.FriendRequestService = nil
	var friendshipService service.FriendshipService = nil
	var blockService service.BlockService = nil
	var relationshipService service.RelationshipService = nil

	socialHandler := handler.NewServer(friendRequestService, friendshipService, blockService, relationshipService)
	socialv1.RegisterSocialServiceServer(grpcServer, socialHandler)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen on %s: %v", port, err)
	}

	log.Printf("Social Service listening on %s", port)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
