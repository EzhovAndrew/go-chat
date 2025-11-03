package handler

import (
	usersv1 "github.com/go-chat/users/pkg/api/users/v1"
)

// Server implements the UserService gRPC interface
type Server struct {
	usersv1.UnimplementedUserServiceServer
}

// NewServer creates a new User service handler
// Dependencies (repository, validator) will be added here as we implement real functionality
func NewServer() *Server {
	return &Server{}
}

