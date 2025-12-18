package server

import (
	"bufio"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
)

type ServerExec struct {
	cmd    *exec.Cmd // Pointer to command
	stdout io.Reader // io.Reader for stdout of command
	stdin  io.Writer // io.Writer for stdin of command
}

// NewServer creates a new instance of struct Server, based on provided path of server jar.
//
// Inits pipes for stdin and stdout.
//
// Returns a pointer to struct Server, or error if pipes could not be created.
func NewServerExec(dataDir string, jarName string, memory string) (*ServerExec, error) {
	cmd := exec.Command("java", "-Xms", memory, "-Xmx", memory, "-jar", jarName, "-nogui")
	cmd.Dir = dataDir
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

	server := ServerExec{
		cmd:    cmd,
		stdin:  cmdStdin,
		stdout: cmdStdout,
	}

	return &server, nil
}

// Runs command/executable found in struct ServerExec.
func (s *ServerExec) Run() error {
	err := s.cmd.Start()
	if err != nil {
		return err
	}

	return s.cmd.Wait()
}

// Returns exit code for command/executable found in struct Server, same as os.ProccessState.ExitCode(),
// or -1 if no os.ProcessState is found for command.
func (s *ServerExec) ExitCode() int {
	if s.cmd.ProcessState == nil {
		return -1
	}

	return s.cmd.ProcessState.ExitCode()
}

// Redirects output of command stdout to provided io.Writer 'dest'.
func (s *ServerExec) RedirectStdout(dest io.Writer) {
	go func() {
		scanner := bufio.NewScanner(s.stdout)
		for scanner.Scan() {
			fmt.Fprintln(dest, scanner.Text())
		}
	}()
}

// Catches signals SIGTERM and SIGINT, and tires to stop server using command 'stop',
// and wait the provided timeout.
//
// If timeout is reached, SIGINT or SIGKILL is sent to server process. SIGINT is sent
// if 'useSigKill' is false, and SIGKILL is sent if 'useSigKill' is true.
func (s *ServerExec) SignalCatcher(timeout int, useSigKill bool) {
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		signal := <-signalChannel
		fmt.Println("runner: Got signal: ", signal.String())
		slog.Info("received signal", "signal", signal)

		// fmt.Println("runner: Sending 'stop' to server")
		slog.Info("sending command 'stop' to server")
		fmt.Fprintln(s.stdin, "stop")

		if timeout != -1 {
			time.Sleep(time.Duration(timeout) * time.Second)
		}
		if !useSigKill {
			slog.Warn("server has not shut down within time limit. sending SIGINT")
			// fmt.Println("runner: Server has not shut down within the time limit; Sending SIGINT")
			s.cmd.Process.Signal(syscall.SIGINT)
		} else {
			slog.Warn("server has not shut down within time limit. sending SIGKILL")
			// fmt.Println("runner: Server has not shut down within the time limit; Sending SIGKILL")
			s.cmd.Process.Kill()
		}
	}()
}
