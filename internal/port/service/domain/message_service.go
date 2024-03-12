package domain

import (
	"context"

	"github.com/google/uuid"
)

type UserChatInfo struct {
	UnreadCount int
	LastMessage string
}

type MessageService interface {
	FetchUserInfo(ctx context.Context, userID uuid.UUID, chatIDs []uuid.UUID) map[uuid.UUID]UserChatInfo
}
