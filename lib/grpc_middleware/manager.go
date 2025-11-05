package grpc_middleware

import (
	"fmt"

	"google.golang.org/grpc"
)

// Manager orchestrates gRPC middleware in the correct order.
// It provides a centralized way to configure and apply middleware across all services.
type Manager struct {
	config *Config
}

// Config holds middleware configuration options.
type Config struct {
	// ValidationEnabled controls whether request validation middleware is active.
	// When enabled, all incoming requests are validated against proto validation rules.
	ValidationEnabled bool
}

// Option is a functional option for configuring the Manager.
type Option func(*Config)

// NewManager creates a new middleware manager with the given options.
// By default, validation is enabled. Use WithValidation(false) to disable it.
func NewManager(opts ...Option) (*Manager, error) {
	cfg := &Config{
		ValidationEnabled: true, // Enabled by default
	}

	for _, opt := range opts {
		opt(cfg)
	}

	return &Manager{config: cfg}, nil
}

// UnaryInterceptors returns all unary server interceptors in the correct order.
// Interceptors are ordered from outermost to innermost in the execution chain.
func (m *Manager) UnaryInterceptors() ([]grpc.UnaryServerInterceptor, error) {
	var interceptors []grpc.UnaryServerInterceptor

	// Validation middleware - validate requests before handler execution
	if m.config.ValidationEnabled {
		validator, err := NewValidationMiddleware()
		if err != nil {
			return nil, fmt.Errorf("failed to create validation middleware: %w", err)
		}
		interceptors = append(interceptors, validator.UnaryServerInterceptor())
	}

	return interceptors, nil
}

// StreamInterceptors returns all stream server interceptors in the correct order.
// Interceptors are ordered from outermost to innermost in the execution chain.
func (m *Manager) StreamInterceptors() ([]grpc.StreamServerInterceptor, error) {
	var interceptors []grpc.StreamServerInterceptor

	// Validation middleware - validate stream messages
	if m.config.ValidationEnabled {
		validator, err := NewValidationMiddleware()
		if err != nil {
			return nil, fmt.Errorf("failed to create validation middleware: %w", err)
		}
		interceptors = append(interceptors, validator.StreamServerInterceptor())
	}

	return interceptors, nil
}

// WithValidation enables or disables validation middleware.
// When disabled, no request validation will be performed.
func WithValidation(enabled bool) Option {
	return func(c *Config) {
		c.ValidationEnabled = enabled
	}
}

