package domain

import (
	"context"
	"time"

	"github.com/Masasak/v1-chat-service/internal/model"
	"github.com/google/uuid"
)

type ChatService interface {
	FetchChatsBefore(ctx context.Context, userID uuid.UUID, base time.Time, take int) []*model.Chat
	FetchByUserIDs(ctx context.Context, userIDs ...uuid.UUID) *model.Chat
	FetchChatInfo(ctx context.Context, chatID uuid.UUID) *model.Chat
	IsParticipant(ctx context.Context, chatID uuid.UUID, userID uuid.UUID) bool
	Create(ctx context.Context, kind model.ChatKind, from uuid.UUID, to uuid.UUID) uuid.UUID
}
