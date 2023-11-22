package app

import (
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/handlers/add"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/handlers/get"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/handlers/shorten"
	"github.com/labstack/echo/v4"
	"time"
)

func (a *App) configureRouter() *echo.Echo {
	e := echo.New()

	e.Use(a.loggerMiddleware())

	e.POST("/", add.Add(a.localStorage, a.tokenGenerator, a.conf.BaseAddress))
	e.GET("/:token", get.Get(a.localStorage))
	e.POST("/api/shorten", shorten.Shorten(a.localStorage, a.tokenGenerator, a.conf.BaseAddress))

	return e
}

func (a *App) loggerMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			err := next(c)

			duration := time.Since(start)

			a.logger.Infow(
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
