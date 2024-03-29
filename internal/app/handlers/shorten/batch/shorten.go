package batch

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/models"
)

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

// Shortener missing godoc.
type Shortener struct {
	saver              batchSaver
	tokenGenerator     tokenGenerator
	uuidGenerator      uuidGenerator
	userContextFetcher userContextFetcher
	addr               string
}

// NewShortener missing godoc.
func NewShortener(
	saver batchSaver,
	tokenGenerator tokenGenerator,
	uuidGenerator uuidGenerator,
	userContextFetcher userContextFetcher,
	addr string,
) *Shortener {
	return &Shortener{
		saver:              saver,
		tokenGenerator:     tokenGenerator,
		uuidGenerator:      uuidGenerator,
		userContextFetcher: userContextFetcher,
		addr:               addr,
	}
}

// Handle missing godoc.
func (h *Shortener) Handle(c echo.Context) error {
	type reqItem struct {
		CorrelationID string `json:"correlation_id"`
		OriginalURL   string `json:"original_url"`
	}

	type req []reqItem

	var requestItems req

	if err := json.NewDecoder(c.Request().Body).Decode(&requestItems); err != nil {
		c.Response().WriteHeader(http.StatusBadRequest)
		return err
	}

	type respItem struct {
		CorrelationID string `json:"correlation_id"`
		ShortURL      string `json:"short_url"`
	}

	type responseItems []respItem

	var response responseItems

	toSave := make([]*models.ShortLink, 0, len(requestItems))

	userID, err := h.userContextFetcher.GetUserIDFromContext(c.Request().Context())
	if err != nil {
		c.Response().WriteHeader(http.StatusInternalServerError)
		return err
	}

	for _, requestItem := range requestItems {
		token, err := h.tokenGenerator.Generate()
		if err != nil {
			c.Response().WriteHeader(http.StatusInternalServerError)
			return err
		}

		shortLink := models.NewShortLink(userID, h.uuidGenerator.Generate(), token, requestItem.OriginalURL)
		toSave = append(toSave, shortLink)

		_respItem := respItem{
			CorrelationID: requestItem.CorrelationID,
			ShortURL:      h.addr + "/" + token,
		}

		response = append(response, _respItem)
	}

	if err := h.saver.SaveBatch(c.Request().Context(), toSave); err != nil {
		c.Response().WriteHeader(http.StatusInternalServerError)
		return err
	}

	c.Response().Header().Set("Content-Type", "application/json")
	c.Response().WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(c.Response().Writer).Encode(response); err != nil {
		c.Response().WriteHeader(http.StatusInternalServerError)
		return err
	}

	return nil
}
