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
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt)
	signal.Notify(sigc, syscall.SIGTERM)
	for sig := range sigc {
		fmt.Println(APP, ": caught signal:", sig)
		killChan <- true
		hold := <-exitChan
		os.Exit(hold)
	}
}
