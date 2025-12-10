package main

import (
	"log"
	"log/slog"
	"os"
)

func InitLogger() {
	log.SetPrefix("runner: ")
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)
	// slog.SetDefault(logger.With("source", "runner"))

	slog.SetLogLoggerLevel(slog.LevelInfo)
}

func SetDebugLogLevel() {
	slog.SetLogLoggerLevel(slog.LevelDebug)
}
