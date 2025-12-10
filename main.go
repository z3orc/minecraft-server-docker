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

	server, err := NewServer(flags.JarName)
	if err != nil {
		slog.Error("failed to initialize server", "err", err)
		// fmt.Printf("runner: Failed to initialize server: %e", err)
		os.Exit(1)
	}
	slog.Debug("current value of server", "server", server)

	server.RedirectStdout(os.Stdout)
	server.SignalCatcher(flags.Timeout, flags.UseSigKill)

	slog.Info("starting server", "dir", flags.dataDir, "jar", flags.JarName, "timeout", flags.Timeout)
	// fmt.Printf("runner: Starting server. jar=%s, timeout=%d\n", flags.JarPath, flags.Timeout)
	server.Run()

	exitCode := server.ExitCode()
	slog.Info("server exited", "exit code", exitCode)

	// fmt.Println("runner: Server exited with code:", exitCode)
	os.Exit(0)
}
