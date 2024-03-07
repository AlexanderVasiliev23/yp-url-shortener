package app

import (
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/middlewares/gzip"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/handlers/add"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/handlers/get"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/handlers/ping"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/handlers/shorten"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/handlers/shorten/batch"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/handlers/user/urls/deleteurl"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/handlers/user/urls/list"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/middlewares/jwt"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/middlewares/logger"
)

func (a *App) configureRouter() *echo.Echo {
	e := echo.New()

	e.Use(
		logger.Middleware(),
		// касмтомная gzip middleware потребляет очень много памяти, но необходима для прохождения тестов
		// для того, чтобы лучше видеть потребление памяти, при замерах отключал этот middleware и включал middleware из echo
		gzip.Middleware(),
		//middleware.Gzip(),
		middleware.Recover(),
		jwt.Middleware(a.conf.JWTSecretKey),
	)

	addHandler := add.NewHandler(a.storage, a.tokenGenerator, a.userContextFetcher, a.conf.BaseAddress)
	getHandler := get.NewHandler(a.storage)
	shortener := shorten.NewShortener(a.storage, a.tokenGenerator, a.userContextFetcher, a.conf.BaseAddress)
	batchShortener := batch.NewShortener(a.storage, a.tokenGenerator, a.uuidGenerator, a.userContextFetcher, a.conf.BaseAddress)
	pingHandler := ping.NewHandler(a.dbConn)
	listHandler := list.NewHandler(a.storage, a.userContextFetcher, a.conf.BaseAddress)
	deleteHandler := deleteurl.NewHandler(a.storage, a.userContextFetcher, a.deleteByTokenCh)

	e.GET("/:token", getHandler.Get)
	e.GET("/ping", pingHandler.Ping)
	e.POST("/", addHandler.Add)
	e.POST("/api/shorten", shortener.Handle)
	e.POST("/api/shorten/batch", batchShortener.Handle)

	g := e.Group("/api/user", jwt.Auth(a.conf.JWTSecretKey))
	g.GET("/urls", listHandler.List)
	g.DELETE("/urls", deleteHandler.Delete)

	return e
}
