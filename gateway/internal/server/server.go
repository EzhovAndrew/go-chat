package server

import (
	"context"
	"log"
	"net/http"

	"github.com/go-chat/gateway/internal/config"
	"github.com/go-chat/gateway/internal/middleware"
	"github.com/go-chat/gateway/internal/proxy"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
)

// Server represents the Gateway HTTP server
type Server struct {
	cfg        *config.Config
	httpServer *http.Server
}

// New creates a new Gateway server
func New(cfg *config.Config) *Server {
	return &Server{
		cfg: cfg,
	}
}

func (s *Server) Start(ctx context.Context) error {
	log.Println("Gateway starting...")

	grpcMux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{}),
	)

	if err := proxy.RegisterServices(ctx, grpcMux, s.cfg); err != nil {
		return err
	}

	handler := middleware.CORS(grpcMux)

	s.httpServer = &http.Server{
		Addr:    s.cfg.HTTPPort,
		Handler: handler,
	}

	log.Printf("Gateway listening on %s", s.cfg.HTTPPort)
	log.Printf("Gateway ready to proxy requests to backend services")

	return s.httpServer.ListenAndServe()
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	if s.httpServer != nil {
		return s.httpServer.Shutdown(ctx)
	}
	return nil
}
