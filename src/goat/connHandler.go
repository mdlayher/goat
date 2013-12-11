package goat

import (
	"io"
	"net/http"
)

type ConnHandler interface {
	Listen()
}

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
