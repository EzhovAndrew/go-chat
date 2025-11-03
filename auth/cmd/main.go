package main

import (
	"log"
	"net"

	"github.com/go-chat/auth/internal/handler"
	authv1 "github.com/go-chat/auth/pkg/api/auth/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":8080"
)

func main() {
	log.Println("Auth Service starting...")

	grpcServer := grpc.NewServer()
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

