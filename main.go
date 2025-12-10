package main

import (
	"fmt"
	"os"
	"runtime"
)

func main() {
	if runtime.GOOS != "linux" {
		fmt.Println("error: Program only works on Linux systems")
		os.Exit(1)
	}

	flags := ParseFlags()
	fmt.Println("runner: flags=", flags)

	server, err := NewServer(flags.JarPath)
	if err != nil {
		fmt.Printf("runner: Failed to initialize server: %e", err)
		os.Exit(1)
	}
	fmt.Println("runner: server=", server)

	server.RedirectStdout(os.Stdout)
	server.SignalCatcher(flags.Timeout, flags.UseSigKill)

	fmt.Printf("runner: Starting server. jar=%s, timeout=%d\n", flags.JarPath, flags.Timeout)
	server.Run()

	exitCode := server.ExitCode()
	fmt.Println("runner: Server exited with code:", exitCode)
	os.Exit(0)
}
