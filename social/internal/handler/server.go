package handler

import (
	"github.com/go-chat/social/internal/service"
	socialv1 "github.com/go-chat/social/pkg/api/social/v1"
)

// Server implements the SocialService gRPC interface
type Server struct {
	socialv1.UnimplementedSocialServiceServer
	friendRequestService service.FriendRequestService
	friendshipService    service.FriendshipService
	blockService         service.BlockService
	relationshipService  service.RelationshipService
}

// NewServer creates a new Social service handler with injected dependencies
func NewServer(
	friendRequestService service.FriendRequestService,
	friendshipService service.FriendshipService,
	blockService service.BlockService,
	relationshipService service.RelationshipService,
) *Server {
	return &Server{
		friendRequestService: friendRequestService,
		friendshipService:    friendshipService,
		blockService:         blockService,
		relationshipService:  relationshipService,
	}
}
