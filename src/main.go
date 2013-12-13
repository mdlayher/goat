package main

import (
	"fmt"
	"goat"
	"os"
	"os/signal"
	"syscall"
)

const APP = "goat"

func main() {
	// Launch manager via goroutine
	killChan := make(chan bool)
	exitChan := make(chan int)
	go goat.Manager(killChan, exitChan)

	// Gracefully handle termination via UNIX signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, syscall.SIGTERM)
	for sig := range sigChan {
		// Trigger manager shutdown
		fmt.Println(APP, ": caught signal:", sig)
		killChan <- true
		break
	}

	// Force terminate if signaled twice
	go func(sigChan chan os.Signal) {
		for sig := range sigChan {
			_ = sig
			fmt.Println(APP, ": force halting now!")
			os.Exit(1)
		}
	}(sigChan)

	// Exit with specified code from manager
	code := <-exitChan
	fmt.Println(APP, ": graceful shutdown complete")
	os.Exit(code)
}
