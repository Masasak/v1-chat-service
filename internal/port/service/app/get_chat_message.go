package app

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type GetChatMessageInput struct {
	ChatID uuid.UUID
	Oldest time.Time
}

type MessageListElem struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	AuthorID  uuid.UUID `json:"authorId"`
	Message   string    `json:"message"`
}

type GetChatMessageOutput struct {
	Messages []MessageListElem `json:"mssages"`
}

type GetChatMessageService interface {
	Execute(ctx context.Context, input GetChatMessageInput) (GetChatMessageOutput, error)
}
