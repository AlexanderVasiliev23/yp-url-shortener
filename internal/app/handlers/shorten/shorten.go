package shorten

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/logger"
	"github.com/labstack/echo/v4"
	"net/http"
)

var (
	ErrURLIsEmpty = errors.New("url is empty")
)

type repository interface {
	Add(ctx context.Context, token, url string) error
}

type tokenGenerator interface {
	Generate() (string, error)
}

func Shorten(repository repository, tokenGenerator tokenGenerator, addr string) echo.HandlerFunc {
	return func(c echo.Context) error {
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

		token, err := tokenGenerator.Generate()
		if err != nil {
			logger.Log.Errorf("tokenGenerator.Generate: %s", err)
			c.Response().WriteHeader(http.StatusInternalServerError)
			return err
		}

		if err := repository.Add(c.Request().Context(), token, req.URL); err != nil {
			logger.Log.Errorf("repository.Add: %s", err)
			c.Response().WriteHeader(http.StatusInternalServerError)
			return err
		}

		resp := struct {
			Result string `json:"result"`
		}{
			Result: fmt.Sprintf("%s/%s", addr, token),
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
