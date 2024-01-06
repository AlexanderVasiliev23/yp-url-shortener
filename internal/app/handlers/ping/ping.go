package ping

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"net/http"
)

type Handler struct {
	dbConn *pgxpool.Pool
}

func NewHandler(dbConn *pgxpool.Pool) *Handler {
	return &Handler{dbConn: dbConn}
}

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
