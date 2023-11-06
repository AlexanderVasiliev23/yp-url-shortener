package main

import (
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/handlers/add"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/handlers/get"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/storage"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/tokengenerator"
	"github.com/labstack/echo/v4"
	"net/http"
)

const (
	addr = "localhost:8080"
)

func main() {
	localStorage := storage.NewLocalStorage()
	tokenGenerator := tokengenerator.New()

	e := echo.New()

	e.POST("/", add.Add(localStorage, tokenGenerator, addr))
	e.GET("/:token", get.Get(localStorage))

	if err := http.ListenAndServe(addr, e); err != nil {
		panic(err)
	}
}
