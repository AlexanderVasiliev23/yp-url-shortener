package logger

import (
	"fmt"

	"go.uber.org/zap"
)

// Log missing godoc.
var Log *zap.SugaredLogger

// Init missing godoc.
func Init() error {
	l, err := zap.NewProduction()
	if err != nil {
		return fmt.Errorf("logger init: %w", err)
	}

	Log = l.Sugar()

	return err
}
