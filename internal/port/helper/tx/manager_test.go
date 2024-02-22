package tx_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Masasak/v1-chat-service/internal/port/helper/tx"
	mocks "github.com/Masasak/v1-chat-service/test/mock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestManagerEvaluate(t *testing.T) {
	var (
		base = context.Background()
		id   = 1

		mockTx      = mocks.NewTxTx(t)
		mockDataSrc = mocks.NewTxDataSource(t)
		txManager   = tx.NewManager()
	)

	mockDataSrc.On("ID").Return(id)
	mockDataSrc.On("NewTx", mock.Anything).Return(mockTx, nil)

	t.Run("no errors found", func(t *testing.T) {
		// given
		mockTx.On("Commit", mock.Anything).Return(nil)

		ctx := txManager.Begin(base)
		utx := tx.Extract(ctx)
		require.NoError(t, utx.Register(base, mockDataSrc))

		// when
		var err error
		txManager.Evaluate(ctx, &err)

		// then
		mockDataSrc.AssertExpectations(t)
		mockTx.AssertExpectations(t)
	})

	t.Run("error found", func(t *testing.T) {
		// given
		mockTx.On("Rollback", mock.Anything).Return(nil)

		ctx := txManager.Begin(base)
		utx := tx.Extract(ctx)
		require.NoError(t, utx.Register(base, mockDataSrc))

		// when
		err := errors.New("")
		txManager.Evaluate(ctx, &err)

		// then
		mockDataSrc.AssertExpectations(t)
		mockTx.AssertExpectations(t)
	})

	t.Run("detect panic", func(t *testing.T) {
		// given
		mockTx.On("Rollback", mock.Anything).Return(nil)

		ctx := txManager.Begin(base)
		utx := tx.Extract(ctx)
		require.NoError(t, utx.Register(base, mockDataSrc))

		// when & then
		assert.Panics(t, func() {
			var err error
			defer txManager.Evaluate(ctx, &err)

			panic("uh-oh")
		})
		
		mockDataSrc.AssertExpectations(t)
		mockTx.AssertExpectations(t)
	})
}
