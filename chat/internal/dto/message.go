package dto

import (
	"github.com/go-chat/chat/internal/domain"
	chatv1 "github.com/go-chat/chat/pkg/api/chat/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ToProtoMessage converts domain.Message to proto Message
func ToProtoMessage(msg *domain.Message) *chatv1.Message {
	return &chatv1.Message{
		MessageId: msg.MessageID.String(),
		ChatId:    msg.ChatID.String(),
		SenderId:  msg.SenderID.String(),
		Text:      msg.Text,
		CreatedAt: timestamppb.New(msg.CreatedAt),
	}
}

// ToProtoMessages converts a slice of domain.Message to proto Message slice
func ToProtoMessages(messages []*domain.Message) []*chatv1.Message {
	result := make([]*chatv1.Message, len(messages))
	for i, msg := range messages {
		result[i] = ToProtoMessage(msg)
	}
	return result
}

