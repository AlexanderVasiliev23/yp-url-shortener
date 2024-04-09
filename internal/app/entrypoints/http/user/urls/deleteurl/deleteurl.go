package deleteurl

import (
	"context"
	"encoding/json"
	"errors"
	deleteusecase "github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/usecases/user/url/delete"
	"net/http"

	"github.com/labstack/echo/v4"
)

type useCase interface {
	Delete(ctx context.Context, tokens []string) error
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

// Delete missing godoc.
func (h *Handler) Delete(c echo.Context) error {
	var tokens []string
	if _err := json.NewDecoder(c.Request().Body).Decode(&tokens); _err != nil {
		c.Response().WriteHeader(http.StatusBadRequest)
		return _err
	}

	err := h.useCase.Delete(c.Request().Context(), tokens)

	if err != nil {
		if errors.Is(err, deleteusecase.ErrUnauthorized) {
			c.Response().WriteHeader(http.StatusUnauthorized)
			return err
		}

		c.Response().WriteHeader(http.StatusInternalServerError)
		return err
	}

	c.Response().WriteHeader(http.StatusAccepted)

	return nil
}
