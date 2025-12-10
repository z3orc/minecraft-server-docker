package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
)

func print_help() {
	fmt.Printf("usage: %s <executable> <timeout>\n- executable: Path to server jar file\n- timeout: How long to for server to gracefully close\n", os.Args[0])
}

type state struct {
	JarPath    string
	Timeout    int
	UseSigKill bool
}

func parse_args(state *state) {
	jarPtr := flag.String("jar", "server.jar", "Path to server.jar")
	timeoutPtr := flag.Int("timeout", 60, "How long to wait for the server to gracefully shut down (in seconds)")
	useSigKill := flag.Bool("sigkill", false, "Use signal SIGKILL to close server if timeout is reached")

	flag.Parse()

	state.JarPath = *jarPtr
	state.Timeout = *timeoutPtr
	state.UseSigKill = *useSigKill
}

func main() {
	if runtime.GOOS != "linux" {
		fmt.Println("error: Program only works on Linux systems")
		os.Exit(1)
	}

	state := state{}
	parse_args(&state)
	fmt.Println(state)

	server, err := NewServer(state.JarPath)
	if err != nil {
		fmt.Printf("runner: Failed to initialize server: %e", err)
		os.Exit(1)
	}

	server.RedirectStdout(os.Stdout)
	server.SignalCatcher(state.Timeout, state.UseSigKill)

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

	fmt.Printf("runner: Starting server. jar=%s, timeout=%d, cpus=%d\n", state.JarPath, state.Timeout, runtime.NumCPU())
	server.Run()
	exitCode := server.ExitCode()
	fmt.Println("runner: Server exited with code:", exitCode)
	os.Exit(0)
}
