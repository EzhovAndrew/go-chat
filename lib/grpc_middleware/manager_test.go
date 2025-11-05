package grpc_middleware

import (
	"testing"
)

func TestNewManager_DefaultConfig(t *testing.T) {
	mgr, err := NewManager()
	if err != nil {
		t.Fatalf("NewManager() failed: %v", err)
	}

	if mgr == nil {
		t.Fatal("NewManager() returned nil manager")
	}

	if !mgr.config.ValidationEnabled {
		t.Error("Expected validation to be enabled by default")
	}
}

func TestNewManager_WithValidationDisabled(t *testing.T) {
	mgr, err := NewManager(WithValidation(false))
	if err != nil {
		t.Fatalf("NewManager() failed: %v", err)
	}

	if mgr.config.ValidationEnabled {
		t.Error("Expected validation to be disabled")
	}
}

func TestNewManager_WithValidationEnabled(t *testing.T) {
	mgr, err := NewManager(WithValidation(true))
	if err != nil {
		t.Fatalf("NewManager() failed: %v", err)
	}

	if !mgr.config.ValidationEnabled {
		t.Error("Expected validation to be enabled")
	}
}

func TestManager_UnaryInterceptors_ValidationEnabled(t *testing.T) {
	mgr, err := NewManager(WithValidation(true))
	if err != nil {
		t.Fatalf("NewManager() failed: %v", err)
	}

	interceptors, err := mgr.UnaryInterceptors()
	if err != nil {
		t.Fatalf("UnaryInterceptors() failed: %v", err)
	}

	if len(interceptors) != 1 {
		t.Errorf("Expected 1 interceptor, got %d", len(interceptors))
	}
}

func TestManager_UnaryInterceptors_ValidationDisabled(t *testing.T) {
	mgr, err := NewManager(WithValidation(false))
	if err != nil {
		t.Fatalf("NewManager() failed: %v", err)
	}

	interceptors, err := mgr.UnaryInterceptors()
	if err != nil {
		t.Fatalf("UnaryInterceptors() failed: %v", err)
	}

	if len(interceptors) != 0 {
		t.Errorf("Expected 0 interceptors, got %d", len(interceptors))
	}
}

func TestManager_StreamInterceptors_ValidationEnabled(t *testing.T) {
	mgr, err := NewManager(WithValidation(true))
	if err != nil {
		t.Fatalf("NewManager() failed: %v", err)
	}

	interceptors, err := mgr.StreamInterceptors()
	if err != nil {
		t.Fatalf("StreamInterceptors() failed: %v", err)
	}

	if len(interceptors) != 1 {
		t.Errorf("Expected 1 interceptor, got %d", len(interceptors))
	}
}

func TestManager_StreamInterceptors_ValidationDisabled(t *testing.T) {
	mgr, err := NewManager(WithValidation(false))
	if err != nil {
		t.Fatalf("NewManager() failed: %v", err)
	}

	interceptors, err := mgr.StreamInterceptors()
	if err != nil {
		t.Fatalf("StreamInterceptors() failed: %v", err)
	}

	if len(interceptors) != 0 {
		t.Errorf("Expected 0 interceptors, got %d", len(interceptors))
	}
}

func TestManager_MultipleOptions(t *testing.T) {
	// Test that multiple options can be applied
	mgr, err := NewManager(
		WithValidation(false),
		WithValidation(true), // Last one wins
	)
	if err != nil {
		t.Fatalf("NewManager() failed: %v", err)
	}

	if !mgr.config.ValidationEnabled {
		t.Error("Expected validation to be enabled (last option should win)")
	}
}
