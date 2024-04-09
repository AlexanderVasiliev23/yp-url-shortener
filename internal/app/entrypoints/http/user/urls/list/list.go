package list

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/usecases/user/url/list"
	"github.com/labstack/echo/v4"
	"net/http"
)

type useCase interface {
	List(ctx context.Context) (*list.OutDTO, error)
}

// Handler missing godoc.
type Handler struct {
	useCase useCase
}

// NewHandler missing godoc.
func NewHandler(useCase useCase) *Handler {
	return &Handler{
		useCase: useCase,
	}
}

// List missing godoc.
func (h *Handler) List(c echo.Context) error {
	outDTO, err := h.useCase.List(c.Request().Context())

	if err != nil {
		if errors.Is(err, list.ErrUnauthorized) {
			c.Response().WriteHeader(http.StatusUnauthorized)
			return err
		}

		if errors.Is(err, list.ErrNoSavedURLs) {
			c.Response().WriteHeader(http.StatusNoContent)
			return nil
		}

		c.Response().WriteHeader(http.StatusInternalServerError)
		return err
	}

	type respItem struct {
		ShortURL    string `json:"short_url"`
		OriginalURL string `json:"original_url"`
	}

	type resp []respItem

	response := make(resp, 0, len(outDTO.Items))

	for _, outDTOItem := range outDTO.Items {
		response = append(response, respItem{
			ShortURL:    outDTOItem.ShortURL,
			OriginalURL: outDTOItem.OriginalURL,
		})
	}

	c.Response().Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(c.Response().Writer).Encode(response); err != nil {
		c.Response().WriteHeader(http.StatusInternalServerError)
		return err
	}

	return nil
}
