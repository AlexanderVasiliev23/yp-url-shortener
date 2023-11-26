package logger

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"time"
)

func Middleware(logger *zap.SugaredLogger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			err := next(c)

			duration := time.Since(start)

			logger.Infow(
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
