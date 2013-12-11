package goat

import (
	"fmt"
	"net"
)

// Listener interface method Listen defines a network listener which accepts connections
type Listener interface {
	Listen(in, out chan bool, port string)
}

// HttpListener listens for HTTP (TCP) connections
type HttpListener struct {
}

// Listen on specified TCP port, accept and handle connections
func (h HttpListener) Listen(in, out chan bool, port string) {
	l, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println(err)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println(err)
		}
		go new (HttpConnHandler).Handle(conn)
	}
}

// UdpListener listens for UDP connections
type UdpListener struct {
}

// Listen on specified UDP port, accept and handle connections
func (u UdpListener) Listen(in, out chan bool) {

}
