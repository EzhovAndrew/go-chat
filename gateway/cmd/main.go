package main

import (
	"context"
	"log"

	"github.com/go-chat/gateway/internal/config"
	"github.com/go-chat/gateway/internal/server"
)

func main() {
	ctx := context.Background()
	cfg := config.New()
	srv := server.New(cfg)

	if err := srv.Start(ctx); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
