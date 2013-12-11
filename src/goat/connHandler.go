package goat

import (
	"net"
)

// ConnHandler interface method Handle defines how to handle incoming network connections
type ConnHandler interface {
	Handle(c net.Conn) bool
}

// HttpConnHandler handles incoming HTTP (TCP) network connections
type HttpConnHandler struct {
}

func (h HttpConnHandler) Handle(c net.Conn) bool {
	return true
}

// UdpConnHandler handles incoming UDP network connections
type UdpConnHandler struct {
}

func (u UdpConnHandler) Handle(c net.Conn) bool {
	return true
}
