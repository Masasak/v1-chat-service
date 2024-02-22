package tx

import (
	"context"
	"errors"
	"fmt"
)

// TODO: Make it into interface if needed.
type Manager struct{}

func NewManager() *Manager {
	return &Manager{}
}

// Begin creates a new UnionTx and injects it into context.
func (m *Manager) Begin(ctx context.Context) context.Context {
	utx := UnionTx{txs: make(map[any]Tx)}
	return context.WithValue(ctx, txKey{}, &utx)
}

// Evaluate decides whether to commit or rollback the transaction depending on the state of the error.
// When it meets panic, it rollbacks the transaction and produces panic again.
// If execution of commit or rollback fails, it appends its error to given error.
// It is recommended to call it as deferred statement right after the execution of Begin.
func (m *Manager) Evaluate(ctx context.Context, err *error) {
	if e := recover(); e != nil {
		err := fmt.Errorf("panic during transaction: %v", err)

		m.doEvaluate(ctx, &err)
		panic(err)
	}
	m.doEvaluate(ctx, err)
}

func (m *Manager) doEvaluate(ctx context.Context, err *error) {
	utx := Extract(ctx)
	if *err == nil {
		*err = utx.commit(ctx)
		return
	}
	*err = errors.Join(*err, utx.rollback(ctx))
}
