package app_svc

import (
	"context"
	"net/http"
	"slices"
	"time"

	"github.com/Masasak/v1-chat-service/internal/model"
	"github.com/Masasak/v1-chat-service/internal/port/err/status"
	"github.com/Masasak/v1-chat-service/internal/port/helper/auth"
	"github.com/Masasak/v1-chat-service/internal/port/helper/tx"
	"github.com/Masasak/v1-chat-service/internal/port/service/app"
	"github.com/Masasak/v1-chat-service/internal/port/service/domain"
	"github.com/google/uuid"
)

type getRandomChatInfoService struct {
	txManager         *tx.Manager
	randomChatService domain.RandomChatService
	msgService        domain.MessageService
	userService       domain.UserService
}

var _ app.GetRandomChatInfoService = (*getRandomChatInfoService)(nil)

func NewGetRandomChatInfoService(
	tm *tx.Manager, rsvc domain.RandomChatService,
	msvc domain.MessageService, usvc domain.UserService,
) *getRandomChatInfoService {
	return &getRandomChatInfoService{
		randomChatService: rsvc,
		msgService:        msvc,
		userService:       usvc,
		txManager:         tm,
	}
}

func (svc *getRandomChatInfoService) Execute(ctx context.Context, input app.GetRandomChatInfoInput) (_ app.GetRandomChatInfoOutput, err error) {
	info := auth.MustExtract(ctx)

	ctx = svc.txManager.Begin(ctx)
	defer svc.txManager.Evaluate(ctx, &err)

	chat := svc.randomChatService.FetchChatInfo(ctx, input.ChatID)
	if chat == nil {
		return app.GetRandomChatInfoOutput{}, status.NewErr(http.StatusNotFound, "chat not found")
	}

	isParticipant := slices.Contains(chat.UserIDs, info.UserID)
	if !isParticipant {
		return app.GetRandomChatInfoOutput{}, status.NewErr(http.StatusForbidden, "you are not participant")
	}

	lastReadMap := svc.msgService.FetchLastRead(ctx, chat.ID, chat.UserIDs)
	users := svc.userService.FetchUsers(ctx, chat.UserIDs)

	return svc.transform(chat, lastReadMap, users), nil
}

func (svc *getRandomChatInfoService) transform(chat *model.RandomChat, lastReadMap map[uuid.UUID]time.Time, users []model.User) (out app.GetRandomChatInfoOutput) {
	out.Participants = make([]app.RandomChatParticipantElem, len(users))
	for idx, user := range users {
		lastRead := lastReadMap[user.ID]
		alias := chat.GetAlias(user.ID)

		out.Participants[idx] = app.RandomChatParticipantElem{
			Nickname: alias,
			LastRead: lastRead,
		}
	}
	return
}
