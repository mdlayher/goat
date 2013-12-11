package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

const APP = "goat"

func main() {
	// Create goroutine to handle termination via UNIX signal
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt)
	signal.Notify(sigc, syscall.SIGTERM)
	go func() {
		for sig := range sigc {
			fmt.Println(APP, ": caught signal:", sig)
			os.Exit(0)
		}
	}()
}
