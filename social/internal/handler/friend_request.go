package handler

import (
	"context"

	"github.com/go-chat/social/internal/domain"
	"github.com/go-chat/social/internal/dto"
	socialv1 "github.com/go-chat/social/pkg/api/social/v1"
)

// SendFriendRequest sends a friend request to another user
func (s *Server) SendFriendRequest(ctx context.Context, req *socialv1.SendFriendRequestRequest) (*socialv1.SendFriendRequestResponse, error) {
	// TODO: Extract authenticated user_id from JWT
	requesterID := domain.NewUserID("authenticated-user-id") // Placeholder

	// Delegate to service layer
	friendRequest, err := s.friendRequestService.SendFriendRequest(
		ctx,
		requesterID,
		domain.NewUserID(req.TargetUserId),
	)
	if err != nil {
		return nil, err // Middleware will map domain error to gRPC status
	}

	// Convert domain model to proto message
	return &socialv1.SendFriendRequestResponse{
		Request: dto.ToProtoFriendRequest(friendRequest),
	}, nil
}

// ListRequests lists pending friend requests for a user
func (s *Server) ListRequests(ctx context.Context, req *socialv1.ListRequestsRequest) (*socialv1.ListRequestsResponse, error) {
	// Default limit if not provided
	limit := req.Limit
	if limit == 0 {
		limit = 20
	}

	// Delegate to service layer
	requests, nextCursor, err := s.friendRequestService.ListRequests(
		ctx,
		domain.NewUserID(req.UserId),
		req.Cursor,
		limit,
	)
	if err != nil {
		return nil, err // Middleware will map domain error to gRPC status
	}

	// Convert domain models to proto messages
	return &socialv1.ListRequestsResponse{
		Requests:   dto.ToProtoFriendRequests(requests),
		NextCursor: nextCursor,
	}, nil
}

// AcceptFriendRequest accepts a pending friend request
func (s *Server) AcceptFriendRequest(ctx context.Context, req *socialv1.AcceptFriendRequestRequest) (*socialv1.AcceptFriendRequestResponse, error) {
	// TODO: Extract authenticated user_id from JWT
	userID := domain.NewUserID("authenticated-user-id") // Placeholder

	// Delegate to service layer
	friendRequest, err := s.friendRequestService.AcceptFriendRequest(
		ctx,
		domain.NewRequestID(req.RequestId),
		userID,
	)
	if err != nil {
		return nil, err // Middleware will map domain error to gRPC status
	}

	// Convert domain model to proto message
	return &socialv1.AcceptFriendRequestResponse{
		Request: dto.ToProtoFriendRequest(friendRequest),
	}, nil
}

// DeclineFriendRequest declines a pending friend request
func (s *Server) DeclineFriendRequest(ctx context.Context, req *socialv1.DeclineFriendRequestRequest) (*socialv1.DeclineFriendRequestResponse, error) {
	// TODO: Extract authenticated user_id from JWT
	userID := domain.NewUserID("authenticated-user-id") // Placeholder

	// Delegate to service layer
	err := s.friendRequestService.DeclineFriendRequest(
		ctx,
		domain.NewRequestID(req.RequestId),
		userID,
	)
	if err != nil {
		return nil, err // Middleware will map domain error to gRPC status
	}

	// Return empty response on success (as per proto)
	return &socialv1.DeclineFriendRequestResponse{
		Request: &socialv1.FriendRequest{
			RequestId: req.RequestId,
			Status:    string(domain.FriendRequestStatusDeclined),
		},
	}, nil
}
