package urls

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/models"
	"github.com/labstack/echo/v4"
	"net/http"
)

type linksStorage interface {
	FindByUserID(ctx context.Context, userID int) ([]*models.ShortLink, error)
}

type userContextFetcher interface {
	GetUserIDFromContext(ctx context.Context) (int, error)
}

func Urls(storage linksStorage, userContextFetcher userContextFetcher, addr string) echo.HandlerFunc {
	return func(c echo.Context) error {
		userID, err := userContextFetcher.GetUserIDFromContext(c.Request().Context())
		if err != nil {
			c.Response().WriteHeader(http.StatusUnauthorized)
			return err
		}

		shortLinks, err := storage.FindByUserID(c.Request().Context(), userID)
		if err != nil {
			c.Response().WriteHeader(http.StatusInternalServerError)
			return err
		}

		if len(shortLinks) == 0 {
			c.Response().WriteHeader(http.StatusNoContent)
			return nil
		}

		type respItem struct {
			ShortUrl    string `json:"short_url"`
			OriginalUrl string `json:"original_url"`
		}

		type resp []respItem

		response := make(resp, 0, len(shortLinks))

		for _, shortLink := range shortLinks {
			response = append(response, respItem{
				ShortUrl:    fmt.Sprintf("%s/%s", addr, shortLink.Token),
				OriginalUrl: shortLink.Original,
			})
		}

		c.Response().Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(c.Response().Writer).Encode(response); err != nil {
			c.Response().WriteHeader(http.StatusInternalServerError)
			return err
		}

		return nil
	}
}
