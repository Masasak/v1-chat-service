package app

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type GetRandomChatInfoInput struct {
	ChatID uuid.UUID
}

type RandomChatParticipantElem struct {
	Nickname string    `json:"nickname"`
	LastRead time.Time `json:"lastRead"`
}

type GetRandomChatInfoOutput struct {
	Participants []RandomChatParticipantElem `json:"participants"`
}

type GetRandomChatInfoService interface {
	Execute(ctx context.Context, input GetRandomChatInfoInput) (output GetRandomChatInfoOutput, err error)
}
