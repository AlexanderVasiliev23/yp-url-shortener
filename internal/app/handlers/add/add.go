package add

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/models"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/storage"
)

type repository interface {
	Add(ctx context.Context, shortLink *models.ShortLink) error
	GetTokenByURL(ctx context.Context, url string) (string, error)
}

type tokenGenerator interface {
	Generate() (string, error)
}

type userContextFetcher interface {
	GetUserIDFromContext(ctx context.Context) (int, error)
}

// Handler missing godoc.
type Handler struct {
	repository         repository
	tokenGenerator     tokenGenerator
	userContextFetcher userContextFetcher
	addr               string
}

// NewHandler missing godoc.
func NewHandler(
	repository repository,
	tokenGenerator tokenGenerator,
	userContextFetcher userContextFetcher,
	addr string,
) *Handler {
	return &Handler{
		repository:         repository,
		tokenGenerator:     tokenGenerator,
		userContextFetcher: userContextFetcher,
		addr:               addr,
	}
}

// Add missing godoc.
func (h *Handler) Add(c echo.Context) error {
	url, err := io.ReadAll(c.Request().Body)
	if err != nil || len(url) == 0 {
		c.Response().WriteHeader(http.StatusBadRequest)
		return nil
	}

	token, err := h.tokenGenerator.Generate()
	if err != nil {
		c.Response().WriteHeader(http.StatusInternalServerError)
		return nil
	}

	userID, err := h.userContextFetcher.GetUserIDFromContext(c.Request().Context())
	if err != nil {
		c.Response().WriteHeader(http.StatusInternalServerError)
		return err
	}

	model := models.NewShortLink(userID, uuid.New(), token, string(url))
	if err := h.repository.Add(c.Request().Context(), model); err != nil {
		if !errors.Is(err, storage.ErrAlreadyExists) {
			c.Response().WriteHeader(http.StatusInternalServerError)
			return nil
		}

		token, err := h.repository.GetTokenByURL(c.Request().Context(), string(url))
		if err != nil {
			c.Response().WriteHeader(http.StatusInternalServerError)
			return nil
		}

		c.Response().WriteHeader(http.StatusConflict)
		_, _ = fmt.Fprintf(c.Response(), "%s/%s", h.addr, token)

		return nil
	}

	c.Response().WriteHeader(http.StatusCreated)
	_, _ = fmt.Fprintf(c.Response(), "%s/%s", h.addr, token)

	return nil
}
