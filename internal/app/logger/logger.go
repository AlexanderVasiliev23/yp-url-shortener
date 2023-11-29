package logger

import (
	"fmt"
	"go.uber.org/zap"
)

var Log *zap.SugaredLogger

func Init() error {
	l, err := zap.NewProduction()
	if err != nil {
		return fmt.Errorf("logger init: %w", err)
	}

	Log = l.Sugar()

	return err
}
