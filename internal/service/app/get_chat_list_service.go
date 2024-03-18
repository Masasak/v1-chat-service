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

type getChatListService struct {
	userService domain.UserService
	chatService domain.ChatService
	msgService  domain.MessageService
	txManager   *tx.Manager
}

var _ app.GetChatListService = (*getChatListService)(nil)

func NewGetChatListService(
	tm *tx.Manager, csvc domain.ChatService,
	msvc domain.MessageService, usvc domain.UserService,
) *getChatListService {
	return &getChatListService{
		chatService: csvc,
		msgService:  msvc,
		userService: usvc,
		txManager:   tm,
	}
}

func (svc *getChatListService) Execute(ctx context.Context, input app.GetChatListInput) (_ app.GetChatListOutput, err error) {
	info := auth.MustExtract(ctx)

	ctx = svc.txManager.Begin(ctx)
	defer svc.txManager.Evaluate(ctx, &err)

	chats := svc.chatService.FetchChatsBefore(ctx, info.UserID, input.Base, input.Count)

	chatIDs := make([]uuid.UUID, len(chats))
	for idx, chat := range chats {
		chatIDs[idx] = chat.ID
	}

	chatInfo := svc.msgService.FetchUserInfo(ctx, info.UserID, chatIDs)
	opponents := svc.userService.FetchOpponents(ctx, info.UserID, chats)

	return svc.transform(chats, chatInfo, opponents), nil
}

func (svc *getChatListService) transform(chats []*model.Chat, chatInfo map[uuid.UUID]domain.UserChatInfo, opponents map[uuid.UUID]model.User) (out app.GetChatListOutput) {
	out.Chats = make([]app.ChatListElem, len(chats))
	for idx, chat := range chats {
		info := chatInfo[chat.ID]
		opponent := opponents[chat.ID]

		out.Chats[idx] = app.ChatListElem{
			ID:          chat.ID,
			LastMessage: info.LastMessage,
			UnreadCount: info.UnreadCount,
			User: app.UserInfoElem{
				ID:              opponent.ID,
				Nickname:        opponent.Nickname,
				ProfileImageURL: &opponent.ProfileImageURL,
			},
		}
	}

	return
}
