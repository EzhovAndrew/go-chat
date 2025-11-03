package main

import (
	"log"
	"net"

	"github.com/go-chat/social/internal/handler"
	socialv1 "github.com/go-chat/social/pkg/api/social/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":8080"
)

func main() {
	log.Println("Social Service starting...")

	grpcServer := grpc.NewServer()
	socialHandler := handler.NewServer()
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

