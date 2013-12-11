package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Create goroutine to handle termination via UNIX signal
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt)
	signal.Notify(sigc, syscall.SIGTERM)
	go func() {
		for sig := range sigc {
			fmt.Println("Got signal:", sig)
			os.Exit(0)
		}
	}()

	// Set up HTTP routes for various tracker functions
	http.HandleFunc("/announce", announce)
	http.HandleFunc("/scrape", scrape)
	http.HandleFunc("/statistics", statistics)

	// Start HTTP server
	fmt.Println("go-tracker: listening on port 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("go-tracker: failed to start: ", err)
		os.Exit(-1)
	}
}

// Tracker announce handling
func announce(res http.ResponseWriter, req *http.Request) {
	io.WriteString(res, "announce successful")
}

// Tracker scrape handling
func scrape(res http.ResponseWriter, req *http.Request) {
	io.WriteString(res, "scrape successful")
}

// Tracker statistics output
func statistics(res http.ResponseWriter, req *http.Request) {
	io.WriteString(res, "statistics")
}
