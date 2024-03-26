package app

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type GetChatInfoInput struct {
	ChatID uuid.UUID
}

type ChatParticipantElem struct {
	ID              uuid.UUID `json:"id"`
	ProfileImageURL string    `json:"profileImageUrl"`
	Nickname        string    `json:"nickname"`
	LastRead        time.Time `json:"lastRead"`
}

type GetChatInfoOutput struct {
	Participants []ChatParticipantElem `json:"participants"`
}

type GetChatInfoService interface {
	Execute(ctx context.Context, input GetChatInfoInput) (output GetChatInfoOutput, err error)
}
