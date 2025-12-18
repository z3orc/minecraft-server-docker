package main

import (
	"log/slog"
	"os"
	"runtime"
	"strconv"

	"github.com/z3orc/minecraft-server-docker/internal/logger"
	"github.com/z3orc/minecraft-server-docker/internal/server"
)

type flags struct {
	gameVersion string
	dataDir     string
	memory      string
	jarName     string
	timeout     int
	useSigKill  bool
	debug       bool
}

func parseFlags() *flags {
	// flags := flags{}

	// flag.StringVar(&flags.gameVersion, "version", "1.21.10",
	// 	"Which version of minecraft the server should run. "+
	// 		"Only works if there does not already exist a server.jar file")

	// flag.StringVar(&flags.dataDir, "dir", "./",
	// 	"Directory of server files. This should be the same location as the server jar.")

	// flag.StringVar(&flags.memory, "memory", "1G",
	// 	"How much memory to allocate to the server. Example: 1000M or 1G")

	// flag.StringVar(&flags.jarName, "jar", "server.jar",
	// 	"Name of server jar that the runner will use")

	// flag.IntVar(&flags.timeout, "timeout", 60,
	// 	"How long to wait for the server to gracefully shut down (in seconds)")

	// flag.BoolVar(&flags.useSigKill, "sigkill", false,
	// 	"Use signal SIGKILL to close server if timeout is reached")

	// flag.BoolVar(&flags.debug, "debug", false, "Use debug mode")

	// flag.Parse()

	gameVersion := os.Getenv("VERSION")
	if gameVersion == "" {
		gameVersion = "1.21.10"
	}

	dataDir := os.Getenv("DATA_DIR")
	if dataDir == "" {
		dataDir = "./"
	}

	memory := os.Getenv("MEMORY")
	if memory == "" {
		memory = "1G"
	}

	jarName := os.Getenv("JAR")
	if jarName == "" {
		jarName = "server.jar"
	}

	timeout, err := strconv.Atoi(os.Getenv("TIMEOUT"))
	if timeout == 0 || err != nil {
		timeout = 60
	}

	sigkill, err := strconv.ParseBool(os.Getenv("SIGKILL"))
	if sigkill != true || err != nil {
		sigkill = false
	}

	debug, err := strconv.ParseBool(os.Getenv("DEBUG"))
	if debug != true || err != nil {
		debug = false
	}

	flags := flags{
		gameVersion: gameVersion,
		dataDir:     "./",
		memory:      memory,
		jarName:     jarName,
		timeout:     timeout,
		useSigKill:  sigkill,
		debug:       debug,
	}

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
	s, err := server.New(flags.gameVersion, flags.dataDir, flags.memory, flags.jarName, flags.timeout, flags.useSigKill)
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
