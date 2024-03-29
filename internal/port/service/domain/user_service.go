package domain

import (
	"context"

	"github.com/Masasak/v1-chat-service/internal/model"
	"github.com/google/uuid"
)

type UserService interface {
	FetchOpponents(ctx context.Context, userID uuid.UUID, chats []*model.Chat) map[uuid.UUID]model.User
	FetchUsers(ctx context.Context, userIDs []uuid.UUID) []model.User
	Exists(ctx context.Context, userID uuid.UUID) bool
}
