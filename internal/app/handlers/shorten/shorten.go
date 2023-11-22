package shorten

import (
	"encoding/json"
	"errors"
	"github.com/labstack/echo/v4"
	"net/http"
)

var (
	ErrURLIsEmpty = errors.New("url is empty")
)

type repository interface {
	Add(token, url string) error
}

type tokenGenerator interface {
	Generate() (string, error)
}

func Shorten(repository repository, tokenGenerator tokenGenerator) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := struct {
			Url string
		}{}

		if err := json.NewDecoder(c.Request().Body).Decode(&req); err != nil {
			c.Response().WriteHeader(http.StatusBadRequest)
			return err
		}

		if req.Url == "" {
			c.Response().WriteHeader(http.StatusBadRequest)
			return ErrURLIsEmpty
		}

		token, err := tokenGenerator.Generate()
		if err != nil {
			c.Response().WriteHeader(http.StatusInternalServerError)
			return err
		}

		if err := repository.Add(token, req.Url); err != nil {
			c.Response().WriteHeader(http.StatusInternalServerError)
			return err
		}

		resp := struct {
			Result string `json:"result"`
		}{
			Result: token,
		}

		c.Response().Header().Set("Content-Type", "application/json")
		c.Response().WriteHeader(http.StatusCreated)
		if err := json.NewEncoder(c.Response().Writer).Encode(resp); err != nil {
			c.Response().WriteHeader(http.StatusInternalServerError)
			return err
		}

		return nil
	}
}
