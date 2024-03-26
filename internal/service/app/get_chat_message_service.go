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

type getChatMessageService struct {
	txManager   *tx.Manager
	chatService domain.ChatService
	msgService  domain.MessageService
}

var _ app.GetChatMessageService = (*getChatMessageService)(nil)

func NewGetChatMessageService(
	tm *tx.Manager, csvc domain.ChatService, msvc domain.MessageService,
) *getChatMessageService {
	return &getChatMessageService{
		chatService: csvc,
		msgService:  msvc,
		txManager:   tm,
	}
}

func (svc *getChatMessageService) Execute(ctx context.Context, input app.GetChatMessageInput) (_ app.GetChatMessageOutput, err error) {
	info := auth.MustExtract(ctx)

	ctx = svc.txManager.Begin(ctx)
	defer svc.txManager.Evaluate(ctx, &err)

	isParticipant := svc.chatService.IsParticipant(ctx, input.ChatID, info.UserID)
	if !isParticipant {
		return app.GetChatMessageOutput{}, status.NewErr(http.StatusForbidden, "you are not participant")
	}

	// TODO: change 15 into constant value.
	msgs := svc.msgService.FetchMessages(ctx, input.ChatID, input.Oldest, 15)

	return svc.transform(msgs), nil
}

func (svc *getChatMessageService) transform(msgs []*model.Message) (out app.GetChatMessageOutput) {
	out.Messages = make([]app.MessageListElem, len(msgs))
	for idx, msg := range msgs {
		out.Messages[idx] = app.MessageListElem{
			ID:        msg.ID,
			CreatedAt: msg.CreatedAt,
			AuthorID:  msg.AuthorID,
			Message:   msg.Content,
		}
	}

	return
}
