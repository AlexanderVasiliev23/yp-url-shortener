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
	Get(ctx context.Context, token string) (link *models.ShortLink, err error)
	SaveBatch(ctx context.Context, shortLinks []*models.ShortLink) error
	GetTokenByURL(ctx context.Context, url string) (string, error)
	FindByUserID(ctx context.Context, userID int) ([]*models.ShortLink, error)
	DeleteTokens(ctx context.Context, userID int, tokens []string) error
}
