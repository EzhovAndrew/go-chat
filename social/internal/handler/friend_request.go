package handler

import (
	"context"
	"log"

	socialv1 "github.com/go-chat/social/pkg/api/social/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// SendFriendRequest sends a friend request to another user
// Returns dummy request until database and Kafka integration is added
func (s *Server) SendFriendRequest(ctx context.Context, req *socialv1.SendFriendRequestRequest) (*socialv1.SendFriendRequestResponse, error) {
	log.Printf("SendFriendRequest called with target_user_id: %s", req.TargetUserId)

	// TODO: Extract authenticated user_id (requester_id) from JWT
	// TODO: Validate users are not already friends
	// TODO: Check if target has blocked requester (return INVALID_ARGUMENT or PERMISSION_DENIED)
	// TODO: Check if request already exists (idempotent - return existing or ALREADY_EXISTS)
	// TODO: Store friend request in database
	// TODO: Publish friend_request.sent event to Kafka for notifications
	// TODO: Return created request

	return &socialv1.SendFriendRequestResponse{
		Request: &socialv1.FriendRequest{
			RequestId:   "550e8400-e29b-41d4-a716-446655440000",
			RequesterId: "550e8400-e29b-41d4-a716-446655440001",
			TargetId:    req.TargetUserId,
			CreatedAt:   timestamppb.Now(),
			Status:      "pending",
		},
	}, nil
}

// ListRequests lists pending friend requests for a user with cursor-based pagination
// Returns dummy requests until database integration is added
func (s *Server) ListRequests(ctx context.Context, req *socialv1.ListRequestsRequest) (*socialv1.ListRequestsResponse, error) {
	log.Printf("ListRequests called for user_id: %s, limit: %d", req.UserId, req.Limit)

	// TODO: Extract authenticated user_id from JWT
	// TODO: Verify req.UserId matches authenticated user
	// TODO: Decode cursor for pagination
	// TODO: Query pending friend requests from database where target_id = user_id
	// TODO: Order by created_at DESC (newest first)
	// TODO: Apply limit + 1 to check for more results
	// TODO: Encode next_cursor if more results exist

	return &socialv1.ListRequestsResponse{
		Requests: []*socialv1.FriendRequest{
			{
				RequestId:   "550e8400-e29b-41d4-a716-446655440001",
				RequesterId: "550e8400-e29b-41d4-a716-446655440002",
				TargetId:    req.UserId,
				CreatedAt:   timestamppb.Now(),
				Status:      "pending",
			},
		},
		NextCursor: "", // Empty means no more results
	}, nil
}

// AcceptFriendRequest accepts a friend request
// Returns dummy request until database and Kafka integration is added
func (s *Server) AcceptFriendRequest(ctx context.Context, req *socialv1.AcceptFriendRequestRequest) (*socialv1.AcceptFriendRequestResponse, error) {
	log.Printf("AcceptFriendRequest called with request_id: %s", req.RequestId)

	// TODO: Extract authenticated user_id from JWT
	// TODO: Query friend request from database
	// TODO: Return NOT_FOUND if request doesn't exist
	// TODO: Verify authenticated user is the target_id (PERMISSION_DENIED if not)
	// TODO: Update request status to "accepted"
	// TODO: Create bidirectional friendship records in database
	// TODO: Publish friend_request.accepted event to Kafka for notifications
	// TODO: Return updated request

	return &socialv1.AcceptFriendRequestResponse{
		Request: &socialv1.FriendRequest{
			RequestId:   req.RequestId,
			RequesterId: "550e8400-e29b-41d4-a716-446655440001",
			TargetId:    "550e8400-e29b-41d4-a716-446655440002",
			CreatedAt:   timestamppb.Now(),
			Status:      "accepted",
		},
	}, nil
}

// DeclineFriendRequest declines a friend request
// Returns dummy request until database integration is added
func (s *Server) DeclineFriendRequest(ctx context.Context, req *socialv1.DeclineFriendRequestRequest) (*socialv1.DeclineFriendRequestResponse, error) {
	log.Printf("DeclineFriendRequest called with request_id: %s", req.RequestId)

	// TODO: Extract authenticated user_id from JWT
	// TODO: Query friend request from database
	// TODO: Return NOT_FOUND if request doesn't exist
	// TODO: Verify authenticated user is the target_id (PERMISSION_DENIED if not)
	// TODO: Update request status to "declined" or delete it
	// TODO: Return updated/deleted request

	return &socialv1.DeclineFriendRequestResponse{
		Request: &socialv1.FriendRequest{
			RequestId:   req.RequestId,
			RequesterId: "550e8400-e29b-41d4-a716-446655440001",
			TargetId:    "550e8400-e29b-41d4-a716-446655440002",
			CreatedAt:   timestamppb.Now(),
			Status:      "declined",
		},
	}, nil
}

