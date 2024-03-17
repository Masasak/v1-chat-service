package domain

import (
	"context"
	"time"

	"github.com/Masasak/v1-chat-service/internal/model"
	"github.com/google/uuid"
)

type ChatService interface {
	FetchChatsBefore(ctx context.Context, userID uuid.UUID, base time.Time, take int) []*model.Chat
}
