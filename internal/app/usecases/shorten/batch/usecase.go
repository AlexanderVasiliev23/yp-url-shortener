package batch

import (
	"context"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/models"
	"github.com/google/uuid"
)

type InDTO struct {
	Items []InDTOItem
}

type InDTOItem struct {
	CorrelationID string
	OriginalURL   string
}

type OutDTO struct {
	Items []OutDTOItem
}

type OutDTOItem struct {
	CorrelationID string
	ShortURL      string
}

type batchSaver interface {
	SaveBatch(ctx context.Context, shortLinks []*models.ShortLink) error
}

type tokenGenerator interface {
	Generate() (string, error)
}

type uuidGenerator interface {
	Generate() uuid.UUID
}

type userContextFetcher interface {
	GetUserIDFromContext(ctx context.Context) (int, error)
}

type UseCase struct {
	saver              batchSaver
	tokenGenerator     tokenGenerator
	uuidGenerator      uuidGenerator
	userContextFetcher userContextFetcher
	addr               string
}

func NewUseCase(saver batchSaver, tokenGenerator tokenGenerator, uuidGenerator uuidGenerator, userContextFetcher userContextFetcher, addr string) *UseCase {
	return &UseCase{saver: saver, tokenGenerator: tokenGenerator, uuidGenerator: uuidGenerator, userContextFetcher: userContextFetcher, addr: addr}
}

func (u *UseCase) Shorten(ctx context.Context, in InDTO) (*OutDTO, error) {
	out := &OutDTO{
		Items: make([]OutDTOItem, 0, len(in.Items)),
	}

	toSave := make([]*models.ShortLink, 0, len(in.Items))

	userID, err := u.userContextFetcher.GetUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	for _, requestItem := range in.Items {
		token, err := u.tokenGenerator.Generate()
		if err != nil {
			return nil, err
		}

		shortLink := models.NewShortLink(userID, u.uuidGenerator.Generate(), token, requestItem.OriginalURL)
		toSave = append(toSave, shortLink)

		out.Items = append(out.Items, OutDTOItem{
			CorrelationID: requestItem.CorrelationID,
			ShortURL:      u.addr + "/" + token,
		})
	}

	if err := u.saver.SaveBatch(ctx, toSave); err != nil {
		return nil, err
	}

	return out, nil
}
