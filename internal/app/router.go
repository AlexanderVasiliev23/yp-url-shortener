package app

import (
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/handlers/add"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/handlers/get"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/handlers/ping"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/handlers/shorten"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/handlers/shorten/batch"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/middlewares/gzip"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/middlewares/logger"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func (a *App) configureRouter() *echo.Echo {
	e := echo.New()

	e.Use(
		logger.Middleware(),
		gzip.Middleware(),
		middleware.Recover(),
	)

	e.POST("/", add.Add(a.storage, a.tokenGenerator, a.conf.BaseAddress))
	e.GET("/:token", get.Get(a.storage))
	e.POST("/api/shorten", shorten.Shorten(a.storage, a.tokenGenerator, a.conf.BaseAddress))
	e.POST("/api/shorten/batch", batch.Shorten(a.storage, a.tokenGenerator, a.conf.BaseAddress))
	e.GET("/ping", ping.Ping(a.dbConn))

	return e
}
