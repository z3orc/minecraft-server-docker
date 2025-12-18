package server

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/z3orc/minecraft-server-docker/internal/data/fabric"
	"github.com/z3orc/minecraft-server-docker/internal/jar"
	"github.com/z3orc/minecraft-server-docker/internal/minecraft/management"
	"github.com/z3orc/minecraft-server-docker/internal/minecraft/properties"
)

type Server struct {
	GameVersion string                 // game version
	DataDir     string                 // path of data directory
	Memory      string                 // mount of memory allocated to the JVM
	JarName     string                 // name of jar
	Properties  *properties.Properties // server properties
	serverExec  *ServerExec            // server executor
}

func New(gameVersion string, dataDir string, memory string, jarName string, timeout int, useSigKill bool) (*Server, error) {
	serverExec, err := NewServerExec(dataDir, jarName, memory)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize server exec: %s", err)
	}

	serverExec.RedirectStdout(os.Stdout)
	serverExec.SignalCatcher(timeout, useSigKill)

	props := properties.New(filepath.Join(dataDir, "server.properties"))
	err = props.LoadFromEnv()
	if err != nil {
		return nil, fmt.Errorf("failed to load values for server.properties from env: %s", err)
	}

	return &Server{
		GameVersion: gameVersion,
		JarName:     jarName,
		Memory:      memory,
		DataDir:     dataDir,
		Properties:  props,
		serverExec:  serverExec,
	}, nil
}

func (s *Server) Start() error {
	slog.Info("preparing to start server")

	// Write properties from env to file
	slog.Info("writing properties to disk")
	err := s.Properties.Write()
	if err != nil {
		return fmt.Errorf("failed to write values for server.properties to disk: %s", err)
	}

	// Download server jar if it does not exist
	_, err = os.Stat(filepath.Join(s.DataDir, s.JarName))
	if err != nil && os.IsNotExist(err) {
		slog.Info("downloading server jar from fabric", "version", s.GameVersion, "jar", s.JarName)

		url, err := fabric.GetDownloadUrl(s.GameVersion)
		if err != nil {
			return fmt.Errorf("failed get download url from fabric: %s", err)
		}

		err = jar.DownloadServerJar(url, s.DataDir, s.JarName)
		if err != nil {
			return fmt.Errorf("error while downloading server jar: %s", err)
		}
	} else if err != nil {
		return fmt.Errorf("failed to check if file '%s' exists: %s", s.JarName, err)
	} else {
		slog.Info("server jar already exists. using existing jar", "jar", s.JarName)
	}

	//Add players to OPs
	ops := os.Getenv("OPS")
	slog.Debug("ops", "env", ops)
	if len(ops) > 0 {
		usernames := strings.SplitSeq(ops, ",")
		for username := range usernames {
			username = strings.TrimSpace(username)
			err := management.AddPlayerToOpsList(username, management.OPS_LIST, s.DataDir)
			if err != nil {
				return err
			}
		}
	}

	//Add players to whitelist
	whitelist := os.Getenv("WHITELIST")
	slog.Debug("whitelist", "env", whitelist)
	if len(ops) > 0 {
		usernames := strings.SplitSeq(whitelist, ",")
		for username := range usernames {
			username = strings.TrimSpace(username)
			err := management.AddPlayerToWhitelist(username, management.WHITELIST, s.DataDir)
			if err != nil {
				return err
			}
		}
	}

	// Run server based on serverExec
	slog.Info("starting server")
	return s.serverExec.Run()
}
