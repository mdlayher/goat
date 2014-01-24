package goat

import (
	"crypto/tls"
	"log"
	"net"
	"os"
	"strconv"
)

// Listen and handle HTTP (TCP) connections
func listenHTTP(sendChan chan bool, recvChan chan bool) {
	// Listen on specified TCP port
	l, err := net.Listen("tcp", ":" + strconv.Itoa(static.Config.Port))
	if err != nil {
		log.Println(err.Error())
		log.Println("Cannot start HTTP server, exiting now.")
		os.Exit(1)
	}

	// Send listener to handler
	go handleHTTP(l, sendChan, recvChan)
}

// Listen and handle HTTPS (SSL over TCP) connections
func listenHTTPS(sendChan chan bool, recvChan chan bool) {
	// Load certificate and key
	cert, err := tls.LoadX509KeyPair(static.Config.SSL.Certificate, static.Config.SSL.Key)
	if err != nil {
		log.Println(err.Error())
		log.Println("Cannot load HTTPS X509 key pair, exiting now.")
		os.Exit(1)
	}

	// SSL configuration
	sslConfig := tls.Config{
		Certificates: []tls.Certificate{cert},
	}

	// Listen on specified SSL port
	l, err := tls.Listen("tcp", ":" + strconv.Itoa(static.Config.SSL.Port), &sslConfig)
	if err != nil {
		log.Println(err.Error())
		log.Println("Cannot start HTTPS server, exiting now.")
		os.Exit(1)
	}

	// Send listener to handler
	go handleHTTP(l, sendChan, recvChan)
}

// Listen on specified UDP port, accept and handle connections
func listenUDP(sendChan chan bool, recvChan chan bool) {
	// Listen on specified UDP port
	addr, err := net.ResolveUDPAddr("udp", ":" + strconv.Itoa(static.Config.Port))
	l, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Println(err.Error())
	}

	// Send listener to handler
	go handleUDP(l, sendChan, recvChan)
}
