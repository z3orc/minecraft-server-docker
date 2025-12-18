package main

import (
	"flag"
	"log/slog"
	"os"
	"runtime"

	"github.com/z3orc/minecraft-server-docker/internal/logger"
	"github.com/z3orc/minecraft-server-docker/internal/server"
)

type flags struct {
	gameVersion string
	dataDir     string
	jarName     string
	timeout     int
	useSigKill  bool
	debug       bool
}

func parseFlags() *flags {
	flags := flags{}

	flag.StringVar(&flags.gameVersion, "version", "1.21.10",
		"Which version of minecraft the server should run. "+
			"Only works if there does not already exist a server.jar file")

	flag.StringVar(&flags.dataDir, "dir", "./",
		"Directory of server files. This should be the same location as the server jar.")

	flag.StringVar(&flags.jarName, "jar", "server.jar",
		"Name of server jar that the runner will use")

	flag.IntVar(&flags.timeout, "timeout", 60,
		"How long to wait for the server to gracefully shut down (in seconds)")

	flag.BoolVar(&flags.useSigKill, "sigkill", false,
		"Use signal SIGKILL to close server if timeout is reached")

	flag.BoolVar(&flags.debug, "debug", false, "Use debug mode")

	flag.Parse()

	return &flags
}

func main() {
	logger.Init()

	if runtime.GOOS != "linux" {
		slog.Error("program only works on Linux systems", "GOOS", runtime.GOOS)
		os.Exit(1)
	}

	flags := parseFlags()
	if flags.debug {
		logger.SetDebugLogLevel()
	}
	slog.Debug("value of flags", "flags", flags)

	slog.Debug("flags", "version", flags.gameVersion, "dir", flags.dataDir, "jar", flags.jarName, "timeout", flags.timeout, "useSigKill", flags.useSigKill, "debug", flags.debug)
	s, err := server.New(flags.gameVersion, flags.dataDir, flags.jarName, flags.timeout, flags.useSigKill)
	if err != nil {
		slog.Error("failed to create server instance:", "err", err)
		os.Exit(1)
	}

	err = s.Start()
	if err != nil {
		slog.Error("failed to start server:", "err", err)
		os.Exit(1)
	}

	os.Exit(0)
}
