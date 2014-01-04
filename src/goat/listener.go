package goat

import (
	"net"
	"strconv"
)

// Listener interface method Listen defines a network listener which accepts connections
type Listener interface {
	Listen(chan bool)
}

// HTTPListener listens for HTTP (TCP) connections
type HTTPListener struct {
}

// Listen and handle HTTP (TCP) connections
func (h HTTPListener) Listen(httpDoneChan chan bool) {
	// Listen on specified TCP port
	l, err := net.Listen("tcp", ":"+strconv.Itoa(Static.Config.Port))
	if err != nil {
		Static.LogChan <- err.Error()
	}

	// Send listener to HttpConnHandler
	go new(HTTPConnHandler).Handle(l, httpDoneChan)
}

// UDPListener listens for UDP connections
type UDPListener struct {
}

// Listen on specified UDP port, accept and handle connections
func (u UDPListener) Listen(udpDoneChan chan bool) {
	// Listen on specified UDP port
	addr, err := net.ResolveUDPAddr("udp", ":"+strconv.Itoa(Static.Config.Port))
	l, err := net.ListenUDP("udp", addr)
	if err != nil {
		Static.LogChan <- err.Error()
	}

	// Send listener to UdpConnHandler
	go new(UDPConnHandler).Handle(l, udpDoneChan)
}
