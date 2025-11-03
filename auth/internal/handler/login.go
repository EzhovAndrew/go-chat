package handler

import (
	"context"
	"log"

	authv1 "github.com/go-chat/auth/pkg/api/auth/v1"
)

func (s *Server) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	log.Println("Login called")

	// TODO: Implement authentication logic:
	// - Query user by email
	// - Verify password with bcrypt
	// - Generate JWT tokens (RS256): access (15min), refresh (30d)
	// - Store refresh token in database
	// - Return UNAUTHENTICATED if invalid

	return &authv1.LoginResponse{
		AccessToken:  "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.dummy.access_token",
		RefreshToken: "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.dummy.refresh_token",
		UserId:       "550e8400-e29b-41d4-a716-446655440000",
	}, nil
}
