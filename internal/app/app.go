package app

import (
	"context"
	"fmt"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/util/tls"
	"net/http"
	"os"

	"github.com/AlexanderVasiliev23/yp-url-shortener/pkg/tokengenerator"

	zap "github.com/jackc/pgx-zap"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
	"github.com/labstack/echo/v4"

	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/configs"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/logger"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/storage"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/storage/dumper"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/storage/local"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/storage/postgres"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/util/auth"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/uuidgenerator"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/uuidgenerator/google"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/workers/deleter"
)

// App missing godoc.
type App struct {
	conf               *configs.Config
	storage            storage.Storage
	tokenGenerator     *tokengenerator.TokenGenerator
	router             *echo.Echo
	dbConn             *pgxpool.Pool
	uuidGenerator      uuidgenerator.UUIDGenerator
	userContextFetcher *auth.UserContextFetcher

	deleteByTokenCh chan deleter.DeleteTask
}

// New missing godoc.
func New(ctx context.Context, conf *configs.Config) *App {
	a := new(App)

	a.deleteByTokenCh = make(chan deleter.DeleteTask)
	a.conf = conf
	a.uuidGenerator = google.UUIDGenerator{}
	a.userContextFetcher = &auth.UserContextFetcher{}

	l := zap.NewLogger(logger.Log.Desugar())

	if a.conf.DatabaseDSN != "" {
		dbConfig, err := pgxpool.ParseConfig(a.conf.DatabaseDSN)
		if err != nil {
			logger.Log.Fatalln("parse db dns for config:", err.Error())
			os.Exit(1)
		}
		if a.conf.Debug {
			dbConfig.ConnConfig.Tracer = &tracelog.TraceLog{
				Logger:   l,
				LogLevel: tracelog.LogLevelTrace,
			}
		}

		pool, err := pgxpool.NewWithConfig(ctx, dbConfig)
		if err != nil {
			logger.Log.Fatalln("connect to db:", err.Error())
			os.Exit(1)
		}
		a.dbConn = pool
	}

	storageObj, err := a.buildStorage(ctx)
	if err != nil {
		logger.Log.Fatalln("creating storage:", err.Error())
		os.Exit(1)
	}
	a.storage = storageObj
	a.tokenGenerator = tokengenerator.New(conf.TokenLen)
	a.router = a.configureRouter(ctx)

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
		s, err := dumper.New(ctx, local.New(a.uuidGenerator), a.uuidGenerator, a.conf.StorageFilePath, a.conf.FileStorageBufferSize)
		if err != nil {
			return nil, fmt.Errorf("creating file dumper storage: %w", err)
		}

		return s, nil
	}

	return local.New(a.uuidGenerator), nil
}

// Run missing godoc.
func (a *App) Run() error {
	logger.Log.Infof("Server is running on %s", a.conf.Addr)

	if a.conf.EnableHTTPS {
		if !tls.PemFilesExist() {
			if err := tls.CreatePemFiles(); err != nil {
				return fmt.Errorf("generate pem files: %w", err)
			}
		}
		return fmt.Errorf("app err: %w", http.ListenAndServeTLS(a.conf.Addr, tls.CertFilePath, tls.KeyFilePath, a.router))
	} else {
		return fmt.Errorf("app err: %w", http.ListenAndServe(a.conf.Addr, a.router))
	}
}

// RunWorkers missing godoc.
func (a *App) RunWorkers() error {
	deleteWorker := deleter.NewDeleteWorker(a.storage, deleter.Options{
		RepoDeletionTimeout: a.conf.DeleteWorkerConfig.RepoTimeout,
	})
	deleteWorker.Consume(a.deleteByTokenCh)

	return nil
}

// Shutdown missing godoc.
func (a *App) Shutdown() error {
	s, ok := a.storage.(*dumper.Storage)
	if ok {
		if err := s.Dump(); err != nil {
			return fmt.Errorf("dump file storage on closing: %w", err)
		}
	}

	if a.dbConn != nil {
		a.dbConn.Close()
	}

	close(a.deleteByTokenCh)

	return nil
}
