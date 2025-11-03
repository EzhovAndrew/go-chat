package handler

import (
	authv1 "github.com/go-chat/auth/pkg/api/auth/v1"
)

type Server struct {
	authv1.UnimplementedAuthServiceServer
}

func NewServer() *Server {
	return &Server{}
}

