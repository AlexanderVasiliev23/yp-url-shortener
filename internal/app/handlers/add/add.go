package add

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
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
		_, _ = fmt.Fprintf(c.Response(), "http://%s/%s", addr, token)

		return nil
	}
}
