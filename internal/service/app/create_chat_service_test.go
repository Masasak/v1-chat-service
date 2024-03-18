package app_svc_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/Masasak/v1-chat-service/internal/model"
	"github.com/Masasak/v1-chat-service/internal/port/err/status"
	"github.com/Masasak/v1-chat-service/internal/port/helper/auth"
	"github.com/Masasak/v1-chat-service/internal/port/helper/tx"
	"github.com/Masasak/v1-chat-service/internal/port/service/app"
	app_svc "github.com/Masasak/v1-chat-service/internal/service/app"
	mocks "github.com/Masasak/v1-chat-service/test/mock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateChatService(t *testing.T) {
	txManager := tx.NewManager()
	mockUserService := mocks.NewUserService(t)
	mockChatService := mocks.NewChatService(t)

	service := app_svc.NewCreateChatService(
		txManager, mockUserService, mockChatService,
	)

	anything := mock.Anything
	targetID := uuid.New()

	tc := []struct {
		desc   string
		input  app.CreateChatInput
		want   app.CreateChatOutput
		status int

		set func(t *testing.T)
	}{
		{
			desc:  "success: found",
			input: app.CreateChatInput{UserID: targetID},
			want:  app.CreateChatOutput{ChatID: targetID},
			set: func(t *testing.T) {
				mockUserService.On("Exists", anything, anything).Return(true).Once()
				mockChatService.On("FetchByUserIDs", anything, anything, anything).Return(&model.Chat{ID: targetID}).Once()
			},
		},
		{
			desc:  "success: create",
			input: app.CreateChatInput{UserID: targetID},
			want:  app.CreateChatOutput{ChatID: targetID},
			set: func(t *testing.T) {
				mockUserService.On("Exists", anything, anything).Return(true).Once()
				mockChatService.On("FetchByUserIDs", anything, anything, anything).Return(nil).Once()
				mockChatService.On("Create", anything, anything, anything, anything).Return(targetID).Once()
			},
		},
		{
			desc:   "same user id",
			input:  app.CreateChatInput{UserID: uuid.Nil},
			status: http.StatusConflict,
			set:    func(t *testing.T) {},
		},
		{
			desc:   "user not found",
			input:  app.CreateChatInput{UserID: targetID},
			status: http.StatusNotFound,
			set: func(t *testing.T) {
				mockUserService.On("Exists", anything, anything).Return(false).Once()
			},
		},
	}

	ctx := auth.Inject(context.Background(), auth.Info{UserID: uuid.Nil})
	for _, tt := range tc {
		t.Run(tt.desc, func(t *testing.T) {
			tt.set(t)

			out, err := service.Execute(ctx, tt.input)

			if tt.status != 0 {
				var statErr *status.Error
				if assert.ErrorAs(t, err, &statErr) {
					assert.Equal(t, tt.status, statErr.Code)
				}
			}

			assert.Equal(t, out, tt.want)

			mockUserService.AssertExpectations(t)
			mockChatService.AssertExpectations(t)
		})
	}
}
