package handler

import (
	"context"
	"log"

	authv1 "github.com/go-chat/auth/pkg/api/auth/v1"
)

func (s *Server) GetPublicKeys(ctx context.Context, req *authv1.GetPublicKeysRequest) (*authv1.GetPublicKeysResponse, error) {
	log.Println("GetPublicKeys called")

	// TODO: Implement JWK endpoint:
	// - Load RSA key pair
	// - Convert public key to JWK format
	// - Support key rotation

	return &authv1.GetPublicKeysResponse{
		Keys: []*authv1.PublicKey{
			{
				Kid: "key-2024-01",
				Alg: "RS256",
				Use: "sig",
				N:   "dummy_modulus_base64url_encoded_value",
				E:   "AQAB",
			},
		},
	}, nil
}
