package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

func print_help() {
	fmt.Printf("usage: %s <executable>\n- executable: Path to server jar file\n", os.Args[0])
}

func main() {
	if len(os.Args) < 2 {
		fmt.Print("Missing arguments!\n\n")
		print_help()
		os.Exit(1)
	}

	fmt.Println("runner: Starting server using jar")

	jar_path := os.Args[1]
	cmd := exec.Command("java", "-jar", jar_path, "-nogui")
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		signal := <-signalChannel
		fmt.Println("Got signal: ", signal.String())
		cmd.Process.Signal(syscall.SIGINT)
	}()

	cmd.Run()
	exitCode := cmd.ProcessState.ExitCode()
	fmt.Println("runner: Server exited with code: ", exitCode)
}
