package get

import (
	"context"
	"errors"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/models"
)

var (
	ErrTokenIsEmpty   = errors.New("token is empty")
	ErrTokenIsDeleted = errors.New("token is deleted")
)

type repository interface {
	Get(ctx context.Context, token string) (*models.ShortLink, error)
}

type UseCase struct {
	repository repository
}

func NewUseCase(repository repository) *UseCase {
	return &UseCase{repository: repository}
}

func (u *UseCase) Get(ctx context.Context, token string) (originalURL string, err error) {
	if token == "" {
		return "", ErrTokenIsEmpty
	}

	shortLink, err := u.repository.Get(ctx, token)
	if err != nil {
		return "", err
	}
	if shortLink.DeletedAt != nil {
		return "", ErrTokenIsDeleted
	}

	return shortLink.Original, nil
}
