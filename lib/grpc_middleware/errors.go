package grpc_middleware

import (
	"fmt"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// formatValidationError converts a protovalidate error into a gRPC status error.
// It returns an InvalidArgument error code with a descriptive message containing
// all validation violations.
func formatValidationError(err error) error {
	if err == nil {
		return nil
	}

	// Format the error message to be client-friendly
	msg := formatErrorMessage(err.Error())

	return status.Error(codes.InvalidArgument, msg)
}

// formatErrorMessage formats the validation error message for better readability.
// It preserves the detailed validation information while ensuring consistent formatting.
func formatErrorMessage(msg string) string {
	// Remove excessive whitespace and normalize the message
	msg = strings.TrimSpace(msg)

	// If message is empty, provide a generic message
	if msg == "" {
		return "validation failed"
	}

	// Ensure message starts with lowercase for consistency (unless it's a field name)
	if len(msg) > 0 && msg[0] >= 'A' && msg[0] <= 'Z' {
		// Check if it looks like a field path (contains dot or starts with capital)
		if !strings.Contains(msg, ".") {
			msg = strings.ToLower(msg[:1]) + msg[1:]
		}
	}

	return fmt.Sprintf("validation failed: %s", msg)
}
