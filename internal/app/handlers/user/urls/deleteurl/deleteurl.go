package deleteurl

import (
	"context"
	"encoding/json"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/workers/deleter"
	"github.com/labstack/echo/v4"
	"net/http"
)

type linksStorage interface {
	FilterOnlyThisUserTokens(ctx context.Context, userID int, tokens []string) ([]string, error)
}

type userContextFetcher interface {
	GetUserIDFromContext(ctx context.Context) (int, error)
}

func Delete(linksStorage linksStorage, userContextFetcher userContextFetcher, deleteByTokenCh chan<- deleter.DeleteTask) echo.HandlerFunc {
	return func(c echo.Context) error {
		userID, err := userContextFetcher.GetUserIDFromContext(c.Request().Context())
		if err != nil {
			c.Response().WriteHeader(http.StatusUnauthorized)
			return err
		}

		var reqBody []string
		if err := json.NewDecoder(c.Request().Body).Decode(&reqBody); err != nil {
			c.Response().WriteHeader(http.StatusBadRequest)
			return err
		}

		tokens, err := linksStorage.FilterOnlyThisUserTokens(c.Request().Context(), userID, reqBody)
		if err != nil {
			c.Response().WriteHeader(http.StatusInternalServerError)
			return err
		}

		deleteByTokenCh <- deleter.DeleteTask{
			Tokens: tokens,
		}

		c.Response().WriteHeader(http.StatusAccepted)

		return nil
	}
}
