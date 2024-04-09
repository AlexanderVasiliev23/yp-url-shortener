package logger

import (
	"context"
	"github.com/AlexanderVasiliev23/yp-url-shortener/internal/app/logger"
	"google.golang.org/grpc"
	"time"
)

func UnaryInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	start := time.Now()

	resp, err = handler(ctx, req)

	duration := time.Since(start)

	if err != nil {
		logger.Log.Error(err.Error())
	}

	logger.Log.Infow(
		"GRPC request handled",
		"method", info.FullMethod,
		"duration", duration,
	)

	return
}
