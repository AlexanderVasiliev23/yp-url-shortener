package local

import (
	"context"
	"errors"
	"fmt"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/models"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/storage"
)

var _ storage.Storage = (*Storage)(nil)

var (
	ErrURLNotFound = errors.New("url is not found")
)

type Storage struct {
	tokenToURLMap map[string]string
	URLToTokenMap map[string]string
}

func New() *Storage {
	return &Storage{
		tokenToURLMap: make(map[string]string),
		URLToTokenMap: make(map[string]string),
	}
}

func (s Storage) Add(ctx context.Context, token, url string) error {
	if _, ok := s.URLToTokenMap[url]; ok {
		return storage.ErrAlreadyExists
	}

	s.tokenToURLMap[token] = url
	s.URLToTokenMap[url] = token

	return nil
}

func (s Storage) Get(ctx context.Context, token string) (string, error) {
	url, ok := s.tokenToURLMap[token]
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

func (s Storage) GetTokenByURL(ctx context.Context, url string) (string, error) {
	token, ok := s.URLToTokenMap[url]
	if !ok {
		return "", storage.ErrNotFound
	}

	return token, nil
}
