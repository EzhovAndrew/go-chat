package handler

import (
	"context"

	authv1 "github.com/go-chat/auth/pkg/api/auth/v1"
)

// Register creates a new user account
func (s *Server) Register(ctx context.Context, req *authv1.RegisterRequest) (*authv1.RegisterResponse, error) {
	user, err := s.authService.Register(ctx, req.Email, req.Password)
	if err != nil {
		return nil, err // Middleware will map domain error to gRPC status
	}

	return &authv1.RegisterResponse{
		UserId: user.ID.String(),
	}, nil
}
