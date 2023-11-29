package main

import (
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/configs"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/logger"
)

func main() {
	conf := configs.Configure()

	if err := logger.Init(); err != nil {
		panic(err)
	}
	defer logger.Log.Sync()

	application := app.New(conf)

	if err := application.Run(); err != nil {
		logger.Log.Fatalln(err)
	}
}
