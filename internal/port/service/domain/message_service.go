package domain

import (
	"context"
	"time"

	"github.com/Masasak/v1-chat-service/internal/model"
	"github.com/google/uuid"
)

type UserChatInfo struct {
	UnreadCount int
	LastMessage string
}

type MessageService interface {
	FetchUserInfo(ctx context.Context, userID uuid.UUID, chatIDs []uuid.UUID) map[uuid.UUID]UserChatInfo
	FetchMessages(ctx context.Context, chatID uuid.UUID, before time.Time, take int) []*model.Message
}
