package list

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/models"
)

type linksStorage interface {
	FindByUserID(ctx context.Context, userID int) ([]*models.ShortLink, error)
}

type userContextFetcher interface {
	GetUserIDFromContext(ctx context.Context) (int, error)
}

// Handler missing godoc.
type Handler struct {
	storage            linksStorage
	userContextFetcher userContextFetcher
	addr               string
}

// NewHandler missing godoc.
func NewHandler(
	storage linksStorage,
	userContextFetcher userContextFetcher,
	addr string,
) *Handler {
	return &Handler{
		storage:            storage,
		userContextFetcher: userContextFetcher,
		addr:               addr,
	}
}

// List missing godoc.
func (h *Handler) List(c echo.Context) error {
	userID, err := h.userContextFetcher.GetUserIDFromContext(c.Request().Context())
	if err != nil {
		c.Response().WriteHeader(http.StatusUnauthorized)
		return err
	}

	shortLinks, err := h.storage.FindByUserID(c.Request().Context(), userID)
	if err != nil {
		c.Response().WriteHeader(http.StatusInternalServerError)
		return err
	}

	if len(shortLinks) == 0 {
		c.Response().WriteHeader(http.StatusNoContent)
		return nil
	}

	type respItem struct {
		ShortURL    string `json:"short_url"`
		OriginalURL string `json:"original_url"`
	}

	type resp []respItem

	response := make(resp, 0, len(shortLinks))

	for _, shortLink := range shortLinks {
		response = append(response, respItem{
			ShortURL:    fmt.Sprintf("%s/%s", h.addr, shortLink.Token),
			OriginalURL: shortLink.Original,
		})
	}

	c.Response().Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(c.Response().Writer).Encode(response); err != nil {
		c.Response().WriteHeader(http.StatusInternalServerError)
		return err
	}

	return nil
}
