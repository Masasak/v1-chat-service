package app

import (
	"context"

	"github.com/google/uuid"
)

type CreateChatInput struct {
	UserID uuid.UUID
}

type CreateChatOutput struct {
	ChatID uuid.UUID
}

type CreateChatService interface {
	Execute(ctx context.Context, input CreateChatInput) (output CreateChatOutput, err error)
}
