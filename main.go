package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	fmt.Println("Hello, World")

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		signal := <-signalChannel
		fmt.Println("Got signal: ", signal.String())
		os.Exit(0)
	}()

	time.Sleep(10 * time.Second)
}
