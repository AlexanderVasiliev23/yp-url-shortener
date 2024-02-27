package main

import (
	"context"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"

	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/configs"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/logger"
)

func main() {
	conf := configs.MustConfigure()

	if err := logger.Init(); err != nil {
		panic(err)
	}
	defer logger.Log.Sync()

	ctx := context.Background()

	application := app.New(ctx, conf)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

	g, gCtx := errgroup.WithContext(ctx)
	g.Go(application.RunWorkers)
	g.Go(application.Run)
	profilerAddr := ":8081"
	g.Go(func() error {
		return http.ListenAndServe(profilerAddr, nil)
	})

	select {
	case <-interrupt:
	case <-gCtx.Done():
	}

	if err := application.Shutdown(); err != nil {
		logger.Log.Fatalln(err)
	}
}
