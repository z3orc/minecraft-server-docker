package main

import (
	"fmt"
	"os"
	"runtime"
)

func print_help() {
	fmt.Printf("usage: %s <executable> <timeout>\n- executable: Path to server jar file\n- timeout: How long to for server to gracefully close\n", os.Args[0])
}

func main() {
	if runtime.GOOS != "linux" {
		fmt.Println("error: Program only works on Linux systems")
		os.Exit(1)
	}

	flags := ParseFlags()
	fmt.Println(flags)

	server, err := NewServer(flags.JarPath)
	if err != nil {
		fmt.Printf("runner: Failed to initialize server: %e", err)
		os.Exit(1)
	}

	server.RedirectStdout(os.Stdout)
	server.SignalCatcher(flags.Timeout, flags.UseSigKill)

	// if len(os.Args) < 3 {
	// 	fmt.Println("error: Missing arguments!")
	// 	print_help()
	// 	os.Exit(1)
	// }

	// jar_path := os.Args[1]
	// timeout, err := strconv.Atoi(os.Args[2])
	// if err != nil {
	// 	fmt.Println("error: Failed to parse value for timeout")
	// 	print_help()
	// 	os.Exit(1)
	// }

	fmt.Printf("runner: Starting server. jar=%s, timeout=%d, cpus=%d\n", flags.JarPath, flags.Timeout, runtime.NumCPU())
	server.Run()
	exitCode := server.ExitCode()
	fmt.Println("runner: Server exited with code:", exitCode)
	os.Exit(0)
}
