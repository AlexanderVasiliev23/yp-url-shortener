package storage

import (
	"context"
	"errors"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/models"
)

var (
	ErrNotFound = errors.New("not found")
)

type Storage interface {
	Add(ctx context.Context, token, url string) error
	Get(ctx context.Context, token string) (string, error)
	SaveBatch(ctx context.Context, shortLinks []*models.ShortLink) error
}
