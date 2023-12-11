package add

import (
	"context"
	"errors"
	"fmt"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/storage"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
)

type repository interface {
	Add(ctx context.Context, token, url string) error
	GetTokenByURL(ctx context.Context, url string) (string, error)
}

type tokenGenerator interface {
	Generate() (string, error)
}

func Add(repository repository, tokenGenerator tokenGenerator, addr string) echo.HandlerFunc {
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

		if err := repository.Add(c.Request().Context(), token, string(url)); err != nil {
			if errors.Is(err, storage.ErrAlreadyExists) {
				token, err := repository.GetTokenByURL(c.Request().Context(), string(url))
				if err != nil {
					c.Response().WriteHeader(http.StatusInternalServerError)
					return nil
				}

				c.Response().WriteHeader(http.StatusConflict)
				_, _ = fmt.Fprintf(c.Response(), "%s/%s", addr, token)
				return nil
			}
			c.Response().WriteHeader(http.StatusInternalServerError)
			return nil
		}

		c.Response().WriteHeader(http.StatusCreated)
		_, _ = fmt.Fprintf(c.Response(), "%s/%s", addr, token)

		return nil
	}
}
