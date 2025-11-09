package handler

import (
	"context"

	authv1 "github.com/go-chat/auth/pkg/api/auth/v1"
)

// Refresh renews the access token using a refresh token
func (s *Server) Refresh(ctx context.Context, req *authv1.RefreshRequest) (*authv1.RefreshResponse, error) {
	tokens, userID, err := s.authService.Refresh(ctx, req.RefreshToken)
	if err != nil {
		return nil, err // Middleware will map domain error to gRPC status
	}

	return &authv1.RefreshResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		UserId:       userID.String(),
	}, nil
}

