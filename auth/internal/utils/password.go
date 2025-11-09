package utils

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

const (
	// Argon2id parameters (OWASP recommended minimums for 2024+)
	// OWASP minimum: m=19456 KiB, t=2, p=1
	argon2Time       = 2         // Number of iterations (OWASP minimum: 2)
	argon2Memory     = 64 * 1024 // Memory in KiB (64 MB, OWASP minimum: 19 MiB)
	argon2Threads    = 4         // Number of threads (parallelism, OWASP minimum: 1)
	argon2KeyLength  = 32        // Length of the derived key (256 bits)
	argon2SaltLength = 16        // Length of the random salt (128 bits)

	// Maximum allowed parameters to prevent DoS attacks
	maxArgon2Time    = 10         // Max iterations
	maxArgon2Memory  = 256 * 1024 // Max 256 MB
	maxArgon2Threads = 16         // Max parallelism
)

// HashPassword hashes a plaintext password using Argon2id
func HashPassword(password string) (string, error) {
	// Generate a cryptographically secure random salt
	salt := make([]byte, argon2SaltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}

	// Hash the password using Argon2id
	hash := argon2.IDKey([]byte(password), salt, argon2Time, argon2Memory, argon2Threads, argon2KeyLength)

	// Encode the hash in a standard format: $argon2id$v=19$m=65536,t=1,p=4$salt$hash
	encodedSalt := base64.RawStdEncoding.EncodeToString(salt)
	encodedHash := base64.RawStdEncoding.EncodeToString(hash)

	return fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version, argon2Memory, argon2Time, argon2Threads, encodedSalt, encodedHash), nil
}

// ComparePassword compares a hashed password with a plaintext password
func ComparePassword(encodedHash, password string) error {
	// Extract parameters and hash from the encoded string
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 {
		return fmt.Errorf("invalid hash format")
	}

	if parts[1] != "argon2id" {
		return fmt.Errorf("incompatible hash algorithm")
	}

	var version int
	if _, err := fmt.Sscanf(parts[2], "v=%d", &version); err != nil {
		return fmt.Errorf("invalid version: %w", err)
	}

	var memory, time uint32
	var threads uint8
	if _, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &time, &threads); err != nil {
		return fmt.Errorf("invalid parameters: %w", err)
	}

	// Validate parameters to prevent DoS attacks
	if time > maxArgon2Time {
		return fmt.Errorf("time parameter exceeds maximum allowed value")
	}
	if memory > maxArgon2Memory {
		return fmt.Errorf("memory parameter exceeds maximum allowed value")
	}
	if threads > maxArgon2Threads {
		return fmt.Errorf("parallelism parameter exceeds maximum allowed value")
	}

	// Validate version (should be 19 for Argon2 v1.3)
	if version != 19 {
		return fmt.Errorf("unsupported argon2 version: %d", version)
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return fmt.Errorf("invalid salt encoding: %w", err)
	}

	hash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return fmt.Errorf("invalid hash encoding: %w", err)
	}

	// Hash the input password with the same parameters
	passwordHash := argon2.IDKey([]byte(password), salt, time, memory, threads, uint32(len(hash)))

	// Compare in constant time to prevent timing attacks
	if subtle.ConstantTimeCompare(hash, passwordHash) == 1 {
		return nil
	}

	return fmt.Errorf("password does not match")
}
