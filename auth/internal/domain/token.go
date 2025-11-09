package domain

import "time"

// TokenPair represents a pair of access and refresh tokens
type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

// AccessTokenClaims represents the claims contained in an access token
type AccessTokenClaims struct {
	UserID    string
	Email     string
	IssuedAt  time.Time
	ExpiresAt time.Time
}

// PublicKey represents a JSON Web Key for JWT validation
type PublicKey struct {
	Kid string // Key ID
	Alg string // Algorithm (RS256)
	Use string // Usage (sig)
	N   string // RSA modulus (base64url)
	E   string // RSA exponent (base64url)
}

