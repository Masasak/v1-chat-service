package auth

import (
	"context"

	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type authInfoKey struct{}

// Info is user's information to determine which user is calling the method.
// It should be injected into context.Context to pass it through services.
type Info struct {
	UserID uuid.UUID
}

// Inject injects auth info to context.
func Inject(ctx context.Context, info Info) context.Context {
	return context.WithValue(ctx, authInfoKey{}, &info)
}

// Extract extracts auth info from context.
func Extract(ctx context.Context) *Info {
	return ctx.Value(authInfoKey{}).(*Info)
}

// MustExtract does the same with Extract. Except it panics when info == nil.
func MustExtract(ctx context.Context) Info {
	info := Extract(ctx)
	if info == nil {
		panic(errors.New("tried to extract info, but none provided"))
	}

	return *info
}
