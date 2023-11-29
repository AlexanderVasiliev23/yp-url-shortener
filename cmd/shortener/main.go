package main

import (
	"context"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/configs"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/logger"
	"golang.org/x/sync/errgroup"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	conf := configs.Configure()

	if err := logger.Init(); err != nil {
		panic(err)
	}
	defer logger.Log.Sync()

	application := app.New(conf)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

	g, gCtx := errgroup.WithContext(context.Background())
	g.Go(application.Run)

	select {
	case <-interrupt:
	case <-gCtx.Done():
	}

	if err := application.Shutdown(); err != nil {
		logger.Log.Fatalln(err)
	}
}
