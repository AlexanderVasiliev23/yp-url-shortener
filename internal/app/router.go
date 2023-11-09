package app

import (
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/handlers/add"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/handlers/get"
	"github.com/labstack/echo/v4"
)

func (a *app) configureRouter() *echo.Echo {
	e := echo.New()

	e.POST("/", add.Add(a.localStorage, a.tokenGenerator, a.conf.BaseAddress))
	e.GET("/:token", get.Get(a.localStorage))

	return e
}
