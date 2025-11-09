package handler

import (
	"context"

	"github.com/go-chat/social/internal/domain"
	socialv1 "github.com/go-chat/social/pkg/api/social/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
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
		Request: &socialv1.FriendRequest{
			RequestId:   friendRequest.RequestID.String(),
			RequesterId: friendRequest.RequesterID.String(),
			TargetId:    friendRequest.TargetID.String(),
			CreatedAt:   timestamppb.New(friendRequest.CreatedAt),
			Status:      string(friendRequest.Status),
		},
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
	pbRequests := make([]*socialv1.FriendRequest, len(requests))
	for i, request := range requests {
		pbRequests[i] = &socialv1.FriendRequest{
			RequestId:   request.RequestID.String(),
			RequesterId: request.RequesterID.String(),
			TargetId:    request.TargetID.String(),
			CreatedAt:   timestamppb.New(request.CreatedAt),
			Status:      string(request.Status),
		}
	}

	return &socialv1.ListRequestsResponse{
		Requests:   pbRequests,
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
		Request: &socialv1.FriendRequest{
			RequestId:   friendRequest.RequestID.String(),
			RequesterId: friendRequest.RequesterID.String(),
			TargetId:    friendRequest.TargetID.String(),
			CreatedAt:   timestamppb.New(friendRequest.CreatedAt),
			Status:      string(friendRequest.Status),
		},
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
