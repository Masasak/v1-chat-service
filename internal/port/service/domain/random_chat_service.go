package domain

import (
	"context"
	"time"

	"github.com/Masasak/v1-chat-service/internal/model"
	"github.com/google/uuid"
)

type RandomChatService interface {
	Create(ctx context.Context, chat model.Chat) model.RandomChat
	MakeRevealVote(ctx context.Context, chatID uuid.UUID, userID uuid.UUID) (model.RevealVote, error)

	FetchChatsBefore(ctx context.Context, userID uuid.UUID, base time.Time, take int) []*model.RandomChat
	FetchChatInfo(ctx context.Context, chatID uuid.UUID) *model.RandomChat
}
