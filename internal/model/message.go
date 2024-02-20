package model

import (
	"time"

	"github.com/google/uuid"
)

// Message is a message created in a chat.
type Message struct {
	ID        uuid.UUID
	ChatID    uuid.UUID
	AuthorID  uuid.UUID
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
}
