package batch

import (
	"context"
	"encoding/json"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/usecases/shorten/batch"
	"net/http"

	"github.com/labstack/echo/v4"
)

type useCase interface {
	Shorten(ctx context.Context, in batch.InDTO) (*batch.OutDTO, error)
}

// Shortener missing godoc.
type Shortener struct {
	useCase useCase
}

// NewShortener missing godoc.
func NewShortener(useCase useCase) *Shortener {
	return &Shortener{
		useCase: useCase,
	}
}

// Handle missing godoc.
func (h *Shortener) Handle(c echo.Context) error {
	var requestItems []struct {
		CorrelationID string `json:"correlation_id"`
		OriginalURL   string `json:"original_url"`
	}

	if err := json.NewDecoder(c.Request().Body).Decode(&requestItems); err != nil {
		c.Response().WriteHeader(http.StatusBadRequest)
		return err
	}

	inDTO := batch.InDTO{Items: make([]batch.InDTOItem, 0, len(requestItems))}
	for _, item := range requestItems {
		inDTO.Items = append(inDTO.Items, batch.InDTOItem{
			CorrelationID: item.CorrelationID,
			OriginalURL:   item.OriginalURL,
		})
	}

	outDTO, err := h.useCase.Shorten(c.Request().Context(), inDTO)

	if err != nil {
		c.Response().WriteHeader(http.StatusInternalServerError)
		return err
	}

	type respItem struct {
		CorrelationID string `json:"correlation_id"`
		ShortURL      string `json:"short_url"`
	}

	response := make([]respItem, 0, len(outDTO.Items))
	for _, item := range outDTO.Items {
		response = append(response, respItem{
			CorrelationID: item.CorrelationID,
			ShortURL:      item.ShortURL,
		})
	}

	c.Response().Header().Set("Content-Type", "application/json")
	c.Response().WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(c.Response().Writer).Encode(response); err != nil {
		c.Response().WriteHeader(http.StatusInternalServerError)
		return err
	}

	return nil
}
