package goat

import (
	"net"
)

// Listener interface method Listen defines a network listener which accepts connections
type Listener interface {
	Listen(port string, logChan chan string)
}

// HttpListener listens for HTTP (TCP) connections
type HttpListener struct {
}

// Listen and handle HTTP (TCP) connections
func (h HttpListener) Listen(port string, logChan chan string) {
	// Listen on specified TCP port
	l, err := net.Listen("tcp", ":"+port)
	if err != nil {
		logChan <- err.Error()
	}

	// Send listener to HttpConnHandler
	go new(HttpConnHandler).Handle(l, logChan)
}

// UdpListener listens for UDP connections
type UdpListener struct {
}

// Listen on specified UDP port, accept and handle connections
func (u UdpListener) Listen(port string, logChan chan string) {

}
