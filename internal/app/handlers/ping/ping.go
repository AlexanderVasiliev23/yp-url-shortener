package ping

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

// Handler missing godoc.
type Handler struct {
	dbConn *pgxpool.Pool
}

// NewHandler missing godoc.
func NewHandler(dbConn *pgxpool.Pool) *Handler {
	return &Handler{dbConn: dbConn}
}

// Ping missing godoc.
func (h *Handler) Ping(c echo.Context) error {
	if h.dbConn == nil {
		c.Response().WriteHeader(http.StatusInternalServerError)
		return nil
	}

	if err := h.dbConn.Ping(c.Request().Context()); err != nil {
		c.Response().WriteHeader(http.StatusInternalServerError)
		return err
	}

	return nil
}
