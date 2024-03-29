package logger

import (
	"time"

	"github.com/labstack/echo/v4"

	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/logger"
)

// Middleware missing godoc.
func Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			err := next(c)

			if err != nil {
				logger.Log.Error(err)
			}

			duration := time.Since(start)

			logger.Log.Infow(
				"request handled",
				"method", c.Request().Method,
				"uri", c.Request().RequestURI,
				"duration", duration,
				"status", c.Response().Status,
				"size", c.Response().Size,
			)

			return err
		}
	}
}
