package main

import (
	"log"
	"net"

	"github.com/go-chat/chat/internal/handler"
	chatv1 "github.com/go-chat/chat/pkg/api/chat/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":8080"
)

func main() {
	log.Println("Chat Service starting...")

	grpcServer := grpc.NewServer()
	chatHandler := handler.NewServer()
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

