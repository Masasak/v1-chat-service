package app

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type GetRandomChatMessageInput struct {
	ChatID uuid.UUID
	Oldest time.Time
}

type RandomChatMessageListElem struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	FromName  string    `json:"fromName"`
	Message   string    `json:"message"`
}

type GetRandomChatMessageOutput struct {
	Messages []RandomChatMessageListElem `json:"message"`
}

type GetRandomChatMessageService interface {
	Execute(ctx context.Context, input GetRandomChatMessageInput) (GetRandomChatMessageOutput, error)
}
