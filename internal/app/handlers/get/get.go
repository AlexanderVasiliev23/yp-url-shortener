package get

import (
	"context"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

type repository interface {
	Get(ctx context.Context, token string) (*models.ShortLink, error)
}

type Handler struct {
	repository repository
}

func NewHandler(repository repository) *Handler {
	return &Handler{repository: repository}
}

func (h *Handler) Get(c echo.Context) error {
	token := c.Param("token")

	if token == "" {
		c.Response().WriteHeader(http.StatusBadRequest)
		return nil
	}

	shortLink, err := h.repository.Get(c.Request().Context(), token)
	if err != nil {
		c.Response().WriteHeader(http.StatusBadRequest)
		return nil
	}
	if shortLink.DeletedAt != nil {
		c.Response().WriteHeader(http.StatusGone)
		return nil
	}

	c.Response().Header().Set("Location", shortLink.Original)
	c.Response().WriteHeader(http.StatusTemporaryRedirect)

	return nil
}
