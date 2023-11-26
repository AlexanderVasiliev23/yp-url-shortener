package app

import (
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/handlers/add"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/handlers/get"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/handlers/shorten"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/middlewares/gzip"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/middlewares/logger"
	"github.com/labstack/echo/v4"
)

func (a *App) configureRouter() *echo.Echo {
	e := echo.New()

	e.Use(logger.Middleware(a.logger))
	e.Use(gzip.Middleware())

	e.POST("/", add.Add(a.localStorage, a.tokenGenerator, a.conf.BaseAddress))
	e.GET("/:token", get.Get(a.localStorage))
	e.POST("/api/shorten", shorten.Shorten(a.localStorage, a.tokenGenerator, a.conf.BaseAddress))

	return e
}
