package handler

import (
	"context"
	"log"

	authv1 "github.com/go-chat/auth/pkg/api/auth/v1"
)

func (s *Server) Refresh(ctx context.Context, req *authv1.RefreshRequest) (*authv1.RefreshResponse, error) {
	log.Println("Refresh called")

	// TODO: Implement token refresh with rotation:
	// - Validate refresh token (signature, expiry, not revoked)
	// - Generate new access and refresh tokens
	// - Revoke old refresh token
	// - Return UNAUTHENTICATED if invalid

	return &authv1.RefreshResponse{
		AccessToken:  "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.dummy.new_access_token",
		RefreshToken: "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.dummy.new_refresh_token",
		UserId:       "550e8400-e29b-41d4-a716-446655440000",
	}, nil
}

