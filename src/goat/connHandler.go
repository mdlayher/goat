package goat

import (
	"fmt"
	"net"
	"strings"
)

// ConnHandler interface method Handle defines how to handle incoming network connections
type ConnHandler interface {
	Handle(c net.Conn) bool
}

// HttpConnHandler handles incoming HTTP (TCP) network connections
type HttpConnHandler struct {
}

// Handle an incoming HTTP request and provide a HTTP response
func (h HttpConnHandler) Handle(c net.Conn) bool {
	// Read in data from socket
	var buf = make([]byte, 1024)
	c.Read(buf)

	// TODO: remove temporary printing and fake response
	fmt.Println("http: ", string(buf))
	res := []string {
		"HTTP/1.1 200 OK\r\n",
		"Content-Type: text/plain\r\n",
		"Content-Length: 4\r\n",
		"Connection: close\r\n\r\n",
		"goat\r\n",
	}

	// Write response
	c.Write(strings.Join(res, ""))
	c.Close()

	return true
}

// UdpConnHandler handles incoming UDP network connections
type UdpConnHandler struct {
}

func (u UdpConnHandler) Handle(c net.Conn) bool {
	return true
}
