package main

import (
	"log/slog"
	"os"
	"runtime"
)

func main() {
	InitLogger()

	if runtime.GOOS != "linux" {
		slog.Error("program only works on Linux systems", "GOOS", runtime.GOOS)
		os.Exit(1)
	}

	flags := ParseFlags()
	slog.Debug("current value of flags", "flags", flags)

	if flags.Debug {
		SetDebugLogLevel()
	}

	server, err := NewServer(flags.dataDir, flags.JarName)
	if err != nil {
		slog.Error("failed to initialize server", "err", err)
		os.Exit(1)
	}
	slog.Debug("current value of server", "server", server)

	server.RedirectStdout(os.Stdout)
	server.SignalCatcher(flags.Timeout, flags.UseSigKill)

	slog.Info("starting server", "dir", flags.dataDir, "jar", flags.JarName, "timeout", flags.Timeout)
	err = server.Run()
	if err != nil {
		slog.Error("failed to start server", "err", err)
	}

	exitCode := server.ExitCode()
	slog.Info("server exited", "exit code", exitCode)

	os.Exit(0)
}
