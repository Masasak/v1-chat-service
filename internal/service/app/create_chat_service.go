package app_svc

import (
	"context"
	"net/http"

	"github.com/Masasak/v1-chat-service/internal/model"
	"github.com/Masasak/v1-chat-service/internal/port/err/status"
	"github.com/Masasak/v1-chat-service/internal/port/helper/auth"
	"github.com/Masasak/v1-chat-service/internal/port/helper/tx"
	"github.com/Masasak/v1-chat-service/internal/port/service/app"
	"github.com/Masasak/v1-chat-service/internal/port/service/domain"
	"github.com/google/uuid"
)

type createChatService struct {
	userService domain.UserService
	chatService domain.ChatService
	txManager   *tx.Manager
}

var _ app.CreateChatService = (*createChatService)(nil)

func NewCreateChatService(tm *tx.Manager, usvc domain.UserService, csvc domain.ChatService) *createChatService {
	return &createChatService{
		userService: usvc,
		chatService: csvc,
		txManager:   tm,
	}
}

func (svc *createChatService) Execute(ctx context.Context, input app.CreateChatInput) (_ app.CreateChatOutput, err error) {
	info := auth.MustExtract(ctx)

	ctx = svc.txManager.Begin(ctx)
	defer svc.txManager.Evaluate(ctx, &err)

	if info.UserID == input.UserID {
		return app.CreateChatOutput{}, status.NewErr(http.StatusConflict, "participants' ids are the same")
	}

	exists := svc.userService.Exists(ctx, input.UserID)
	if !exists {
		return app.CreateChatOutput{}, status.NewErr(http.StatusNotFound, "user not found")
	}

	var chatID uuid.UUID

	chat := svc.chatService.FetchByUserIDs(ctx, info.UserID, input.UserID)
	if chat != nil && chat.Kind != model.ChatKindRandom {
		chatID = chat.ID
	} else {
		chatID = svc.chatService.Create(ctx, model.ChatKindGeneral, info.UserID, input.UserID)
	}

	return svc.transform(chatID), nil
}

func (svc *createChatService) transform(chatID uuid.UUID) (out app.CreateChatOutput) {
	return app.CreateChatOutput{
		ChatID: chatID,
	}
}
