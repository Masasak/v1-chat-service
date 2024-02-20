package model

import "github.com/google/uuid"

type ChatKind uint8

const (
	ChatKindGeneral ChatKind = iota
	ChatKindRandom
)

// Chat is a chat room where users can send message to each other.
type Chat struct {
	ID      uuid.UUID
	Kind    ChatKind
	UserIDs []uuid.UUID
}
