package add

import (
	"fmt"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
)

type repository interface {
	Add(token, url string) error
}

type tokenGenerator interface {
	Generate() string
}

func Add(repository repository, tokenGenerator tokenGenerator, addr string) echo.HandlerFunc {
	return func(c echo.Context) error {
		url, err := io.ReadAll(c.Request().Body)
		if err != nil || len(url) == 0 {
			c.Response().WriteHeader(http.StatusBadRequest)
			return nil
		}

		token := tokenGenerator.Generate()
		if err := repository.Add(token, string(url)); err != nil {
			c.Response().WriteHeader(http.StatusInternalServerError)
			return nil
		}

		c.Response().WriteHeader(http.StatusCreated)
		_, _ = fmt.Fprintf(c.Response(), "%s/%s", addr, token)

		return nil
	}
}
