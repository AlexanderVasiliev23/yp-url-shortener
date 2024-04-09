package get

import (
	"context"
	"errors"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/usecases/get"
	"net/http"

	"github.com/labstack/echo/v4"
)

type useCase interface {
	Get(ctx context.Context, token string) (originalURL string, err error)
}

// Handler missing godoc.
type Handler struct {
	useCase useCase
}

func NewHandler(useCase useCase) *Handler {
	return &Handler{useCase: useCase}
}

// Get missing godoc.
func (h *Handler) Get(c echo.Context) error {
	token := c.Param("token")

	originalURL, err := h.useCase.Get(context.Background(), token)

	if err != nil {
		if errors.Is(err, get.ErrTokenIsEmpty) {
			c.Response().WriteHeader(http.StatusBadRequest)
			return err
		}

		if errors.Is(err, get.ErrTokenIsDeleted) {
			c.Response().WriteHeader(http.StatusGone)
			return err
		}

		c.Response().WriteHeader(http.StatusInternalServerError)
		return err
	}

	c.Response().Header().Set("Location", originalURL)
	c.Response().WriteHeader(http.StatusTemporaryRedirect)

	return nil
}
