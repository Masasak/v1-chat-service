package tx

import (
	"context"
)

// DataSource is abstraction of data source that can execute transactions.
type DataSource interface {
	// ID returns unique identifier of data source.
	// It is recommended to return a private struct which is defined in provider's package.
	// Use it like context.WithValue.
	ID() any

	// NewTx starts a new transaction and returns its abstraction.
	NewTx(ctx context.Context) (Tx, error)
}
