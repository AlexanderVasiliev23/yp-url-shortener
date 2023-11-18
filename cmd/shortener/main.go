package main

import (
	"go.uber.org/zap"
	"log"

	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/configs"
)

func main() {
	conf := configs.Configure()

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalln(err)
	}
	defer func() { _ = logger.Sync() }()
	sugaredLogger := logger.Sugar()

	application := app.New(conf, sugaredLogger)

	if err := application.Run(); err != nil {
		log.Fatalln(err)
	}
}
