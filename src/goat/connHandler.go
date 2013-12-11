package goat

import (
	"net"
)

// ConnHandler interface method Handle defines how to handle incoming network connections
type ConnHandler interface {
	Handle(net.Conn) bool
}

// HttpConnHandler handles incoming HTTP (TCP) network connections
type HttpConnHandler struct {
}

func (h HttpConnHandler) Handle(http net.Conn) bool {
	return true
}

// UdpConnHandler handles incoming UDP network connections
type UdpConnHandler struct {
}

func (u UdpConnHandler) Handle(udp net.Conn) bool {
	return true
}
