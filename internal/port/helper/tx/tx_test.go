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

func TestExtract(t *testing.T) {
	var (
		base      = context.Background()
		txManager = tx.NewManager()
	)

	t.Run("extract from valid context", func(t *testing.T) {
		ctx := txManager.Begin(base)

		utx := tx.Extract(ctx)

		assert.NotNil(t, utx)
	})

	t.Run("extract from invalid context", func(t *testing.T) {
		utx := tx.Extract(base)

		assert.Nil(t, utx)
	})
}

func TestUnionTxRegister(t *testing.T) {
	var (
		base = context.Background()
		id   = 1

		mockTx      = mocks.NewTxTx(t)
		mockDataSrc = mocks.NewTxDataSource(t)
		txManager   = tx.NewManager()
	)

	t.Run("register new data source", func(t *testing.T) {
		// given
		mockDataSrc.On("ID").Return(id).Once()
		mockDataSrc.On("NewTx", mock.Anything).Return(mockTx, nil).Once()

		utx := tx.Extract(txManager.Begin(base))

		// when
		err := utx.Register(base, mockDataSrc)

		// then
		assert.NoError(t, err)
		assert.Equal(t, mockTx, utx.GetTx(id))

		mockDataSrc.AssertExpectations(t)
	})

	t.Run("duplicate data sources", func(t *testing.T) {
		// given
		mockDataSrc.On("ID").Return(id).Twice()
		mockDataSrc.On("NewTx", mock.Anything).Return(mockTx, nil).Once()

		utx := tx.Extract(txManager.Begin(base))
		require.NoError(t, utx.Register(base, mockDataSrc))

		// when
		err := utx.Register(base, mockDataSrc)

		// then
		assert.NoError(t, err)
		assert.Equal(t, mockTx, utx.GetTx(id))

		mockDataSrc.AssertExpectations(t)
	})

	t.Run("creating tx failed", func(t *testing.T) {
		// given
		tmpDataSrc := mocks.NewTxDataSource(t)
		tmpDataSrc.On("ID").Return(id + 1).Once()
		tmpDataSrc.On("NewTx", mock.Anything).Return(mockTx, nil).Once()

		mockDataSrc.On("ID").Return(id).Once()
		mockDataSrc.On("NewTx", mock.Anything).Return(nil, errors.New("")).Once()
		mockTx.On("Rollback", mock.Anything).Return(nil).Once()

		utx := tx.Extract(txManager.Begin(base))
		require.NoError(t, utx.Register(base, tmpDataSrc))

		// when
		err := utx.Register(base, mockDataSrc)

		// then
		assert.Error(t, err)

		mockTx.AssertExpectations(t)
		tmpDataSrc.AssertExpectations(t)
		mockDataSrc.AssertExpectations(t)
	})
}
