package logger

import (
	"go.uber.org/zap"
)

// Log global variable of logger
var Log *zap.SugaredLogger

// Init logger
func Init(level string) error {
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return err
	}
	cfg := zap.NewProductionConfig()
	cfg.Level = lvl
	zl, err := cfg.Build()
	if err != nil {
		return err
	}
	Log = zl.Sugar()
	return nil
}
