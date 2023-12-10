package get

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
)

type repository interface {
	Get(ctx context.Context, token string) (url string, err error)
}

func Get(repository repository) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Param("token")

		if token == "" {
			c.Response().WriteHeader(http.StatusBadRequest)
			return nil
		}

		url, err := repository.Get(c.Request().Context(), token)
		if err != nil {
			c.Response().WriteHeader(http.StatusBadRequest)
			return nil
		}

		c.Response().Header().Set("Location", url)
		c.Response().WriteHeader(http.StatusTemporaryRedirect)

		return nil
	}
}
