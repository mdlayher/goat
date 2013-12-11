package goat

import (
	"net"
	"net/http"
)

// Listener interface method Listen defines a network listener which accepts connections
type Listener interface {
	Listen(port string, logChan chan string)
}

// HttpListener listens for HTTP (TCP) connections
type HttpListener struct {
}

// Listen on specified TCP port, accept and handle connections
func (h HttpListener) Listen(port string, logChan chan string) {
	l, err := net.Listen("tcp", ":"+port)
	if err != nil {
		logChan <- err.Error()
	}
	for {
		// Tracker announce requests
		http.HandleFunc("/announce", parseHttp)

		// Start listening and serving HTTP requests
		err := http.Serve(l, nil)
		if err != nil {
			logChan <- err.Error()
		}
	}
}

// Parse HTTP request to be handled by standard tracker functions
func parseHttp(w http.ResponseWriter, r *http.Request) {
	// Add header identifying server
	w.Header().Add("Server", APP+"-git")
	w.Write([]byte("announce successful"))
}

// UdpListener listens for UDP connections
type UdpListener struct {
}

// Listen on specified UDP port, accept and handle connections
func (u UdpListener) Listen(port string, logChan chan string) {

}
