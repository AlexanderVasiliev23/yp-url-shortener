package app

import (
	"context"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/entrypoints/http/add"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/entrypoints/http/get"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/entrypoints/http/ping"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/entrypoints/http/shorten"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/entrypoints/http/shorten/batch"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/entrypoints/http/stats"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/entrypoints/http/user/urls/deleteurl"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/entrypoints/http/user/urls/list"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/middlewares/gzip"
	add2 "github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/usecases/add"
	get2 "github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/usecases/get"
	batch2 "github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/usecases/shorten/batch"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/usecases/shorten/single"
	stats2 "github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/usecases/stats"
	delete2 "github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/usecases/user/url/delete"
	list2 "github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/usecases/user/url/list"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/middlewares/jwt"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/middlewares/logger"
)

func (a *App) configureRouter(ctx context.Context) *echo.Echo {
	e := echo.New()

	e.Use(
		logger.Middleware(),
		// касмтомная gzip middleware потребляет очень много памяти, но необходима для прохождения тестов
		// для того, чтобы лучше видеть потребление памяти, при замерах отключал этот middleware и включал middleware из echo
		gzip.Middleware(),
		//middleware.Gzip(),
		middleware.Recover(),
		jwt.Middleware(ctx, a.conf.JWTSecretKey),
	)

	addUseCase := add2.NewUseCase(a.storage, a.tokenGenerator, a.userContextFetcher, a.conf.BaseAddress)
	getUseCase := get2.NewUseCase(a.storage)
	singleUseCase := single.NewUseCase(a.storage, a.tokenGenerator, a.userContextFetcher, a.conf.BaseAddress)
	batchUseCase := batch2.NewUseCase(a.storage, a.tokenGenerator, a.uuidGenerator, a.userContextFetcher, a.conf.BaseAddress)
	listUseCase := list2.NewUseCase(a.storage, a.userContextFetcher, a.conf.BaseAddress)
	deleteUseCase := delete2.NewUseCase(a.storage, a.userContextFetcher, a.deleteByTokenCh)
	statsUseCase := stats2.NewUseCase(a.storage, a.conf.TrustedSubnet)

	addHandler := add.NewHandler(addUseCase)
	getHandler := get.NewHandler(getUseCase)
	shortener := shorten.NewShortener(singleUseCase)
	batchShortener := batch.NewShortener(batchUseCase)
	pingHandler := ping.NewHandler(a.dbConn)
	listHandler := list.NewHandler(listUseCase)
	deleteHandler := deleteurl.NewHandler(deleteUseCase)
	statsHandler := stats.NewHandler(statsUseCase)

	e.GET("/:token", getHandler.Get)
	e.GET("/ping", pingHandler.Ping)
	e.POST("/", addHandler.Add)
	e.POST("/api/shorten", shortener.Handle)
	e.POST("/api/shorten/batch", batchShortener.Handle)
	e.GET("/api/internal/stats", statsHandler.Handle)

	g := e.Group("/api/user", jwt.Auth(a.conf.JWTSecretKey))
	g.GET("/urls", listHandler.List)
	g.DELETE("/urls", deleteHandler.Delete)

	return e
}
