package main

import (
	"flag"
)

type Flags struct {
	dataDir    string
	JarName    string
	Timeout    int
	UseSigKill bool
	Debug      bool
}

func ParseFlags() *Flags {
	flags := Flags{}

	dataDirPtr := flag.String("dir", "./", "Directory of server files. This should be the same location as the server jar.")
	jarNamePtr := flag.String("jar", "server.jar", "Name of server jar that the runner will use")
	timeoutPtr := flag.Int("timeout", 60, "How long to wait for the server to gracefully shut down (in seconds)")
	useSigKillPtr := flag.Bool("sigkill", false, "Use signal SIGKILL to close server if timeout is reached")
	debugPtr := flag.Bool("debug", false, "Use debug mode")

	flag.Parse()

	flags.dataDir = *dataDirPtr
	flags.JarName = *jarNamePtr
	flags.Timeout = *timeoutPtr
	flags.UseSigKill = *useSigKillPtr
	flags.Debug = *debugPtr

	return &flags
}
