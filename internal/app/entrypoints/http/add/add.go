package add

import (
	"context"
	"errors"
	"fmt"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/usecases/add"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
)

type useCase interface {
	Add(ctx context.Context, originalURL string) (shortURL string, err error)
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

// Add missing godoc.
func (h *Handler) Add(c echo.Context) error {
	originalURL, err := io.ReadAll(c.Request().Body)
	if err != nil {
		c.Response().WriteHeader(http.StatusBadRequest)
		return nil
	}

	shortURL, err := h.useCase.Add(c.Request().Context(), string(originalURL))

	if err != nil {
		if errors.Is(err, add.ErrOriginURLAlreadyExists) {
			c.Response().WriteHeader(http.StatusConflict)
			_, _ = fmt.Fprint(c.Response(), shortURL)
			return nil
		}

		if errors.Is(err, add.ErrOriginalURLIsEmpty) {
			c.Response().WriteHeader(http.StatusBadRequest)
			return nil
		}

		c.Response().WriteHeader(http.StatusInternalServerError)
		return nil
	}

	c.Response().WriteHeader(http.StatusCreated)
	_, _ = fmt.Fprint(c.Response(), shortURL)

	return nil
}
