package list

import (
	"context"
	"errors"
	"fmt"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/models"
)

var (
	ErrUnauthorized = errors.New("unauthorized")
	ErrNoSavedURLs  = errors.New("no saved URLs")
)

type OutDTO struct {
	Items []OutDTOItem
}

type OutDTOItem struct {
	ShortURL    string
	OriginalURL string
}

type linksStorage interface {
	FindByUserID(ctx context.Context, userID int) ([]*models.ShortLink, error)
}

type userContextFetcher interface {
	GetUserIDFromContext(ctx context.Context) (int, error)
}

type UseCase struct {
	storage            linksStorage
	userContextFetcher userContextFetcher
	addr               string
}

func NewUseCase(storage linksStorage, userContextFetcher userContextFetcher, addr string) *UseCase {
	return &UseCase{storage: storage, userContextFetcher: userContextFetcher, addr: addr}
}

func (u *UseCase) List(ctx context.Context) (*OutDTO, error) {
	userID, err := u.userContextFetcher.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, ErrUnauthorized
	}

	shortLinks, err := u.storage.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if len(shortLinks) == 0 {
		return nil, ErrNoSavedURLs
	}

	outDTO := &OutDTO{
		Items: make([]OutDTOItem, len(shortLinks)),
	}
	for i, shortLink := range shortLinks {
		item := OutDTOItem{
			ShortURL:    fmt.Sprintf("%s/%s", u.addr, shortLink.Token),
			OriginalURL: shortLink.Original,
		}
		outDTO.Items[i] = item
	}

	return outDTO, nil
}
