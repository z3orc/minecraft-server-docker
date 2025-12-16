package logger

import (
	"log"
	"log/slog"
)

func Init() {
	log.SetPrefix("runner: ")
	// logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	// slog.SetDefault(logger)

	slog.SetLogLoggerLevel(slog.LevelInfo)
}

func SetDebugLogLevel() {
	slog.SetLogLoggerLevel(slog.LevelDebug)
}
