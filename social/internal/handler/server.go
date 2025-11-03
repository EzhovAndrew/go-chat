package handler

import (
	socialv1 "github.com/go-chat/social/pkg/api/social/v1"
)

// Server implements the SocialService gRPC interface
type Server struct {
	socialv1.UnimplementedSocialServiceServer
}

// NewServer creates a new Social service handler
// Dependencies (repository, kafka producer) will be added as needed
func NewServer() *Server {
	return &Server{}
}

