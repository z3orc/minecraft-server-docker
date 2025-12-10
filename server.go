package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	cmd    *exec.Cmd
	stdout io.Reader
	stdin  io.Writer
}

func NewServer(jarPath string) (*Server, error) {
	cmd := exec.Command("java", "-jar", jarPath, "-nogui")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

	cmdStdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}

	cmdStdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	server := Server{
		cmd:    cmd,
		stdin:  cmdStdin,
		stdout: cmdStdout,
	}

	return &server, nil
}

func (s *Server) Run() error {
	err := s.cmd.Start()
	if err != nil {
		return err
	}

	for {
		dialTimeout := time.Second * 5
		_, err := net.DialTimeout("tcp", net.JoinHostPort("127.0.0.1", "25565"), dialTimeout)
		if err == nil {
			fmt.Println("runner: Server now listening on TCP!")
			break
		}
	}

	return s.cmd.Wait()
}

func (s *Server) ExitCode() int {
	if s.cmd.ProcessState == nil {
		return -1
	}

	return s.cmd.ProcessState.ExitCode()
}

// func (s *Server) Wait() {
// 	s.cmd.Wait()
// }
//
// func (s *Server) isRunning() bool {
// 	if s.cmd.ProcessState == nil {
// 		return false
// 	}
//
// 	err := s.cmd.Process.Signal(syscall.Signal(0))
// 	fmt.Println("process returned:", err)
// 	return true
// }

// func (s *Server) Wait() {
// 	for {
// 		dialTimeout := time.Second * 5
// 		_, err := net.DialTimeout("tcp", net.JoinHostPort("127.0.0.1", "25565"), dialTimeout)
// 		if err == nil {
// 			fmt.Println("runner: Server now listening on TCP!")
// 			return
// 		}
// 	}
// }

func (s *Server) RedirectStdout(dest io.Writer) {
	go func() {
		scanner := bufio.NewScanner(s.stdout)
		for scanner.Scan() {
			fmt.Fprintln(dest, scanner.Text())
		}
	}()
}

func (s *Server) SignalCatcher(timeout int, useSigKill bool) {
	signalChannel := make(chan os.Signal, 0)
	signal.Notify(signalChannel, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		signal := <-signalChannel
		fmt.Println("runner: Got signal: ", signal.String())

		fmt.Println("runner: Sending 'stop' to server")
		fmt.Fprintln(s.stdin, "stop")

		if timeout != -1 {
			time.Sleep(time.Duration(timeout) * time.Second)
		}
		if !useSigKill {
			fmt.Println("runner: Server has not shut down within the time limit; Sending SIGINT")
			s.cmd.Process.Signal(syscall.SIGINT)
		} else {
			fmt.Println("runner: Server has not shut down within the time limit; Sending SIGKILL")
			s.cmd.Process.Kill()
		}
	}()
}
