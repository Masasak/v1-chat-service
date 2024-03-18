package app_svc

import (
	"context"

	"github.com/Masasak/v1-chat-service/internal/model"
	"github.com/Masasak/v1-chat-service/internal/port/helper/auth"
	"github.com/Masasak/v1-chat-service/internal/port/helper/tx"
	"github.com/Masasak/v1-chat-service/internal/port/service/app"
	"github.com/Masasak/v1-chat-service/internal/port/service/domain"
	"github.com/google/uuid"
)

type getRandomChatListService struct {
	randomChatService domain.RandomChatService
	msgService        domain.MessageService
	txManager         *tx.Manager
}

var _ app.GetRandomChatListService = (*getRandomChatListService)(nil)

func NewGetRandomChatListService(tm *tx.Manager, rsvc domain.RandomChatService, msvc domain.MessageService) *getRandomChatListService {
	return &getRandomChatListService{
		randomChatService: rsvc,
		msgService:        msvc,
		txManager:         tm,
	}
}

func (svc *getRandomChatListService) Execute(ctx context.Context, input app.GetRandomChatListInput) (_ app.GetRandomChatListOutput, err error) {
	info := auth.MustExtract(ctx)

	ctx = svc.txManager.Begin(ctx)
	defer svc.txManager.Evaluate(ctx, &err)

	chats := svc.randomChatService.FetchChatsBefore(ctx, info.UserID, input.Base, input.Count)

	chatIDs := make([]uuid.UUID, len(chats))
	for idx, chat := range chats {
		chatIDs[idx] = chat.ID
	}

	chatInfo := svc.msgService.FetchUserInfo(ctx, info.UserID, chatIDs)
	return svc.transform(chats, chatInfo), nil
}

func (svc *getRandomChatListService) transform(chats []*model.RandomChat, chatInfo map[uuid.UUID]domain.UserChatInfo) (out app.GetRandomChatListOutput) {
	out.Chats = make([]app.RandomChatListElem, len(chats))
	for idx, chat := range chats {
		info := chatInfo[chat.ID]

		out.Chats[idx] = app.RandomChatListElem{
			ID:          chat.ID,
			LastMessage: info.LastMessage,
			UnreadCount: info.UnreadCount,
		}
	}

	return
}
