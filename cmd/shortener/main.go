package main

import (
	"log"

	"github.com/AlexanderVasiliev23/yp-url-shortener/configs"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app"
)

func main() {
	conf := configs.Configure()

	application := app.New(conf)

	if err := application.Run(); err != nil {
		log.Fatalln(err)
	}
}
