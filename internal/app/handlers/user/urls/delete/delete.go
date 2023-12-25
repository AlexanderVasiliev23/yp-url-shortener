package delete

import (
	"context"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"net/http"
)

type linksStorage interface {
	DeleteTokens(ctx context.Context, userID int, tokens []string) error
}

type userContextFetcher interface {
	GetUserIDFromContext(ctx context.Context) (int, error)
}

func Delete(storage linksStorage, userContextFetcher userContextFetcher) echo.HandlerFunc {
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

		deleteTokensByUser(userID, reqBody, storage)

		c.Response().WriteHeader(http.StatusAccepted)

		return nil
	}
}

func deleteTokensByUser(userId int, tokens []string, storage linksStorage) {

}
