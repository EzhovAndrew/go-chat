package dto

import (
	"github.com/go-chat/chat/internal/domain"
	chatv1 "github.com/go-chat/chat/pkg/api/chat/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ToProtoChat converts domain.Chat to proto Chat
func ToProtoChat(chat *domain.Chat) *chatv1.Chat {
	// Convert []UserID to []string
	participantIDs := make([]string, len(chat.ParticipantIDs))
	for i, id := range chat.ParticipantIDs {
		participantIDs[i] = id.String()
	}

	return &chatv1.Chat{
		ChatId:         chat.ChatID.String(),
		ParticipantIds: participantIDs,
		CreatedAt:      timestamppb.New(chat.CreatedAt),
	}
}

// ToProtoChats converts a slice of domain.Chat to proto Chat slice
func ToProtoChats(chats []*domain.Chat) []*chatv1.Chat {
	result := make([]*chatv1.Chat, len(chats))
	for i, chat := range chats {
		result[i] = ToProtoChat(chat)
	}
	return result
}

