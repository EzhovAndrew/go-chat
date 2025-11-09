package handler

import (
	"context"

	authv1 "github.com/go-chat/auth/pkg/api/auth/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetPublicKeys returns public keys for JWT validation (internal endpoint)
func (s *Server) GetPublicKeys(ctx context.Context, req *authv1.GetPublicKeysRequest) (*authv1.GetPublicKeysResponse, error) {
	keys, err := s.tokenService.GetPublicKeys(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to retrieve public keys")
	}

	// Convert domain keys to protobuf
	pbKeys := make([]*authv1.PublicKey, len(keys))
	for i, key := range keys {
		pbKeys[i] = &authv1.PublicKey{
			Kid: key.Kid,
			Alg: key.Alg,
			Use: key.Use,
			N:   key.N,
			E:   key.E,
		}
	}

	return &authv1.GetPublicKeysResponse{Keys: pbKeys}, nil
}
