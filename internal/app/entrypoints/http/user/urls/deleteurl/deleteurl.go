package deleteurl

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/workers/deleter"
)

type linksStorage interface {
	FilterOnlyThisUserTokens(ctx context.Context, userID int, tokens []string) ([]string, error)
}

type userContextFetcher interface {
	GetUserIDFromContext(ctx context.Context) (int, error)
}

// Handler missing godoc.
type Handler struct {
	linksStorage       linksStorage
	userContextFetcher userContextFetcher
	deleteByTokenCh    chan<- deleter.DeleteTask
}

// NewHandler missing godoc.
func NewHandler(
	linksStorage linksStorage,
	userContextFetcher userContextFetcher,
	deleteByTokenCh chan<- deleter.DeleteTask,
) *Handler {
	return &Handler{
		linksStorage:       linksStorage,
		userContextFetcher: userContextFetcher,
		deleteByTokenCh:    deleteByTokenCh,
	}
}

// Delete missing godoc.
func (h *Handler) Delete(c echo.Context) error {
	userID, err := h.userContextFetcher.GetUserIDFromContext(c.Request().Context())
	if err != nil {
		c.Response().WriteHeader(http.StatusUnauthorized)
		return err
	}

	var reqBody []string
	if _err := json.NewDecoder(c.Request().Body).Decode(&reqBody); _err != nil {
		c.Response().WriteHeader(http.StatusBadRequest)
		return _err
	}

	tokens, err := h.linksStorage.FilterOnlyThisUserTokens(c.Request().Context(), userID, reqBody)
	if err != nil {
		c.Response().WriteHeader(http.StatusInternalServerError)
		return err
	}

	h.deleteByTokenCh <- deleter.DeleteTask{
		Tokens: tokens,
	}

	c.Response().WriteHeader(http.StatusAccepted)

	return nil
}
