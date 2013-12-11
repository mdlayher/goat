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
	// Create goroutine to handle termination via UNIX signal
	killChan := make(chan bool)
	doneChan := make(chan int)
	go goat.Manager(killChan, doneChan)

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt)
	signal.Notify(sigc, syscall.SIGTERM)
	go func() {
		for sig := range sigc {
			fmt.Println(APP, ": caught signal:", sig)
			killChan <- true
			hold := <-doneChan
			os.Exit(hold)
		}
	}()

}
