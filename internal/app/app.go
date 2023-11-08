package app

import (
	"fmt"
	"net/http"

	"github.com/AlexanderVasiliev23/yp-url-shortener/configs"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/storage"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/tokengenerator"
	"github.com/labstack/echo/v4"
)

type app struct {
	conf           *configs.Config
	localStorage   *storage.LocalStorage
	tokenGenerator *tokengenerator.TokenGenerator
	router         *echo.Echo
}

func New(conf *configs.Config) *app {
	a := new(app)

	a.conf = conf
	a.localStorage = storage.NewLocalStorage()
	a.tokenGenerator = tokengenerator.New(conf.TokenLen)
	a.router = a.configureRouter()

	return a
}

func (a *app) Run() error {
	fmt.Println("Server is running on", a.conf.Addr)

	return fmt.Errorf("app err: %w", http.ListenAndServe(a.conf.Addr, a.router))
}
