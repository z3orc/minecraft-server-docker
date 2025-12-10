package main

import "flag"

type Flags struct {
	JarPath    string
	Timeout    int
	UseSigKill bool
}

func ParseFlags() *Flags {
	flags := Flags{}

	jarPtr := flag.String("jar", "server.jar", "Path to server.jar")
	timeoutPtr := flag.Int("timeout", 60, "How long to wait for the server to gracefully shut down (in seconds)")
	useSigKill := flag.Bool("sigkill", false, "Use signal SIGKILL to close server if timeout is reached")

	flag.Parse()

	flags.JarPath = *jarPtr
	flags.Timeout = *timeoutPtr
	flags.UseSigKill = *useSigKill

	return &flags
}
