package local

import (
	"context"
	"errors"
)

var (
	ErrURLNotFound = errors.New("url is not found")
)

type Storage map[string]string

func New() *Storage {
	s := make(Storage)
	return &s
}

func (s Storage) Add(ctx context.Context, token, url string) error {
	s[token] = url

	return nil
}

func (s Storage) Get(ctx context.Context, token string) (string, error) {
	url, ok := s[token]
	if ok {
		return url, nil
	}

	return "", ErrURLNotFound
}
