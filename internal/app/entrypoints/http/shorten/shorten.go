package shorten

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/usecases/shorten/single"
	"github.com/labstack/echo/v4"
	"net/http"
)

type useCase interface {
	Shorten(ctx context.Context, jsonString string) (shortURL string, err error)
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
	req := struct {
		URL string
	}{}

	if err := json.NewDecoder(c.Request().Body).Decode(&req); err != nil {
		c.Response().WriteHeader(http.StatusBadRequest)
		return err
	}

	shortURL, err := h.useCase.Shorten(c.Request().Context(), req.URL)

	type resp struct {
		Result string `json:"result"`
	}

	if err != nil {
		if errors.Is(err, single.ErrInvalidJSON) || errors.Is(err, single.ErrEmptyOriginalURL) {
			c.Response().WriteHeader(http.StatusBadRequest)
			return err
		}

		if errors.Is(err, single.ErrAlreadyExists) {
			c.Response().WriteHeader(http.StatusConflict)
			response := resp{Result: shortURL}
			c.Response().Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(c.Response().Writer).Encode(response); err != nil {
				c.Response().WriteHeader(http.StatusInternalServerError)
				return err
			}
		}

		c.Response().WriteHeader(http.StatusInternalServerError)
		return err
	}

	c.Response().WriteHeader(http.StatusCreated)
	response := resp{Result: shortURL}
	c.Response().Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(c.Response().Writer).Encode(response); err != nil {
		c.Response().WriteHeader(http.StatusInternalServerError)
		return err
	}

	return nil
}
