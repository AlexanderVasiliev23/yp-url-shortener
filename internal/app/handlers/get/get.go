package get

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type repository interface {
	Get(string) (url string, err error)
}

func Get(repository repository) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.Param("token")

		if token == "" {
			c.Response().WriteHeader(http.StatusBadRequest)
			return nil
		}

		url, err := repository.Get(token)
		if err != nil {
			c.Response().WriteHeader(http.StatusBadRequest)
			return nil
		}

		c.Response().Header().Set("Location", url)
		c.Response().WriteHeader(http.StatusTemporaryRedirect)

		return nil
	}
}