package goat

import (
	"fmt"
	"net"
	"net/url"
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
	// Close socket upon request completion
	defer c.Close()

	// Read in data from socket
	var buf = make([]byte, 1024)
	c.Read(buf)

	// Get the HTTP request line
	line := strings.Split(string(buf), "\r\n")[0]

	// Validate that client is making an actual HTTP request
	req := strings.Split(line, " ")
	if req[0] != "GET" || req[2] != "HTTP/1.1" {
		fmt.Println(APP, ": malformed HTTP request")
		return false
	}

	// Parse request function
	uri, err := url.Parse(req[1])
	if err != nil {
		fmt.Println(APP, ": malformed HTTP request")
		return false
	}

	fmt.Println("path:", uri.Path)
	query, err := url.ParseQuery(uri.RawQuery)
	if err != nil {
		fmt.Println(APP, ": malformed HTTP request")
		return false
	}

	switch uri.Path {
	// Torrent announce request
	case "/announce":
		return announce(c, query)
	}

	// No matching function call found
	c.Write(http_header("404 Not Found", 0))
	return true
}

// UdpConnHandler handles incoming UDP network connections
type UdpConnHandler struct {
}

func (u UdpConnHandler) Handle(c net.Conn) bool {
	return true
}

// Create a byte array of headers with specified HTTP response and content length
func http_header(res string, length int) []byte {
	return []byte(fmt.Sprintf("HTTP/1.1 %s\r\nServer: %s\r\nContent-Type: text/plain\r\nContent-Length: %d\r\nConnection: close\r\n\r\n", res, APP, length))
}

// Create a byte array of a full HTTP response with headers and a 200 OK status
func http_response(body string) []byte {
	return append(http_header("200 OK", len(body))[:], []byte(body)[:]...)
}

// Announce to tracker
func announce(c net.Conn, query map[string][]string) bool {
	c.Write(http_response("announce successful"))
	return true
}
