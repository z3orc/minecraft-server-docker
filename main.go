package main

import (
	"bufio"
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
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}
	cmdStdin, err := cmd.StdinPipe()
	if err != nil {
		fmt.Printf("error: Failed to open Stdin Pipe: %e\n", err)
		cmd.Process.Kill()
		os.Exit(-1)
	}

	cmdStdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Printf("error: Failed to open Stdout Pipe: %e\n", err)
		cmd.Process.Kill()
		os.Exit(-1)
	}

	go func() {
		scanner := bufio.NewScanner(cmdStdout)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	}()

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGTERM, syscall.SIGINT, os.Interrupt)

	go func() {
		signal := <-signalChannel
		fmt.Println("runner: Got signal: ", signal.String())
		fmt.Println("runner: Ignoring signal")
		// cmd.Process.Signal(syscall.SIGINT)
		fmt.Fprintln(cmdStdin, "stop")
	}()

	cmd.Run()

	cmdStdin.Close()
	cmdStdout.Close()

	exitCode := cmd.ProcessState.ExitCode()
	fmt.Println("runner: Server exited with code: ", exitCode)
}
