package handler

import (
	"github.com/go-chat/users/internal/service"
	usersv1 "github.com/go-chat/users/pkg/api/users/v1"
)

// Server implements the UserService gRPC interface
type Server struct {
	usersv1.UnimplementedUserServiceServer
	userService service.UserService
}

// NewServer creates a new User service handler with injected dependencies
func NewServer(userService service.UserService) *Server {
	return &Server{
		userService: userService,
	}
}

