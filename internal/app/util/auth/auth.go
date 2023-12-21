package auth

import (
	"context"
	"errors"
)

const (
	contextUserIdFieldName = "user_id"
)

var (
	ErrNotFound = errors.New("user not found")
)

type UserContextFetcher struct {
}

func (f *UserContextFetcher) GetUserIdFromContext(ctx context.Context) (int, error) {
	val, ok := ctx.Value(contextUserIdFieldName).(int)
	if ok {
		return val, nil
	}

	return 0, ErrNotFound
}

func WithUserId(ctx context.Context, userId int) context.Context {
	return context.WithValue(ctx, contextUserIdFieldName, userId)
}
