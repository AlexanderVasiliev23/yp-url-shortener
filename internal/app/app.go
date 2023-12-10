package app

import (
	"context"
	"fmt"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/configs"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/logger"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/storage"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/storage/dumper"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/storage/local"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/storage/postgres"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/tokengenerator"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"net/http"
	"os"
)

type App struct {
	conf           *configs.Config
	storage        storage.Storage
	tokenGenerator *tokengenerator.TokenGenerator
	router         *echo.Echo
	dbConn         *pgx.Conn
}

func New(ctx context.Context, conf *configs.Config) *App {
	a := new(App)

	a.conf = conf

	conn, err := pgx.Connect(ctx, a.conf.DatabaseDSN)
	if err != nil {
		logger.Log.Fatalln("connect to db: ", err.Error())
		os.Exit(1)
	}
	a.dbConn = conn

	storageObj, err := a.buildStorage(ctx)
	if err != nil {
		logger.Log.Fatalln("creating storage: ", err.Error())
		os.Exit(1)
	}
	a.storage = storageObj
	a.tokenGenerator = tokengenerator.New(conf.TokenLen)
	a.router = a.configureRouter()

	return a
}

func (a *App) buildStorage(ctx context.Context) (storage.Storage, error) {
	if a.dbConn != nil {
		pgStorage, err := postgres.New(ctx, a.dbConn)
		if err != nil {
			return nil, fmt.Errorf("creating pg storage: %w", err)
		}

		return pgStorage, nil
	}

	if a.conf.StorageFilePath != "" {
		s, err := dumper.New(ctx, local.New(), a.conf.StorageFilePath, a.conf.FileStorageBufferSize)
		if err != nil {
			return nil, fmt.Errorf("creating file dumpres storage: %w", err)
		}

		return s, nil
	}

	return local.New(), nil
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
