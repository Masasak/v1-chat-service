package app

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type GetChatListInput struct {
	Base  time.Time
	Count int
}

type UserInfoElem struct {
	ID              uuid.UUID `json:"id"`
	Nickname        string    `json:"nickname"`
	ProfileImageURL *string   `json:"profileImageUrl"`
}

type ChatListElem struct {
	ID          uuid.UUID `json:"id"`
	LastMessage string    `json:"lastMessage"`
	UnreadCount int       `json:"unreadCount"`

	User UserInfoElem `json:"user"`
}

type GetChatListOutput struct {
	Chats []ChatListElem `json:"chats"`
}

type GetChatListService interface {
	Execute(ctx context.Context, input GetChatListInput) (output GetChatListOutput, err error)
}
