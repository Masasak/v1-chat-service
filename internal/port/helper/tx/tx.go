package tx

import (
	"context"
	"errors"
)

// Tx is abstraction of transaction.
type Tx interface {
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

// txKey is unique identifier of UnionTx in ctx.
type txKey struct{}

// Extract extracts UnionTx from context.
func Extract(ctx context.Context) *UnionTx {
	if utx, ok := ctx.Value(txKey{}).(*UnionTx); ok {
		return utx
	}
	// TODO: Panic or do something.
	return nil
}

// UnionTx is union of transactions.
// it holds platform-independent transactions in it.
type UnionTx struct {
	txs map[any]Tx
}

// Register registers data source to UnionTx. It starts new transaction from data source internally.
// If duplicate data source is found, then no-op.
func (utx *UnionTx) Register(ctx context.Context, dataSources ...DataSource) error {
	for _, ds := range dataSources {
		id := ds.ID()
		if _, ok := utx.txs[id]; !ok {
			tx, err := ds.NewTx(ctx)
			if err != nil {
				e := utx.rollback(ctx)
				return errors.Join(err, e)
			}
			utx.txs[id] = tx
		}
	}
	return nil
}

// GetTx finds platform-specific transaction from UnionTx.
func (utx *UnionTx) GetTx(id any) Tx { return utx.txs[id] }

func (utx *UnionTx) commit(ctx context.Context) error {
	errs := make([]error, 0)
	for _, tx := range utx.txs {
		if e := tx.Commit(ctx); e != nil {
			// TODO: We might need to implement retry.
			errs = append(errs, e)
		}
	}
	return errors.Join(errs...)
}

func (utx *UnionTx) rollback(ctx context.Context) error {
	errs := make([]error, 0)
	for _, tx := range utx.txs {
		if e := tx.Rollback(ctx); e != nil {
			// TODO: We might need to implement retry.
			errs = append(errs, e)
		}
	}
	return errors.Join(errs...)
}
