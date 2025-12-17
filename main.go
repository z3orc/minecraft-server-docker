package main

import (
	"flag"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"

	"github.com/z3orc/minecraft-server-docker/internal/fabric"
	"github.com/z3orc/minecraft-server-docker/internal/jar"
	"github.com/z3orc/minecraft-server-docker/internal/logger"
	"github.com/z3orc/minecraft-server-docker/internal/properties"
	"github.com/z3orc/minecraft-server-docker/internal/serverexec"
)

type flags struct {
	dataDir    string
	JarName    string
	Timeout    int
	UseSigKill bool
	Debug      bool
}

func parseFlags() *flags {
	flags := flags{}

	dataDirPtr := flag.String("dir", "./", "Directory of server files. This should be the same location as the server jar.")
	jarNamePtr := flag.String("jar", "server.jar", "Name of server jar that the runner will use")
	timeoutPtr := flag.Int("timeout", 60, "How long to wait for the server to gracefully shut down (in seconds)")
	useSigKillPtr := flag.Bool("sigkill", false, "Use signal SIGKILL to close server if timeout is reached")
	debugPtr := flag.Bool("debug", false, "Use debug mode")

	flag.Parse()

	flags.dataDir = *dataDirPtr
	flags.JarName = *jarNamePtr
	flags.Timeout = *timeoutPtr
	flags.UseSigKill = *useSigKillPtr
	flags.Debug = *debugPtr

	return &flags
}

func main() {
	logger.Init()

	if runtime.GOOS != "linux" {
		slog.Error("program only works on Linux systems", "GOOS", runtime.GOOS)
		os.Exit(1)
	}

	flags := parseFlags()
	slog.Debug("current value of flags", "flags", flags)

	if flags.Debug {
		logger.SetDebugLogLevel()
	}

	props := properties.New(filepath.Join(flags.dataDir, "server.properties"))
	err := props.LoadFromEnv()
	if err != nil {
		slog.Error("failed to load values for server.properties from env:", "err", err)
		os.Exit(1)
	}

	err = props.Write()
	if err != nil {
		slog.Error("failed to write values for server.properties to disk:", "err", err)
		os.Exit(1)
	}

	url, err := fabric.GetDownloadUrl("1.21.10")
	if err != nil {
		slog.Error("failed get download url from fabric:", "err", err)
		os.Exit(1)
	}

	err = jar.DownloadServerJar(url, flags.dataDir)
	if err != nil {
		slog.Error("error while downloading server jar", "err", err)
		os.Exit(1)
	}

	server, err := serverexec.New(flags.dataDir, flags.JarName)
	if err != nil {
		slog.Error("failed to initialize server:", "err", err)
		os.Exit(1)
	}
	slog.Debug("current value of server", "server", server)

	server.RedirectStdout(os.Stdout)
	server.SignalCatcher(flags.Timeout, flags.UseSigKill)

	slog.Info("starting server", "dir", flags.dataDir, "jar", flags.JarName, "timeout", flags.Timeout)
	err = server.Run()
	if err != nil {
		slog.Error("failed to start server:", "err", err)
		os.Exit(1)
	}

	exitCode := server.ExitCode()
	slog.Info("server exited", "exit code", exitCode)

	os.Exit(0)
}
