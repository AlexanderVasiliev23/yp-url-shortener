package main

import (
	"log"

	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/configs"
)

func main() {
	conf := configs.Configure()

	application := app.New(conf)

	if err := application.Run(); err != nil {
		log.Fatalln(err)
	}
}
