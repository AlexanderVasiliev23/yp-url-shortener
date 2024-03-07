package local

import (
	"context"
	"errors"
	"fmt"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/models"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/storage"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/uuidgenerator"
)

var _ storage.Storage = (*Storage)(nil)

// ErrURLNotFound missing godoc.
var (
	ErrURLNotFound = errors.New("url is not found")
)

// Storage missing godoc.
type Storage struct {
	uuidGenerator uuidgenerator.UUIDGenerator

	tokenToShortLinkMap  map[string]*models.ShortLink
	urlToShortLinkMap    map[string]*models.ShortLink
	userIDToShortLinkMap map[int][]*models.ShortLink
}

// New missing godoc.
func New(uuidGenerator uuidgenerator.UUIDGenerator) *Storage {
	return &Storage{
		uuidGenerator:        uuidGenerator,
		tokenToShortLinkMap:  make(map[string]*models.ShortLink),
		urlToShortLinkMap:    make(map[string]*models.ShortLink),
		userIDToShortLinkMap: make(map[int][]*models.ShortLink),
	}
}

// Add missing godoc.
func (s Storage) Add(ctx context.Context, shortLink *models.ShortLink) error {
	if _, ok := s.urlToShortLinkMap[shortLink.Original]; ok {
		return storage.ErrAlreadyExists
	}

	s.tokenToShortLinkMap[shortLink.Token] = shortLink
	s.urlToShortLinkMap[shortLink.Original] = shortLink
	s.userIDToShortLinkMap[shortLink.UserID] = append(s.userIDToShortLinkMap[shortLink.UserID], shortLink)

	return nil
}

// Get missing godoc.
func (s Storage) Get(ctx context.Context, token string) (*models.ShortLink, error) {
	shortLink, ok := s.tokenToShortLinkMap[token]
	if ok {
		return shortLink, nil
	}

	return nil, ErrURLNotFound
}

// SaveBatch missing godoc.
func (s Storage) SaveBatch(ctx context.Context, shortLinks []*models.ShortLink) error {
	for _, shortLink := range shortLinks {
		if err := s.Add(ctx, shortLink); err != nil {
			return fmt.Errorf("add one short link: %w", err)
		}
	}

	return nil
}

// GetTokenByURL missing godoc.
func (s Storage) GetTokenByURL(ctx context.Context, url string) (string, error) {
	shortLink, ok := s.urlToShortLinkMap[url]
	if !ok {
		return "", storage.ErrNotFound
	}

	return shortLink.Token, nil
}

// FindByUserID missing godoc.
func (s Storage) FindByUserID(ctx context.Context, userID int) ([]*models.ShortLink, error) {
	return s.userIDToShortLinkMap[userID], nil
}

// DeleteByTokens missing godoc.
func (s Storage) DeleteByTokens(ctx context.Context, tokens []string) error {
	for _, token := range tokens {
		shortLink, ok := s.tokenToShortLinkMap[token]
		if !ok {
			continue
		}

		shortLink.Delete()
	}

	return nil
}

// FilterOnlyThisUserTokens missing godoc.
func (s Storage) FilterOnlyThisUserTokens(ctx context.Context, userID int, tokens []string) ([]string, error) {
	result := make([]string, 0, len(tokens))

	for _, token := range tokens {
		shortLink, ok := s.tokenToShortLinkMap[token]
		if !ok {
			continue
		}
		if shortLink.UserID != userID {
			continue
		}
		result = append(result, token)
	}

	return result, nil
}
