package dto

import (
	"github.com/go-chat/notifications/internal/domain"
	notificationsv1 "github.com/go-chat/notifications/pkg/api/notifications/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ToProtoNotification converts domain.Notification to proto Notification
func ToProtoNotification(n *domain.Notification) *notificationsv1.Notification {
	return &notificationsv1.Notification{
		NotificationId: n.NotificationID.String(),
		UserId:         n.UserID.String(),
		Type:           toProtoNotificationType(n.Type),
		Title:          n.Title,
		Body:           n.Body,
		Data:           n.Data,
		CreatedAt:      timestamppb.New(n.CreatedAt),
		Read:           n.Read,
	}
}

// ToProtoNotifications converts a slice of domain.Notification to proto Notification slice
func ToProtoNotifications(notifications []*domain.Notification) []*notificationsv1.Notification {
	result := make([]*notificationsv1.Notification, len(notifications))
	for i, n := range notifications {
		result[i] = ToProtoNotification(n)
	}
	return result
}

// toProtoNotificationType converts domain.NotificationType to proto NotificationType
func toProtoNotificationType(t domain.NotificationType) notificationsv1.NotificationType {
	switch t {
	case domain.NotificationTypeMessage:
		return notificationsv1.NotificationType_NOTIFICATION_TYPE_MESSAGE
	case domain.NotificationTypeFriendRequest:
		return notificationsv1.NotificationType_NOTIFICATION_TYPE_FRIEND_REQUEST
	case domain.NotificationTypeFriendAccepted:
		return notificationsv1.NotificationType_NOTIFICATION_TYPE_FRIEND_ACCEPTED
	default:
		return notificationsv1.NotificationType_NOTIFICATION_TYPE_UNSPECIFIED
	}
}

