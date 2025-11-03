package handler

import (
	"context"
	"log"

	authv1 "github.com/go-chat/auth/pkg/api/auth/v1"
)

func (s *Server) Register(ctx context.Context, req *authv1.RegisterRequest) (*authv1.RegisterResponse, error) {
	log.Println("Register called")

	// TODO: Implement registration logic:
	// - Validate email and password
	// - Hash password with bcrypt
	// - Store user in database
	// - Return ALREADY_EXISTS if email exists

	return &authv1.RegisterResponse{
		UserId: "550e8400-e29b-41d4-a716-446655440000",
	}, nil
}
