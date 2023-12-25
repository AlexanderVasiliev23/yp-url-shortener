package ping

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"net/http"
)

func Ping(dbConn *pgxpool.Pool) echo.HandlerFunc {
	return func(c echo.Context) error {
		if dbConn == nil {
			c.Response().WriteHeader(http.StatusInternalServerError)
			return nil
		}

		if err := dbConn.Ping(c.Request().Context()); err != nil {
			c.Response().WriteHeader(http.StatusInternalServerError)
			return err
		}

		return nil
	}
}
