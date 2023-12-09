package app

import (
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/handlers/add"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/handlers/get"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/handlers/ping"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/handlers/shorten"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/middlewares/gzip"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/middlewares/logger"
	"github.com/labstack/echo/v4"
)

func (a *App) configureRouter() *echo.Echo {
	e := echo.New()

	e.Use(
		logger.Middleware(),
		gzip.Middleware(),
	)

	e.POST("/", add.Add(a.storage, a.tokenGenerator, a.conf.BaseAddress))
	e.GET("/:token", get.Get(a.storage))
	e.POST("/api/shorten", shorten.Shorten(a.storage, a.tokenGenerator, a.conf.BaseAddress))
	e.GET("/ping", ping.Ping(a.dbConn))

	return e
}
