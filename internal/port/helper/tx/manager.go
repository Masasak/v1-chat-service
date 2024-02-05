package tx

import (
	"context"
	"errors"
)

type Manager struct{}

// Begin creates a new UnionTx and injects it into context.
func (m *Manager) Begin(ctx context.Context) context.Context {
	utx := UnionTx{txs: make(map[any]Tx)}
	return context.WithValue(ctx, txKey{}, &utx)
}

// Evaluate decides whether to commit or roll back the transaction depending on the state of the error.
// If execution of commit or rollback fails, it appends its error to given error.
// It is recommended to call it as deferred statement right after the excution of Begin.
func (m *Manager) Evaluate(ctx context.Context, err *error) {
	utx := Extract(ctx)
	if *err == nil {
		*err = utx.commit(ctx)
		return
	}
	*err = errors.Join(*err, utx.rollback(ctx))
}
