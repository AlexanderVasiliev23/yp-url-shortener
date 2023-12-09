package ping

import (
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"net/http"
)

func Ping(dbConn *pgx.Conn) echo.HandlerFunc {
	return func(c echo.Context) error {
		if err := dbConn.Ping(c.Request().Context()); err != nil {
			c.Response().WriteHeader(http.StatusInternalServerError)
			return nil
		}

		return nil
	}
}
