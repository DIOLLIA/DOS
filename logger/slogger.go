package logger

import (
	"log/slog"
	"os"
	"sync"
)

var (
	L    *slog.Logger
	once sync.Once
)

func InitLogger(logLevel string) {
	var optsLevel slog.Leveler

	switch logLevel {
	case "debug":
		optsLevel = slog.LevelDebug
	case "error":
		optsLevel = slog.LevelError
	case "info":
		optsLevel = slog.LevelInfo
	case "warn":
		optsLevel = slog.LevelWarn
	default:
		optsLevel = slog.LevelInfo
	}

	once.Do(func() {
		opts := &slog.HandlerOptions{Level: optsLevel}
		L = slog.New(slog.NewJSONHandler(os.Stdout, opts))
		slog.SetDefault(L)

		L.Info("application configured to use log", "level", optsLevel)

	})
}
