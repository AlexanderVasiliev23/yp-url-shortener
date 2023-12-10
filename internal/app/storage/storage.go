package storage

import (
	"context"
	"errors"
)

var (
	ErrNotFound = errors.New("not found")
)

type Storage interface {
	Add(ctx context.Context, token, url string) error
	Get(ctx context.Context, token string) (string, error)
}
