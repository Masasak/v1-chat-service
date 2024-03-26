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

type getChatInfoService struct {
	txManager   *tx.Manager
	chatService domain.ChatService
	msgService  domain.MessageService
	userService domain.UserService
}

var _ app.GetChatInfoService = (*getChatInfoService)(nil)

func NewGetChatInfoService(
	tm *tx.Manager, csvc domain.ChatService,
	msvc domain.MessageService, usvc domain.UserService,
) *getChatInfoService {
	return &getChatInfoService{
		chatService: csvc,
		msgService:  msvc,
		userService: usvc,
		txManager:   tm,
	}
}

func (svc *getChatInfoService) Execute(ctx context.Context, input app.GetChatInfoInput) (_ app.GetChatInfoOutput, err error) {
	info := auth.MustExtract(ctx)

	ctx = svc.txManager.Begin(ctx)
	defer svc.txManager.Evaluate(ctx, &err)

	chat := svc.chatService.FetchChatInfo(ctx, input.ChatID)
	if chat == nil {
		return app.GetChatInfoOutput{}, status.NewErr(http.StatusNotFound, "chat not found")
	}

	isParticipant := slices.Contains(chat.UserIDs, info.UserID)
	if !isParticipant {
		return app.GetChatInfoOutput{}, status.NewErr(http.StatusForbidden, "you are not participant")
	}

	lastReadMap := svc.msgService.FetchLastRead(ctx, chat.ID, chat.UserIDs)
	users := svc.userService.FetchUsers(ctx, chat.UserIDs)

	return svc.transform(lastReadMap, users), nil
}

func (svc *getChatInfoService) transform(lastReadMap map[uuid.UUID]time.Time, users []model.User) (out app.GetChatInfoOutput) {
	out.Participants = make([]app.ChatParticipantElem, len(users))
	for idx, user := range users {
		lastRead := lastReadMap[user.ID]

		out.Participants[idx] = app.ChatParticipantElem{
			ID:              user.ID,
			ProfileImageURL: user.ProfileImageURL,
			Nickname:        user.Nickname,
			LastRead:        lastRead,
		}
	}
	return
}
