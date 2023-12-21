package add

import (
	"context"
	"errors"
	"fmt"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/models"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/storage"
	"github.com/google/uuid"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
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

func Add(repository repository, tokenGenerator tokenGenerator, userContextFetcher userContextFetcher, addr string) echo.HandlerFunc {
	return func(c echo.Context) error {
		url, err := io.ReadAll(c.Request().Body)
		if err != nil || len(url) == 0 {
			c.Response().WriteHeader(http.StatusBadRequest)
			return nil
		}

		token, err := tokenGenerator.Generate()
		if err != nil {
			c.Response().WriteHeader(http.StatusInternalServerError)
			return nil
		}

		userID, err := userContextFetcher.GetUserIDFromContext(c.Request().Context())
		if err != nil {
			c.Response().WriteHeader(http.StatusInternalServerError)
			return err
		}

		model := models.NewShortLink(userID, uuid.New(), token, string(url))
		if err := repository.Add(c.Request().Context(), model); err != nil {
			if !errors.Is(err, storage.ErrAlreadyExists) {
				c.Response().WriteHeader(http.StatusInternalServerError)
				return nil
			}

			token, err := repository.GetTokenByURL(c.Request().Context(), string(url))
			if err != nil {
				c.Response().WriteHeader(http.StatusInternalServerError)
				return nil
			}

			c.Response().WriteHeader(http.StatusConflict)
			_, _ = fmt.Fprintf(c.Response(), "%s/%s", addr, token)

			return nil
		}

		c.Response().WriteHeader(http.StatusCreated)
		_, _ = fmt.Fprintf(c.Response(), "%s/%s", addr, token)

		return nil
	}
}
