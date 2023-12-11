package local

import (
	"context"
	"errors"
	"fmt"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/models"
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

func (s Storage) SaveBatch(ctx context.Context, shortLinks []*models.ShortLink) error {
	for _, shortLink := range shortLinks {
		if err := s.Add(ctx, shortLink.Token, shortLink.Original); err != nil {
			return fmt.Errorf("add one short link: %w", err)
		}
	}

	return nil
}
