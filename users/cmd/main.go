package main

import (
	"log"
	"net"

	"github.com/go-chat/users/internal/handler"
	usersv1 "github.com/go-chat/users/pkg/api/users/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":8080"
)

func main() {
	log.Println("Users Service starting...")

	grpcServer := grpc.NewServer()
	userHandler := handler.NewServer()
	usersv1.RegisterUserServiceServer(grpcServer, userHandler)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen on %s: %v", port, err)
	}

	log.Printf("Users Service listening on %s", port)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

