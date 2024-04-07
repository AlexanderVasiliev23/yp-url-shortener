package stats

import (
	"context"
	"encoding/json"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/storage"
	"github.com/labstack/echo/v4"
	"net/http"
)

type repository interface {
	Stats(ctx context.Context) (*storage.StatsOutDTO, error)
}

type Handler struct {
	repo          repository
	trustedSubnet string
}

func NewHandler(repo repository, trustedSubnet string) *Handler {
	return &Handler{repo: repo, trustedSubnet: trustedSubnet}
}

func (h *Handler) Handle(c echo.Context) error {
	if h.trustedSubnet == "" {
		c.Response().WriteHeader(http.StatusForbidden)
		return nil
	}

	stats, err := h.repo.Stats(c.Request().Context())
	if err != nil {
		c.Response().WriteHeader(http.StatusInternalServerError)
		return err
	}

	resp := struct {
		Urls  int `json:"urls"`
		Users int `json:"users"`
	}{
		Urls:  stats.UrlsCount,
		Users: stats.UsersCount,
	}

	c.Response().Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(c.Response()).Encode(resp); err != nil {
		c.Response().WriteHeader(http.StatusInternalServerError)
		return err
	}

	return nil
}
