package goat

import (
	"bencode"
	"net"
	"net/http"
)

// ConnHandler interface method Handle defines how to handle incoming network connections
type ConnHandler interface {
	Handle(l net.Listener) bool
}

// HttpConnHandler handles incoming HTTP (TCP) network connections
type HttpConnHandler struct {
}

// Handle incoming HTTP connections and serve
func (h HttpConnHandler) Handle(l net.Listener, logChan chan string) bool {
	http.HandleFunc("/announce", parseHttp)

	err := http.Serve(l, nil)
	if err != nil {
		logChan <- err.Error()
	}

	return true
}

// Parse incoming HTTP connections before making tracker calls
func parseHttp(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server", APP+"-git")

	// Failure response for announce to wrong location
	res := map[string][]byte {
		"failure reason": bencode.EncBytes([]byte("Malformed announce")),
		"min interval": bencode.EncInt(1800),
		"interval": bencode.EncInt(3600),
	}
	w.Write(bencode.EncDictMap(res))
}

// UdpConnHandler handles incoming UDP network connections
type UdpConnHandler struct {
}

// Handle incoming UDP connections and return response
func (u UdpConnHandler) Handle(l net.Listener) bool {
	return true
}
