package handler

import (
	"context"

	authv1 "github.com/go-chat/auth/pkg/api/auth/v1"
)

// Login authenticates a user and returns JWT tokens
func (s *Server) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	tokens, userID, err := s.authService.Login(ctx, req.Email, req.Password)
	if err != nil {
		return nil, err // Middleware will map domain error to gRPC status
	}

	return &authv1.LoginResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		UserId:       userID.String(),
	}, nil
}
