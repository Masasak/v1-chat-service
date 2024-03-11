package app

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type GetRandomChatListInput struct {
	Base  time.Time
	Count int
}

type RandomChatListElem struct {
	ID          uuid.UUID `json:"id"`
	LastMessage string    `json:"lastMessage"`
	UnreadCount int       `json:"unreadCount"`
}

type GetRandomChatListOutput struct {
	Chats []RandomChatListElem `json:"chats"`
}

type GetRandomChatListService interface {
	Execute(ctx context.Context, input GetRandomChatListInput) (output GetRandomChatListOutput, err error)
}
