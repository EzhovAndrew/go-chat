package handler

import (
	"github.com/go-chat/auth/internal/service"
	authv1 "github.com/go-chat/auth/pkg/api/auth/v1"
)

// Server implements the AuthService gRPC server
type Server struct {
	authv1.UnimplementedAuthServiceServer
	authService  service.AuthService
	tokenService service.TokenService
}

// NewServer creates a new auth service server with injected dependencies
func NewServer(authService service.AuthService, tokenService service.TokenService) *Server {
	return &Server{
		authService:  authService,
		tokenService: tokenService,
	}
}
