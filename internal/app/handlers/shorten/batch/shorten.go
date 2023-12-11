package batch

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/logger"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/models"
	"github.com/labstack/echo/v4"
	"net/http"
)

type batchSaver interface {
	SaveBatch(ctx context.Context, shortLinks []*models.ShortLink) error
}

type tokenGenerator interface {
	Generate() (string, error)
}

func Shorten(saver batchSaver, tokenGenerator tokenGenerator, addr string) echo.HandlerFunc {
	return func(c echo.Context) error {
		type reqItem struct {
			CorrelationId string `json:"correlation_id"`
			OriginalURL   string `json:"original_url"`
		}

		type req []reqItem

		var requestItems req

		if err := json.NewDecoder(c.Request().Body).Decode(&requestItems); err != nil {
			c.Response().WriteHeader(http.StatusBadRequest)
			return err
		}

		type respItem struct {
			CorrelationId string `json:"correlation_id"`
			ShortURL      string `json:"short_url"`
		}

		type responseItems []respItem

		var response responseItems

		toSave := make([]*models.ShortLink, 0, len(requestItems))

		for _, requestItem := range requestItems {
			token, err := tokenGenerator.Generate()
			if err != nil {
				logger.Log.Errorf("tokenGenerator.Generate: %s", err)
				c.Response().WriteHeader(http.StatusInternalServerError)
				return err
			}

			shortLink := models.NewShortLink(token, requestItem.OriginalURL)
			toSave = append(toSave, shortLink)

			respItem := respItem{
				CorrelationId: requestItem.CorrelationId,
				ShortURL:      fmt.Sprintf("%s/%s", addr, token),
			}

			response = append(response, respItem)
		}

		if err := saver.SaveBatch(c.Request().Context(), toSave); err != nil {
			logger.Log.Errorf("saver.SaveBatch: %s", err)
			c.Response().WriteHeader(http.StatusInternalServerError)
			return err
		}

		c.Response().Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(c.Response().Writer).Encode(response); err != nil {
			logger.Log.Errorf("encode response: %s", err)
			c.Response().WriteHeader(http.StatusInternalServerError)
			return err
		}

		return nil
	}
}