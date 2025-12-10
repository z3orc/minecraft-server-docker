package main

import (
	"log"
	"log/slog"
)

func InitLogger() {
	log.SetPrefix("runner: ")
	slog.SetLogLoggerLevel(slog.LevelInfo)
}

func SetDebugLogLevel() {
	slog.SetLogLoggerLevel(slog.LevelDebug)
}
