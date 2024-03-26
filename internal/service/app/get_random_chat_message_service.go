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
)

type getRandomChatMessageService struct {
	txManager         *tx.Manager
	randomChatService domain.RandomChatService
	msgService        domain.MessageService
}

var _ app.GetRandomChatMessageService = (*getRandomChatMessageService)(nil)

func NewGetRandomChatMessageService(
	tm *tx.Manager, rsvc domain.RandomChatService, msvc domain.MessageService,
) *getRandomChatMessageService {
	return &getRandomChatMessageService{
		randomChatService: rsvc,
		msgService:        msvc,
		txManager:         tm,
	}
}

func (svc *getRandomChatMessageService) Execute(ctx context.Context, input app.GetRandomChatMessageInput) (_ app.GetRandomChatMessageOutput, err error) {
	info := auth.MustExtract(ctx)

	ctx = svc.txManager.Begin(ctx)
	defer svc.txManager.Evaluate(ctx, &err)

	chat := svc.randomChatService.FetchChatInfo(ctx, input.ChatID)
	if chat == nil {
		return app.GetRandomChatMessageOutput{}, status.NewErr(http.StatusNotFound, "chat not found")
	}

	found := false
	for _, userID := range chat.UserIDs {
		if info.UserID == userID {
			found = true
			break
		}
	}

	if !found {
		return app.GetRandomChatMessageOutput{}, status.NewErr(http.StatusForbidden, "you are not participant")
	}

	// TODO: change 15 into constant value.
	msgs := svc.msgService.FetchMessages(ctx, input.ChatID, input.Oldest, 15)

	return svc.transform(chat, msgs), nil
}

func (svc *getRandomChatMessageService) transform(chat *model.RandomChat, msgs []*model.Message) (out app.GetRandomChatMessageOutput) {
	out.Messages = make([]app.RandomChatMessageListElem, len(msgs))
	for idx, msg := range msgs {
		alias := chat.GetAlias(msg.AuthorID)

		out.Messages[idx] = app.RandomChatMessageListElem{
			ID:        msg.ID,
			CreatedAt: msg.CreatedAt,
			FromName:  alias,
			Message:   msg.Content,
		}
	}

	return
}
