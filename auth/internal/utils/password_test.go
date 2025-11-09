package utils

import (
	"strings"
	"testing"
)

func TestHashPassword_ValidPassword_ReturnsHash(t *testing.T) {
	password := "SecurePassword123!"

	hash, err := HashPassword(password)

	if err != nil {
		t.Fatalf("HashPassword() returned error: %v", err)
	}

	if hash == "" {
		t.Fatal("HashPassword() returned empty hash")
	}

	// Verify hash format: $argon2id$v=19$m=65536,t=1,p=4$salt$hash
	if !strings.HasPrefix(hash, "$argon2id$") {
		t.Errorf("Hash doesn't have argon2id prefix: %s", hash)
	}

	parts := strings.Split(hash, "$")
	if len(parts) != 6 {
		t.Errorf("Hash has incorrect format, expected 6 parts, got %d", len(parts))
	}
}

func TestHashPassword_SamePasswordTwice_ReturnsDifferentHashes(t *testing.T) {
	password := "SecurePassword123!"

	hash1, err1 := HashPassword(password)
	hash2, err2 := HashPassword(password)

	if err1 != nil || err2 != nil {
		t.Fatalf("HashPassword() returned errors: %v, %v", err1, err2)
	}

	if hash1 == hash2 {
		t.Error("HashPassword() returned identical hashes for same password (should use random salt)")
	}
}

func TestComparePassword_ValidPassword_ReturnsNil(t *testing.T) {
	password := "SecurePassword123!"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	err = ComparePassword(hash, password)

	if err != nil {
		t.Errorf("ComparePassword() returned error for valid password: %v", err)
	}
}

func TestComparePassword_InvalidPassword_ReturnsError(t *testing.T) {
	password := "SecurePassword123!"
	wrongPassword := "WrongPassword456!"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	err = ComparePassword(hash, wrongPassword)

	if err == nil {
		t.Error("ComparePassword() should return error for invalid password")
	}
}

func TestComparePassword_InvalidHashFormat_ReturnsError(t *testing.T) {
	tests := []struct {
		name string
		hash string
	}{
		{"empty hash", ""},
		{"invalid format", "not-a-hash"},
		{"incomplete parts", "$argon2id$v=19$m=65536"},
		{"wrong algorithm", "$bcrypt$v=19$m=65536,t=1,p=4$salt$hash"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ComparePassword(tt.hash, "password")
			if err == nil {
				t.Error("ComparePassword() should return error for invalid hash format")
			}
		})
	}
}

func TestComparePassword_CorruptedHash_ReturnsError(t *testing.T) {
	password := "SecurePassword123!"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("Failed to hash password: %v", err)
	}

	// Corrupt the hash by changing last character
	corruptedHash := hash[:len(hash)-1] + "X"

	err = ComparePassword(corruptedHash, password)

	if err == nil {
		t.Error("ComparePassword() should return error for corrupted hash")
	}
}

func TestHashPassword_EmptyPassword_ReturnsHash(t *testing.T) {
	// Even empty passwords should be hashable
	hash, err := HashPassword("")

	if err != nil {
		t.Fatalf("HashPassword() returned error for empty password: %v", err)
	}

	if hash == "" {
		t.Fatal("HashPassword() returned empty hash")
	}

	// Verify we can still compare
	err = ComparePassword(hash, "")
	if err != nil {
		t.Errorf("ComparePassword() failed for empty password: %v", err)
	}
}

func TestHashPassword_LongPassword_ReturnsHash(t *testing.T) {
	// Test with a very long password
	longPassword := strings.Repeat("a", 1000)

	hash, err := HashPassword(longPassword)

	if err != nil {
		t.Fatalf("HashPassword() returned error for long password: %v", err)
	}

	err = ComparePassword(hash, longPassword)
	if err != nil {
		t.Errorf("ComparePassword() failed for long password: %v", err)
	}
}

func TestComparePassword_ExcessiveTimeParameter_ReturnsError(t *testing.T) {
	// Create a hash with excessive time parameter (DoS attack attempt)
	maliciousHash := "$argon2id$v=19$m=65536,t=999,p=4$c29tZXNhbHQxMjM0NTY=$aGFzaGhlcmVoYXNoaGVyZWhhc2hoZXJl"

	err := ComparePassword(maliciousHash, "password")

	if err == nil {
		t.Error("ComparePassword() should reject hash with excessive time parameter")
	}

	if !strings.Contains(err.Error(), "time parameter exceeds maximum") {
		t.Errorf("Expected error about time parameter, got: %v", err)
	}
}

func TestComparePassword_ExcessiveMemoryParameter_ReturnsError(t *testing.T) {
	// Create a hash with excessive memory parameter (DoS attack attempt)
	maliciousHash := "$argon2id$v=19$m=999999999,t=2,p=4$c29tZXNhbHQxMjM0NTY=$aGFzaGhlcmVoYXNoaGVyZWhhc2hoZXJl"

	err := ComparePassword(maliciousHash, "password")

	if err == nil {
		t.Error("ComparePassword() should reject hash with excessive memory parameter")
	}

	if !strings.Contains(err.Error(), "memory parameter exceeds maximum") {
		t.Errorf("Expected error about memory parameter, got: %v", err)
	}
}

func TestComparePassword_ExcessiveThreadsParameter_ReturnsError(t *testing.T) {
	// Create a hash with excessive threads parameter (DoS attack attempt)
	// Note: uint8 max is 255, so we test with 250 which is within parsing range but above our limit
	maliciousHash := "$argon2id$v=19$m=65536,t=2,p=250$c29tZXNhbHQxMjM0NTY=$aGFzaGhlcmVoYXNoaGVyZWhhc2hoZXJl"

	err := ComparePassword(maliciousHash, "password")

	if err == nil {
		t.Error("ComparePassword() should reject hash with excessive parallelism parameter")
	}

	if !strings.Contains(err.Error(), "parallelism parameter exceeds maximum") {
		t.Errorf("Expected error about parallelism parameter, got: %v", err)
	}
}

func TestComparePassword_UnsupportedVersion_ReturnsError(t *testing.T) {
	// Create a hash with unsupported version
	maliciousHash := "$argon2id$v=99$m=65536,t=2,p=4$c29tZXNhbHQxMjM0NTY=$aGFzaGhlcmVoYXNoaGVyZWhhc2hoZXJl"

	err := ComparePassword(maliciousHash, "password")

	if err == nil {
		t.Error("ComparePassword() should reject hash with unsupported version")
	}

	if !strings.Contains(err.Error(), "unsupported argon2 version") {
		t.Errorf("Expected error about unsupported version, got: %v", err)
	}
}
