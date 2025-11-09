package dto

import (
	"github.com/go-chat/social/internal/domain"
	socialv1 "github.com/go-chat/social/pkg/api/social/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ToProtoFriendRequest converts domain.FriendRequest to proto FriendRequest
func ToProtoFriendRequest(fr *domain.FriendRequest) *socialv1.FriendRequest {
	return &socialv1.FriendRequest{
		RequestId:   fr.RequestID.String(),
		RequesterId: fr.RequesterID.String(),
		TargetId:    fr.TargetID.String(),
		CreatedAt:   timestamppb.New(fr.CreatedAt),
		Status:      string(fr.Status),
	}
}

// ToProtoFriendRequests converts a slice of domain.FriendRequest to proto FriendRequest slice
func ToProtoFriendRequests(requests []*domain.FriendRequest) []*socialv1.FriendRequest {
	result := make([]*socialv1.FriendRequest, len(requests))
	for i, fr := range requests {
		result[i] = ToProtoFriendRequest(fr)
	}
	return result
}

