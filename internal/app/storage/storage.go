package storage

import (
	"context"
	"errors"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/models"
)

var (
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("already exists")
)

type Storage interface {
	Add(ctx context.Context, shortLink *models.ShortLink) error
	Get(ctx context.Context, token string) (string, error)
	SaveBatch(ctx context.Context, shortLinks []*models.ShortLink) error
	GetTokenByURL(ctx context.Context, url string) (string, error)
	FindByUserId(ctx context.Context, userId int) ([]*models.ShortLink, error)
}
