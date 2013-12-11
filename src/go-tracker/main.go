package main

import (
	"fmt"
	"io"
	"net/http"
)

func main() {
	// Set up HTTP routes for various tracker functions
	http.HandleFunc("/announce", announce)
	http.HandleFunc("/scrape", scrape)
	http.HandleFunc("/statistics", statistics)

	// Start HTTP server
	fmt.Println("go-tracker: listening on port 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("go-tracker: failed to start: ", err)
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
