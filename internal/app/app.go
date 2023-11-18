package app

import (
	"fmt"
	"go.uber.org/zap"
	"net/http"

	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/configs"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/storage"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/tokengenerator"
	"github.com/labstack/echo/v4"
)

type App struct {
	conf           *configs.Config
	localStorage   *storage.LocalStorage
	tokenGenerator *tokengenerator.TokenGenerator
	router         *echo.Echo
	logger         *zap.SugaredLogger
}

func New(conf *configs.Config, logger *zap.SugaredLogger) *App {
	a := new(App)

	a.conf = conf
	a.logger = logger
	a.localStorage = storage.NewLocalStorage()
	a.tokenGenerator = tokengenerator.New(conf.TokenLen)
	a.router = a.configureRouter()

	return a
}

func (a *App) Run() error {
	a.logger.Infof("Server is running on %s", a.conf.Addr)

	return fmt.Errorf("app err: %w", http.ListenAndServe(a.conf.Addr, a.router))
}
