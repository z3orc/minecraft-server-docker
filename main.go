package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"syscall"
)

func print_help() {
	fmt.Printf("usage: %s <executable>\n- executable: Path to server jar file\n", os.Args[0])
}

func main() {
	if runtime.GOOS != "linux" {
		fmt.Println("error: Program only works on Linux systems")
		os.Exit(1)
	}

	if len(os.Args) < 2 {
		fmt.Print("Missing arguments!\n\n")
		print_help()
		os.Exit(1)
	}

	jar_path := os.Args[1]
	cmd := exec.Command("java", "-jar", jar_path, "-nogui")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

	// cmdStdin, err := cmd.StdinPipe()
	// if err != nil {
	// 	fmt.Printf("error: Failed to open Stdin Pipe: %e\n", err)
	// 	cmd.Process.Kill()
	// 	os.Exit(-1)
	// }

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
	signal.Notify(signalChannel, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		for {
			signal := <-signalChannel
			fmt.Println("runner: Got signal: ", signal.String())
			cmd.Process.Signal(syscall.SIGINT)
			// fmt.Fprintln(cmdStdin, "stop")
		}
	}()

	fmt.Printf("runner: Starting server using jar=%s", jar_path)
	cmd.Run()

	exitCode := cmd.ProcessState.ExitCode()
	fmt.Println("runner: Server exited with code: ", exitCode)
	os.Exit(0)
}
