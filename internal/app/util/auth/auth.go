package auth

import (
	"context"
	"errors"
)

type contextKey string

const (
	contextUserIDFieldName contextKey = "user_id"
)

// ErrNotFound missing godoc.
var (
	ErrNotFound = errors.New("user not found")
)

// UserContextFetcher missing godoc.
type UserContextFetcher struct {
}

// GetUserIDFromContext missing godoc.
func (f *UserContextFetcher) GetUserIDFromContext(ctx context.Context) (int, error) {
	val, ok := ctx.Value(contextUserIDFieldName).(int)
	if ok {
		return val, nil
	}

	return 0, ErrNotFound
}

// WithUserID missing godoc.
func WithUserID(ctx context.Context, userID int) context.Context {
	return context.WithValue(ctx, contextUserIDFieldName, userID)
}
