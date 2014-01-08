package goat

import (
	"log"
	"net"
	"os"
	"strconv"
)

// Listen and handle HTTP (TCP) connections
func listenHTTP(httpDoneChan chan bool) {
	// Listen on specified TCP port
	l, err := net.Listen("tcp", ":"+strconv.Itoa(static.Config.Port))
	if err != nil {
		log.Println(err.Error())
		log.Println("Cannot start HTTP server, exiting now.")
		os.Exit(1)
	}

	// Send listener to handler
	go handleHTTP(l, httpDoneChan)
}

// Listen on specified UDP port, accept and handle connections
func listenUDP(udpDoneChan chan bool) {
	// Listen on specified UDP port
	addr, err := net.ResolveUDPAddr("udp", ":"+strconv.Itoa(static.Config.Port))
	l, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Println(err.Error())
	}

	// Send listener to handler
	go handleUDP(l, udpDoneChan)
}
