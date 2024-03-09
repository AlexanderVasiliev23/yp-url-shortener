package shorten

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/models"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/storage"
)

// ErrURLIsEmpty missing godoc.
var (
	ErrURLIsEmpty = errors.New("url is empty")
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

// Shortener missing godoc.
type Shortener struct {
	repository         repository
	tokenGenerator     tokenGenerator
	userContextFetcher userContextFetcher
	addr               string
}

// NewShortener missing godoc.
func NewShortener(
	repository repository,
	tokenGenerator tokenGenerator,
	userContextFetcher userContextFetcher,
	addr string,
) *Shortener {
	return &Shortener{
		repository:         repository,
		tokenGenerator:     tokenGenerator,
		userContextFetcher: userContextFetcher,
		addr:               addr,
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

	if req.URL == "" {
		c.Response().WriteHeader(http.StatusBadRequest)
		return ErrURLIsEmpty
	}

	token, err := h.tokenGenerator.Generate()
	if err != nil {
		c.Response().WriteHeader(http.StatusInternalServerError)
		return err
	}

	type resp struct {
		Result string `json:"result"`
	}

	c.Response().Header().Set("Content-Type", "application/json")

	userID, err := h.userContextFetcher.GetUserIDFromContext(c.Request().Context())
	if err != nil {
		c.Response().WriteHeader(http.StatusInternalServerError)
		return err
	}
	model := models.NewShortLink(userID, uuid.New(), token, req.URL)
	if err := h.repository.Add(c.Request().Context(), model); err != nil {
		if !errors.Is(err, storage.ErrAlreadyExists) {
			c.Response().WriteHeader(http.StatusInternalServerError)
			return err
		}

		_token, err := h.repository.GetTokenByURL(c.Request().Context(), req.URL)
		if err != nil {
			c.Response().WriteHeader(http.StatusInternalServerError)
			return err
		}

		c.Response().WriteHeader(http.StatusConflict)
		response := resp{Result: h.addr + "/" + _token}
		if err := json.NewEncoder(c.Response().Writer).Encode(response); err != nil {
			c.Response().WriteHeader(http.StatusInternalServerError)
			return err
		}

		return nil
	}

	response := resp{Result: h.addr + "/" + token}
	c.Response().WriteHeader(http.StatusCreated)
	c.Response().Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(c.Response().Writer).Encode(response); err != nil {
		c.Response().WriteHeader(http.StatusInternalServerError)
		return err
	}

	return nil
}
