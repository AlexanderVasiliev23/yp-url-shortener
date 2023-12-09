package app

import (
	"context"
	"fmt"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/configs"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/logger"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/storage"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/storage/dumper"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/storage/local"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/tokengenerator"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"net/http"
)

type App struct {
	conf           *configs.Config
	storage        storage.Storage
	tokenGenerator *tokengenerator.TokenGenerator
	router         *echo.Echo
	dbConn         *pgx.Conn
}

func New(conf *configs.Config) *App {
	a := new(App)

	if conf.DatabaseDSN != "" {
		conn, err := pgx.Connect(context.Background(), conf.DatabaseDSN)
		if err != nil {
			logger.Log.Fatalln("connect to db: ", err.Error())
		}
		a.dbConn = conn
	}

	if conf.StorageFilePath == "" {
		a.storage = local.New()
	} else {
		s, err := dumper.New(local.New(), conf.StorageFilePath, conf.FileStorageBufferSize)
		if err != nil {
			logger.Log.Fatalln("failed to create storage", err.Error())
		}
		a.storage = s
	}

	a.conf = conf
	a.tokenGenerator = tokengenerator.New(conf.TokenLen)
	a.router = a.configureRouter()

	return a
}

func (a *App) Run() error {
	logger.Log.Infof("Server is running on %s", a.conf.Addr)

	return fmt.Errorf("app err: %w", http.ListenAndServe(a.conf.Addr, a.router))
}

func (a *App) Shutdown() error {
	s, ok := a.storage.(*dumper.Storage)
	if ok {
		if err := s.Dump(); err != nil {
			return fmt.Errorf("dump file storage on closing: %w", err)
		}
	}

	a.dbConn.Close(context.Background())

	return nil
}
