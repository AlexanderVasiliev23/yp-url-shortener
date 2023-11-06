package main

import (
	"fmt"
	"github.com/AlexanderVasiliev23/yp-url-shortener/config"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/handlers/add"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/handlers/get"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/storage"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/tokengenerator"
	"github.com/labstack/echo/v4"
	"net/http"
)

func main() {
	conf := config.Configure()

	localStorage := storage.NewLocalStorage()
	tokenGenerator := tokengenerator.New()

	e := echo.New()

	e.POST("/", add.Add(localStorage, tokenGenerator, conf.BaseAddress))
	e.GET("/:token", get.Get(localStorage))

	fmt.Println("Server is running on", conf.Addr)
	if err := http.ListenAndServe(conf.Addr, e); err != nil {
		panic(err)
	}
}
