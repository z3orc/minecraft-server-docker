package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strconv"
	"syscall"
	"time"
)

func print_help() {
	fmt.Printf("usage: %s <executable> <timeout>\n- executable: Path to server jar file\n- timeout: How long to for server to gracefully close\n", os.Args[0])
}

func main() {
	if runtime.GOOS != "linux" {
		fmt.Println("error: Program only works on Linux systems")
		os.Exit(1)
	}

	if len(os.Args) < 3 {
		fmt.Println("error: Missing arguments!")
		print_help()
		os.Exit(1)
	}

	jar_path := os.Args[1]
	timeout, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println("error: Failed to parse value for timeout")
		print_help()
		os.Exit(1)
	}

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
	signal.Notify(signalChannel, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		signal := <-signalChannel
		fmt.Println("runner: Got signal: ", signal.String())

		fmt.Println("runner: Sending 'stop' to server")
		fmt.Fprintln(cmdStdin, "stop")

		time.Sleep(time.Duration(timeout) * time.Second)
		fmt.Println("runner: Server has not shut down within the time limit; Sending SIGINT")
		cmd.Process.Signal(syscall.SIGINT)
	}()

	go func() {
		for {
			dialTimeout := time.Second * 5
			_, err := net.DialTimeout("tcp", net.JoinHostPort("127.0.0.1", "25565"), dialTimeout)
			if err == nil {
				fmt.Println("runner: Server now listening on TCP!")
				return
			}
		}
	}()

	fmt.Printf("runner: Starting server. jar=%s, timeout=%d, cpus=%d\n", jar_path, timeout, runtime.NumCPU())
	cmd.Run()

	exitCode := cmd.ProcessState.ExitCode()
	fmt.Println("runner: Server exited with code: ", exitCode)
	os.Exit(0)
}
