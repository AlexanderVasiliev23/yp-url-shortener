package single

import (
	"context"
	"errors"
	"fmt"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/models"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/storage"
	"github.com/google/uuid"
)

var (
	ErrInvalidJSON      = errors.New("invalid json")
	ErrEmptyOriginalURL = errors.New("empty original url")
	ErrAlreadyExists    = errors.New("already exists")
)

type repository interface {
	Add(ctx context.Context, shortLink *models.ShortLink) error
	GetTokenByURL(ctx context.Context, url string) (string, error)
}

type tokenGenerator interface {
	Generate() (string, error)
}

type userContextFetcher interface {
	GetUserIDFromContext(ctx context.Context) (int, error)
}

type UseCase struct {
	repository         repository
	tokenGenerator     tokenGenerator
	userContextFetcher userContextFetcher
	addr               string
}

func NewUseCase(repository repository, tokenGenerator tokenGenerator, userContextFetcher userContextFetcher, addr string) *UseCase {
	return &UseCase{repository: repository, tokenGenerator: tokenGenerator, userContextFetcher: userContextFetcher, addr: addr}
}

func (u *UseCase) Shorten(ctx context.Context, originalURL string) (shortURL string, err error) {
	if originalURL == "" {
		return "", ErrEmptyOriginalURL
	}

	token, err := u.tokenGenerator.Generate()
	if err != nil {
		return "", err
	}

	userID, err := u.userContextFetcher.GetUserIDFromContext(ctx)
	if err != nil {
		return "", err
	}

	model := models.NewShortLink(userID, uuid.New(), token, originalURL)
	if err := u.repository.Add(ctx, model); err != nil {
		if !errors.Is(err, storage.ErrAlreadyExists) {
			return "", err
		}

		_token, err := u.repository.GetTokenByURL(ctx, originalURL)
		if err != nil {
			return "", err
		}

		return fmt.Sprintf("%s/%s", u.addr, _token), ErrAlreadyExists
	}

	return fmt.Sprintf("%s/%s", u.addr, token), nil
}
