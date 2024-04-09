package stats

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/usecases/stats"
	iputil "github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/util/ip"
	"github.com/labstack/echo/v4"
	"net/http"
)

type useCase interface {
	Stats(ctx context.Context, ip string) (*stats.OutDTO, error)
}

// Handler missing godoc.
type Handler struct {
	useCase useCase
}

// NewHandler missing godoc.
func NewHandler(useCase useCase) *Handler {
	return &Handler{useCase: useCase}
}

// Handle missing godoc.
func (h *Handler) Handle(c echo.Context) error {
	ip := iputil.IPFromRequest(c.Request())
	statsOutDTO, err := h.useCase.Stats(c.Request().Context(), ip)
	if err != nil {
		if errors.Is(err, stats.ErrNotTrustedIP) {
			c.Response().WriteHeader(http.StatusForbidden)
			return nil
		}

		c.Response().WriteHeader(http.StatusInternalServerError)
		return err
	}

	resp := struct {
		Urls  int `json:"urls"`
		Users int `json:"users"`
	}{
		Urls:  statsOutDTO.Urls,
		Users: statsOutDTO.Users,
	}

	c.Response().Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(c.Response()).Encode(resp); err != nil {
		c.Response().WriteHeader(http.StatusInternalServerError)
		return err
	}

	return nil
}
