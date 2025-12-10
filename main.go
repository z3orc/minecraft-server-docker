package main

import (
	"fmt"
	"os"
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

	// fmt.Println("Hello, World")
	//
	// signalChannel := make(chan os.Signal, 1)
	// signal.Notify(signalChannel, syscall.SIGTERM, syscall.SIGINT)
	//
	// go func() {
	// 	signal := <-signalChannel
	// 	fmt.Println("Got signal: ", signal.String())
	// 	os.Exit(0)
	// }()
	//
	// time.Sleep(10 * time.Second)
}
